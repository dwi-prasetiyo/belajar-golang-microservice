package mockgrpcclient

import (
	"context"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
	"github.com/stretchr/testify/mock"
)

type Product struct {
	mock.Mock
}

func NewProduct() *Product {
	return &Product{
		Mock: mock.Mock{},
	}
}

func (m *Product) ReduceStocks(ctx context.Context, data *pb.ReduceStocksReq) error {
	ars := m.Called(ctx, data)
	return ars.Error(0)
}

func (m *Product) RollbackStocks(ctx context.Context, data *pb.RollbackStocksReq) error {
	ars := m.Called(ctx, data)
	return ars.Error(0)
}

func (m *Product) Close() {}
