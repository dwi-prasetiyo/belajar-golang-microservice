package server

import (
	"user-service/env"
	"user-service/src/common/log"
	"user-service/src/factory"
	"user-service/src/restful/handler"
	"user-service/src/restful/middleware"
	"user-service/src/restful/router"

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

	authHandler := handler.NewAuth(f)
	router.Auth(app, authHandler)

	profileHandler := handler.NewProfile(f)
	router.Profile(app, profileHandler, m)

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
