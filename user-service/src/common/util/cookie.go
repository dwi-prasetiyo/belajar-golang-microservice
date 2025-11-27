package util

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func ClearCookie(c *fiber.Ctx, name, path string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		HTTPOnly: true,
		Path:     path,
		Expires:  time.Now().Add(-time.Hour),
	})
}
		
		