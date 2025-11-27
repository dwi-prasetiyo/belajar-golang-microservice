package middleware

import (
	"os"
	"product-service/src/common/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/h2non/filetype"
)

func (m *Middleware) ValidateProductImage(c *fiber.Ctx) error {
	filename := c.Locals("filename").(string)

	file, err := os.Open("./tmp/" + filename)
	if err != nil {
		return err
	}

	defer file.Close()

	head := make([]byte, 261)
	file.Read(head)

	if !filetype.IsImage(head) {
		return &errors.Response{HttpCode: 400, Message: "Invalid image"}
	}

	return c.Next()
}
