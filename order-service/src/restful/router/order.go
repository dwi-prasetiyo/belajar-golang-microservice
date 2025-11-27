package router

import (
	"order-service/src/restful/handler"
	"order-service/src/restful/middleware"

	"github.com/gofiber/fiber/v2"
)

func Order(app *fiber.App, h *handler.Order, m *middleware.Middleware) {
	app.Post("/api/v1/orders", m.AccessToken, h.Create)
}
