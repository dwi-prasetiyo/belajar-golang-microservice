package service

import (
	"product-service/src/common/dto/request"
	"product-service/src/common/errors"
	"product-service/src/common/model"
	v "product-service/src/common/pkg/validator"
	"product-service/src/common/util"
	"product-service/src/factory"
	"product-service/src/repository"
	"product-service/src/restful/client"

	"github.com/gofiber/fiber/v2"
)

type Product interface {
	Create(c *fiber.Ctx, req *request.CreateProduct) (string, error)
	FindMany(c *fiber.Ctx, data *request.FindManyProduct) ([]*model.Product, error)
}

type productImpl struct {
	productRepository repository.Product
	imagekitClient    *client.ImageKit
}

func NewProduct(f *factory.Factory) Product {
	return &productImpl{
		productRepository: f.ProductRepository,
		imagekitClient:    f.ImageKitClient,
	}
}

func (s *productImpl) Create(c *fiber.Ctx, req *request.CreateProduct) (string, error) {
	filename := c.Locals("filename").(string)
	if filename == "" {
		return "", &errors.Response{HttpCode: 400, Message: "File not found"}
	}

	path := "./tmp/" + filename
	defer util.DeleteFile(path)

	res, err := s.imagekitClient.UploadFile(c.Context(), path, filename)
	if err != nil {
		return "", err
	}

	req.ImageID = res.FileId
	req.Image = res.Url

	if err := v.Validate.Struct(req); err != nil {
		return res.FileId, err
	}

	err = s.productRepository.Create(c.Context(), &model.Product{
		Name:        req.Name,
		Sku:         req.Sku,
		ImageID:     req.ImageID,
		Image:       req.Image,
		Price:       req.Price,
		Stock:       req.Stock,
		Length:      req.Length,
		Width:       req.Width,
		Height:      req.Height,
		Weight:      req.Weight,
		Description: req.Description,
	})

	return res.FileId, err
}

func (s *productImpl) FindMany(c *fiber.Ctx, data *request.FindManyProduct) ([]*model.Product, error) {
	if err := v.Validate.Struct(data); err != nil {
		return nil, err
	}

	res, err := s.productRepository.FindMany(c.Context(), data)
	if err != nil {
		return nil, err
	}

	return res, nil
}
