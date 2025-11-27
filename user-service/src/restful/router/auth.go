package router

import (
	"user-service/src/restful/handler"

	"github.com/gofiber/fiber/v2"
)

func Auth(app *fiber.App, h *handler.Auth) {
	app.Post("/api/v1/auth/register", h.Register)
	app.Post("/api/v1/auth/register/verify", h.VerifyRegister)
	app.Post("/api/v1/auth/login", h.Login)
	app.Delete("/api/v1/auth/logout", h.Logout)
	app.Patch("/api/v1/auth/refresh-token", h.RefreshToken)
}
