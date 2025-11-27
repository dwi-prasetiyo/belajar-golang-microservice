package service

import (
	"context"
	"encoding/json"
	"order-service/src/common/dto/request"
	"order-service/src/factory"
	"order-service/src/grpc/client"
	"order-service/src/repository"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
)

type RabbitMQ struct {
	orderRepository        repository.Order
	productOrderRepository repository.ProductOrder
	productClient          client.Product
}

func NewRabbitMQ(f *factory.Factory) *RabbitMQ {
	return &RabbitMQ{
		orderRepository:        f.OrderRepository,
		productOrderRepository: f.ProductOrderRepository,
		productClient:          f.ProductClient,
	}
}

func (s *RabbitMQ) HandleMidtransTxNotif(ctx context.Context, data []byte) error {
	req := new(request.MidtransTxNotif)

	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	if req.TransactionStatus == "capture" && req.FraudStatus == "accept" || req.TransactionStatus == "settlement" {
		err := s.orderRepository.Update(ctx, "id = ?", []any{req.OrderID}, map[string]any{"status": "PAID", "payment_method": req.PaymentType})
		if err != nil {
			return err
		}
	}

	if req.TransactionStatus == "cancel" {
		order, err := s.orderRepository.Find(ctx, "id = ?", []any{req.OrderID})
		if err != nil {
			return err
		}

		if order == nil || order.Status == "CANCELLED" {
			return nil
		}

		products, err := s.productOrderRepository.Find(ctx, req.OrderID)
		if err != nil {
			return err
		}

		var productOrders []*pb.ProductOrder

		for _, product := range products {
			productOrders = append(productOrders, &pb.ProductOrder{
				ProductId: uint32(product.ProductID),
				Quantity:  uint32(product.Quantity),
			})
		}

		err = s.productClient.RollbackStocks(ctx, &pb.RollbackStocksReq{ProductOrders: productOrders})
		if err != nil {
			return err
		}

		err = s.orderRepository.Update(ctx, "id = ?", []any{req.OrderID}, map[string]any{"status": "CANCELLED"})
		if err != nil {
			return err
		}
	}

	if req.TransactionStatus == "deny" || req.TransactionStatus == "expire" {
		order, err := s.orderRepository.Find(ctx, "id = ?", []any{req.OrderID})
		if err != nil {
			return err
		}

		if order == nil || order.Status == "FAILED" {
			return nil
		}

		products, err := s.productOrderRepository.Find(ctx, req.OrderID)
		if err != nil {
			return err
		}

		var productOrders []*pb.ProductOrder

		for _, product := range products {
			productOrders = append(productOrders, &pb.ProductOrder{
				ProductId: uint32(product.ProductID),
				Quantity:  uint32(product.Quantity),
			})
		}

		err = s.productClient.RollbackStocks(ctx, &pb.RollbackStocksReq{ProductOrders: productOrders})
		if err != nil {
			return err
		}

		err = s.orderRepository.Update(ctx, "id = ?", []any{req.OrderID}, map[string]any{"status": "FAILED"})
		if err != nil {
			return err
		}
	}

	return nil
}
