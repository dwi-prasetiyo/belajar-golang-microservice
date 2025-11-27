package mockrepository

import (
	"context"
	"order-service/src/common/model"

	"github.com/stretchr/testify/mock"
)

type ProductOrder struct {
	mock.Mock
}

func NewProductOrder() *ProductOrder {
	return &ProductOrder{
		Mock: mock.Mock{},
	}
}

func (m *ProductOrder) Create(ctx context.Context, productOrder []*model.ProductOrder) error {
	args := m.Called(ctx, productOrder)
	return args.Error(0)
}

func (m *ProductOrder) Find(ctx context.Context, orderID string) ([]*model.ProductOrder, error) {
	arguments := m.Called(ctx, orderID)
	if arguments.Get(0) != nil {
		return arguments.Get(0).([]*model.ProductOrder), arguments.Error(1)
	}
	return nil, arguments.Error(1)
}
