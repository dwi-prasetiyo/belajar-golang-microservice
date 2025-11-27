package middleware

import (
	"product-service/src/common/errors"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Admin(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "ADMIN" {
		return &errors.Response{HttpCode: 403, Message: "Forbidden"}
	}

	return c.Next()
}
