package service

import (
	"context"
	"order-service/src/common/constant"
	"order-service/src/common/dto/request"
	"order-service/src/common/dto/response"
	v "order-service/src/common/pkg/validator"
	"order-service/src/factory"
	grpcclient "order-service/src/grpc/client"
	"order-service/src/repository"

	"order-service/src/restful/client"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Order interface {
	CreateOrder(c *fiber.Ctx, order *request.CreateOrder) (*response.MidtransTx, error)
}

type orderImpl struct {
	midtransClient  client.Midtrans
	orderRepository repository.Order
	productClient   grpcclient.Product
	txRepository    repository.TxBeginner
}

func NewOrder(f *factory.Factory) Order {
	return &orderImpl{
		midtransClient:  f.MidtransClient,
		orderRepository: f.OrderRepository,
		productClient:   f.ProductClient,
		txRepository:    f.TxRepository,
	}
}

func (s *orderImpl) CreateOrder(c *fiber.Ctx, order *request.CreateOrder) (*response.MidtransTx, error) {
	ctx := c.UserContext()

	if err := v.Validate.Struct(order); err != nil {
		return nil, err
	}

	orderID, err := gonanoid.New()
	if err != nil {
		return nil, err
	}

	order.Order.ID = orderID

	midtransTxReq := &request.MidtransTransaction{
		TransactionDetails: request.MidtransTransactionDetails{
			OrderID:     orderID,
			GrossAmount: order.Order.GrossAmount,
		},
		Metadata: request.MidtransTransactionMetadata{
			OriginalOrderID: orderID,
		},
	}

	res, err := s.midtransClient.CreateTransaction(ctx, midtransTxReq)
	if err != nil {
		return nil, err
	}

	order.Order.PaymentURL = res.PaymentURL
	order.Order.Status = "PENDING_PAYMENT"

	tx, err := s.txRepository.Begin()
	if err != nil {
		return nil, err
	}

	for _, p := range order.ProductOrders {
		p.OrderID = orderID
	}

	err = tx.OrderRepository().CreateOrder(ctx, order.Order)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.ProductOrderRepository().Create(ctx, order.ProductOrders)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var req []*pb.ProductOrder

	for _, v := range order.ProductOrders {
		req = append(req, &pb.ProductOrder{
			ProductId: uint32(v.ProductID),
			Quantity:  uint32(v.Quantity),
		})
	}

	ctx = context.WithValue(ctx, constant.RequestID, c.Locals("request_id"))
	ctx = context.WithValue(ctx, constant.UserID, c.Locals("user_id"))

	err = s.productClient.ReduceStocks(ctx, &pb.ReduceStocksReq{ProductOrders: req})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return res, nil
}
