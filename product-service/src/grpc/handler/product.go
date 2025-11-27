package handler

import (
	"context"
	"product-service/src/factory"
	"product-service/src/grpc/service"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productImpl struct {
	service service.Product
	pb.UnimplementedProductServiceServer
}

func NewProduct(f *factory.Factory) pb.ProductServiceServer {
	return &productImpl{
		service: service.NewProduct(f),
	}
}

func (p *productImpl) ReduceStocks(ctx context.Context, req *pb.ReduceStocksReq) (*emptypb.Empty, error) {
	err := p.service.ReduceStocks(ctx, req)
	return nil, err
}

func (p *productImpl) RollbackStocks(ctx context.Context, req *pb.RollbackStocksReq) (*emptypb.Empty, error) {
	err := p.service.RollbackStocks(ctx, req)
	return nil, err
}