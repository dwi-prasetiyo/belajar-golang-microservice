package middleware

import (
	"notification-service/src/common/dto/response"
	"notification-service/src/common/errors"
	"notification-service/src/common/log"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) ErrorHandler(c *fiber.Ctx, err error) error {
	log.Logger.Error(err.Error())

	if e, ok := err.(*errors.Response); ok {
		return c.Status(e.HttpCode).JSON(response.Common{Message: e.Message})
	}

	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(response.Common{Message: e.Message})
	}

	return c.Status(500).JSON(response.Common{Message: "internal server error"})
}
