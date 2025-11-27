package handler

import (
	"order-service/src/common/dto/request"
	"order-service/src/common/dto/response"
	"order-service/src/factory"
	"order-service/src/restful/service"

	"github.com/gofiber/fiber/v2"
)

type Order struct {
	service service.Order
}

func NewOrder(f *factory.Factory) *Order {
	return &Order{
		service: service.NewOrder(f),
	}
}

func (h *Order) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	req := new(request.CreateOrder)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	req.Order.UserID = userID

	res, err := h.service.CreateOrder(c, req)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(response.Common{Message: "success create order", Data: res})
}
