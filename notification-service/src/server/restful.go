package server

import (
	"notification-service/env"
	"notification-service/src/common/log"
	"notification-service/src/factory"
	"notification-service/src/restful/handler"
	"notification-service/src/restful/middleware"
	"notification-service/src/restful/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Restful struct {
	app *fiber.App
}

func NewRestful(f *factory.Factory) *Restful {
	m := middleware.New()

	app := fiber.New(fiber.Config{
		ErrorHandler: m.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	notifHandler := handler.NewNotif(f)
	router.Notif(app, notifHandler, m)

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
