package handler

import (
	"user-service/src/common/dto/response"
	"user-service/src/factory"
	"user-service/src/restful/service"

	"github.com/gofiber/fiber/v2"
)

type Profile struct {
	service service.Profile
}

func NewProfile(f *factory.Factory) *Profile {
	return &Profile{
		service: service.NewProfile(f),
	}
}

func (h *Profile) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(response.Common{Message: "get profile successfully", Data: user})
}
