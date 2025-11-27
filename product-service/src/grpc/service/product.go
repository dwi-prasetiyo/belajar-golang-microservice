package service

import (
	"context"
	"product-service/src/factory"
	"product-service/src/repository"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
)

type Product interface {
	ReduceStocks(ctx context.Context, req *pb.ReduceStocksReq) error
	RollbackStocks(ctx context.Context, req *pb.RollbackStocksReq) error
}

type productImpl struct {
	productRepository repository.Product
}

func NewProduct(f *factory.Factory) Product {
	return &productImpl{
		productRepository: f.ProductRepository,
	}
}

func (p *productImpl) ReduceStocks(ctx context.Context, req *pb.ReduceStocksReq) error {
	return p.productRepository.ReduceStocks(ctx, req.ProductOrders)
}

func (p *productImpl) RollbackStocks(ctx context.Context, req *pb.RollbackStocksReq) error {
	return p.productRepository.RollbackStocks(ctx, req.ProductOrders)
}