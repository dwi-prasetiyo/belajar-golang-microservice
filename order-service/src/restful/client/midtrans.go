package client

import (
	"context"
	"encoding/json"
	"order-service/env"
	"order-service/src/common/dto/request"
	"order-service/src/common/dto/response"
	"order-service/src/common/errors"
	"order-service/src/common/pkg/cbreaker"
	"order-service/src/common/util"

	"github.com/gofiber/fiber/v2"
	"github.com/sony/gobreaker/v2"
)

type Midtrans interface {
	CreateTransaction(ctx context.Context, order *request.MidtransTransaction) (*response.MidtransTx, error)
}

type midtransImpl struct {
	cb *gobreaker.CircuitBreaker[any]
}

func NewMidtrans() Midtrans {
	cb := cbreaker.NewRestful("midtrans")
	return &midtransImpl{cb: cb}
}

func (c *midtransImpl) CreateTransaction(ctx context.Context, order *request.MidtransTransaction) (*response.MidtransTx, error) {
	res, err := c.cb.Execute(func() (any, error) {
		url := env.Conf.Midtrans.ApiHostUrl + "/v1/payment-links"

		a := fiber.AcquireAgent()
		defer fiber.ReleaseAgent(a)

		jsonData, err := json.Marshal(order)
		if err != nil {
			return nil, err
		}

		auth := util.CreateMidtransBasicAuth()

		req := a.Request()
		req.Header.Set("Authorization", auth)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.SetMethod("POST")
		req.SetBodyRaw(jsonData)
		req.SetRequestURI(url)

		if err := a.Parse(); err != nil {
			return nil, err
		}

		code, body, _ := a.Bytes()

		if code != 200 {
			return nil, &errors.Response{HttpCode: code, Message: string(body)}
		}

		res := new(response.MidtransTx)

		if err := json.Unmarshal(body, &res); err != nil {
			return nil, err
		}

		return res, nil
	})

	if err != nil {
		return nil, err
	}

	txRes, ok := res.(*response.MidtransTx)
	if !ok {
		return nil, &errors.Response{HttpCode: 500, Message: "failed to cast response to midtrans transaction response"}
	}

	return txRes, nil
}
