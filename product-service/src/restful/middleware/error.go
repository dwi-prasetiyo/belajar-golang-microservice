package middleware

import (
	"product-service/src/common/dto/response"
	"product-service/src/common/errors"
	"product-service/src/common/log"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) ErrorHandler(c *fiber.Ctx, err error) error {

	if e, ok := err.(*errors.Response); ok {
		return c.Status(e.HttpCode).JSON(response.Common{Message: e.Message})
	}

	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(response.Common{Message: e.Message})
	}

	log.Logger.Error(err.Error())

	return c.Status(500).JSON(response.Common{Message: "internal server error"})
}
