package server

import (
	"product-service/env"
	"product-service/src/common/log"
	"product-service/src/factory"
	"product-service/src/restful/handler"
	"product-service/src/restful/middleware"
	"product-service/src/restful/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Restful struct {
	app *fiber.App
}

func NewRestful(f *factory.Factory) *Restful {
	m := middleware.New(f)

	app := fiber.New(fiber.Config{
		ErrorHandler: m.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(m.Log)
	app.Use(recover.New())

	productHandler := handler.NewProduct(f)
	router.Product(app, m, productHandler)

	return &Restful{
		app: app,
	}
}

func (s *Restful) Start() {
	if err := s.app.Listen(env.Conf.CurrentApp.RestfulAddr); err != nil {
		log.Logger.Fatal(err.Error())
	}
}

func (s *Restful) Stop() {
	if err := s.app.Shutdown(); err != nil {
		log.Logger.Error(err.Error())
	}
}
