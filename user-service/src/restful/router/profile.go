package router

import (
	"user-service/src/restful/handler"
	"user-service/src/restful/middleware"

	"github.com/gofiber/fiber/v2"
)

func Profile(app *fiber.App, h *handler.Profile, m *middleware.Middleware) {
	app.Get("/api/v1/profile", m.AccessToken, h.GetProfile)
}
