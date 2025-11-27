package service

import (
	"time"
	"user-service/src/cache"
	"user-service/src/common/dto/request"
	"user-service/src/common/dto/response"
	"user-service/src/common/errors"
	"user-service/src/common/log"
	"user-service/src/common/model"
	v "user-service/src/common/pkg/validator"
	"user-service/src/common/util"
	"user-service/src/factory"
	"user-service/src/publisher"
	"user-service/src/repository"

	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Register(c *fiber.Ctx, data *request.Register) (email string, err error)
	VerifyRegister(c *fiber.Ctx, data *request.VerifyRegister) error
	Login(c *fiber.Ctx, data *request.Login) (*response.Login, error)
	Logout(c *fiber.Ctx, refreshToken string) error
	RefreshToken(c *fiber.Ctx, refreshToken string) (*response.Token, error)
}

type authImpl struct {
	userRepository         repository.User
	credentialRepository   repository.Credential
	refreshTokenRepository repository.RefreshToken
	cacheRepository        cache.Cache
	rabbitMQPublisher      *publisher.RabbitMQ
}

func NewAuth(f *factory.Factory) Auth {
	return &authImpl{
		userRepository:         f.UserRepository,
		credentialRepository:   f.CredentialRepository,
		refreshTokenRepository: f.RefreshTokenRepository,
		cacheRepository:        f.CacheRepository,
		rabbitMQPublisher:      f.RabbitMQPublisher,
	}
}

func (s *authImpl) Register(c *fiber.Ctx, data *request.Register) (email string, err error) {
	if err := v.Validate.Struct(data); err != nil {
		return "", err
	}

	user, err := s.userRepository.Find(c.Context(), "email = ?", []any{data.Email})
	if err != nil {
		return "", err
	}

	if user != nil {
		return "", &errors.Response{HttpCode: 409, Message: "email already registered"}
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	data.Password = string(bcryptPassword)

	if err := s.cacheRepository.Set(c.Context(), "register:"+data.Email, data, 30*time.Minute); err != nil {
		return "", err
	}

	otp, err := util.GenerateOTP()
	if err != nil {
		return "", err
	}

	sendOtpReq := &request.SendOtp{
		Email: data.Email,
		Otp:   otp,
	}

	if err := s.cacheRepository.Set(c.Context(), "otp:"+data.Email, sendOtpReq, 30*time.Minute); err != nil {
		return "", err
	}

	if err := s.rabbitMQPublisher.Publish("user", "otp", &request.RabbitMQMessage{
		RequestID: c.Locals("request_id"),
		Message:   sendOtpReq,
	}); err != nil {
		return "", err
	}

	return data.Email, nil
}

func (s *authImpl) VerifyRegister(c *fiber.Ctx, data *request.VerifyRegister) error {
	if err := v.Validate.Struct(data); err != nil {
		return err
	}

	otp, err := s.cacheRepository.GetSendOtp(c.Context(), "otp:"+data.Email)
	if err != nil {
		return err
	}

	if otp == nil {
		return &errors.Response{HttpCode: 404, Message: "otp not found"}
	}

	if otp.Otp != data.Otp {
		return &errors.Response{HttpCode: 400, Message: "invalid otp"}
	}

	register, err := s.cacheRepository.GetRegister(c.Context(), "register:"+data.Email)
	if err != nil {
		return err
	}

	if register == nil {
		return &errors.Response{HttpCode: 404, Message: "register not found"}
	}

	userID, err := gonanoid.New()
	if err != nil {
		return err
	}

	if err := s.userRepository.Create(c.Context(), &model.User{ID: userID, Name: register.FullName, Email: register.Email, Role: "USER"}, register.Password); err != nil {
		return err
	}

	if err := s.cacheRepository.Delete(c.Context(), "register:"+data.Email); err != nil {
		log.Logger.Error(err.Error())
	}

	if err := s.cacheRepository.Delete(c.Context(), "otp:"+data.Email); err != nil {
		log.Logger.Error(err.Error())
	}

	return nil
}

func (s *authImpl) Login(c *fiber.Ctx, data *request.Login) (*response.Login, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	user, err := s.userRepository.Find(c.Context(), "email = ?", []any{data.Email})
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &errors.Response{HttpCode: 404, Message: "user not found"}
	}

	credential, err := s.credentialRepository.Find(c.Context(), "u.email = ?", []any{data.Email})
	if err != nil {
		return nil, err
	}

	if credential == nil {
		return nil, &errors.Response{HttpCode: 404, Message: "credential not found"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credential.Password), []byte(data.Password)); err != nil {
		return nil, &errors.Response{HttpCode: 400, Message: "invalid password"}
	}

	accessToken, err := util.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepository.Create(c.Context(), &model.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	return &response.Login{
		Data:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authImpl) Logout(c *fiber.Ctx, refreshToken string) error {
	if err := s.refreshTokenRepository.Delete(c.Context(), refreshToken); err != nil {
		return err
	}

	return nil
}

func (s *authImpl) RefreshToken(c *fiber.Ctx, refreshToken string) (*response.Token, error) {
	token, err := s.refreshTokenRepository.Find(c.Context(), refreshToken)
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, &errors.Response{HttpCode: 404, Message: "refresh token not found"}
	}

	accessToken, err := util.GenerateAccessToken(token.UserID, token.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := util.GenerateRefreshToken(token.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepository.Update(c.Context(), &model.RefreshToken{
		Token:     newRefreshToken,
		UserID:    token.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, refreshToken); err != nil {
		return nil, err
	}

	return &response.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
