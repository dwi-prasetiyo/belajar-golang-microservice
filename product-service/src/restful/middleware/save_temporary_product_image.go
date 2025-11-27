package middleware

import (
	"product-service/src/common/errors"
	"product-service/src/common/util"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) SaveTemporaryProductImage(c *fiber.Ctx) error {
	file, err := c.FormFile("product_image")
	if err != nil {
		return err
	}

	var maxSize int64 = 1 * 1000 * 1000 // 1MB
	if file.Size > maxSize {
		return &errors.Response{HttpCode: 400, Message: "File size is too large"}
	}

	filename := util.CreateUnixFileName(file.Filename)
	path := "./tmp/" + filename

	if err := util.CheckExistDir("./tmp"); err != nil {
		return err
	}

	if err := c.SaveFile(file, path); err != nil {
		return err
	}

	c.Locals("filename", filename)
	return c.Next()
}