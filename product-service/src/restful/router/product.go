package router

import (
	"product-service/src/restful/handler"
	"product-service/src/restful/middleware"

	"github.com/gofiber/fiber/v2"
)

func Product(app *fiber.App, m *middleware.Middleware, h *handler.Product) {
	app.Post("/api/v1/products", m.AccessToken, m.Admin, m.SaveTemporaryProductImage, m.ValidateProductImage, h.Create)
	app.Get("/api/v1/products", m.AccessToken, h.FindMany)
}