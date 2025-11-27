package handler

import (
	"encoding/base64"
	"time"
	"user-service/src/common/dto/request"
	"user-service/src/common/dto/response"
	"user-service/src/common/util"
	"user-service/src/factory"
	"user-service/src/restful/service"

	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	service service.Auth
}

func NewAuth(f *factory.Factory) *Auth {
	return &Auth{
		service: service.NewAuth(f),
	}
}

func (h *Auth) Register(c *fiber.Ctx) error {
	req := new(request.Register)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	logBody := *req
	logBody.Password = "********"
	c.Locals("request_body", logBody)

	email, err := h.service.Register(c, req)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "pending_register",
		Value:    base64.StdEncoding.EncodeToString([]byte(email)),
		HTTPOnly: true,
		Path:     "/api/v1/auth/register/verify",
		Expires:  time.Now().Add(30 * time.Minute),
	})

	return c.Status(200).JSON(response.Common{Message: "request success, please check your email for otp"})
}

func (h *Auth) VerifyRegister(c *fiber.Ctx) error {
	req := new(request.VerifyRegister)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	logBody := *req
	logBody.Otp = "********"
	c.Locals("request_body", logBody)

	email, err := base64.StdEncoding.DecodeString(c.Cookies("pending_register"))
	if err != nil {
		return err
	}

	req.Email = string(email)

	if err := h.service.VerifyRegister(c, req); err != nil {
		return err
	}

	util.ClearCookie(c, "pending_register", "/api/v1/auth/register/verify")

	return c.Status(200).JSON(response.Common{Message: "verify register success"})
}

func (h *Auth) Login(c *fiber.Ctx) error {
	req := new(request.Login)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	logBody := *req
	logBody.Password = "********"
	c.Locals("request_body", logBody)

	user, err := h.service.Login(c, req)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    user.AccessToken,
		HTTPOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(30 * time.Minute),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    user.RefreshToken,
		HTTPOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return c.Status(200).JSON(response.Common{Message: "login success", Data: user.Data})
}

func (h *Auth) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if err := h.service.Logout(c, refreshToken); err != nil {
		return err
	}

	util.ClearCookie(c, "access_token", "/")
	util.ClearCookie(c, "refresh_token", "/")

	return c.Status(200).JSON(response.Common{Message: "logout success"})
}

func (h *Auth) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	res, err := h.service.RefreshToken(c, refreshToken)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		HTTPOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HTTPOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return c.Status(200).JSON(response.Common{Message: "refresh token success"})
}
