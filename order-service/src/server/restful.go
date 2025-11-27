package server

import (
	"order-service/env"
	"order-service/src/common/log"
	"order-service/src/factory"
	"order-service/src/restful/handler"
	"order-service/src/restful/middleware"
	"order-service/src/restful/router"

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

	orderHandler := handler.NewOrder(f)
	router.Order(app, orderHandler, m)

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
