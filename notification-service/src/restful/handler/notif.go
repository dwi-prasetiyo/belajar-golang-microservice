package handler

import (
	"notification-service/src/common/dto/request"
	"notification-service/src/common/dto/response"
	"notification-service/src/common/errors"
	"notification-service/src/factory"
	"notification-service/src/publisher"

	"github.com/gofiber/fiber/v2"
)

type Notif struct {
	rabbitMQPublisher *publisher.RabbitMQ
}

func NewNotif(f *factory.Factory) *Notif {
	return &Notif{
		rabbitMQPublisher: f.RabbitMQPublisher,
	}
}

func (n *Notif) MidtransTx(c *fiber.Ctx) error {
	midtransTx, ok := c.Locals("midtransTx").(*request.MidtransTx)
	if !ok {
		return &errors.Response{HttpCode: 400, Message: "failed to validate midtrans tx"}
	}

	if err := n.rabbitMQPublisher.Publish("notification", "midtrans-tx", midtransTx); err != nil {
		return err
	}

	return c.Status(200).JSON(response.Common{Message: "success"})
}
