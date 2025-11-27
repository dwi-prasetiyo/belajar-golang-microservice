package handler

import (
	"product-service/src/common/dto/request"
	"product-service/src/common/dto/response"
	"product-service/src/common/log"
	"product-service/src/factory"
	"product-service/src/restful/client"
	"product-service/src/restful/service"

	"github.com/gofiber/fiber/v2"
)

type Product struct {
	service service.Product
	ik      *client.ImageKit
}

func NewProduct(f *factory.Factory) *Product {
	return &Product{
		service: service.NewProduct(f),
		ik:      f.ImageKitClient,
	}
}

func (h *Product) Create(c *fiber.Ctx) error {
	req := new(request.CreateProduct)

	if err := c.BodyParser(req); err != nil {
		return err
	}

	fileID, err := h.service.Create(c, req)
	if err != nil {
		if fileID != "" {
			if err := h.ik.DeleteFile(c.Context(), fileID); err != nil {
				log.Logger.Error(err.Error())
			}
		}

		return err
	}

	return c.Status(201).JSON(response.Common{
		Message: "create product success",
	})
}


func (h *Product) FindMany(c *fiber.Ctx) error {
	req := new(request.FindManyProduct)

	if err := c.QueryParser(req); err != nil {
		return err
	}

	res, err := h.service.FindMany(c, req)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(response.Common{Message: "get products successfully", Data: res})
}
