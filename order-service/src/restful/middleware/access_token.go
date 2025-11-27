package middleware

import (
	"context"
	"fmt"
	"order-service/env"
	"order-service/src/common/constant"
	"order-service/src/common/errors"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) AccessToken(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token")

	if accessToken == "" {
		return &errors.Response{HttpCode: 401, Message: "Unauthorized"}
	}

	jwtToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, &errors.Response{HttpCode: 401, Message: "Invalid token"}
		}
		return env.Conf.Jwt.PublicKey, nil
	})

	if err != nil {
		return &errors.Response{HttpCode: 401, Message: err.Error()}
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			return &errors.Response{HttpCode: 401, Message: "Invalid token"}
		}

		c.Locals("user_id", userID)
		c.Locals("role", claims["role"])

		ctx := context.WithValue(c.UserContext(), constant.RequestID, c.Locals("request_id"))
		ctx = context.WithValue(ctx, constant.UserID, userID)

		userBlockInfo, err := m.userClient.CheckUserBlock(ctx, &pb.CheckUserBlockReq{
			UserId: userID,
		})
		if err != nil {
			return &errors.Response{HttpCode: 500, Message: err.Error()}
		}
		if userBlockInfo.IsBlocked {
			return &errors.Response{HttpCode: 403, Message: fmt.Sprintf("User blocked: %s", userBlockInfo.Reason)}
		}
	}

	return c.Next()
}
