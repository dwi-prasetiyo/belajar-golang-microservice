package mockrepository

import (
	"context"
	"order-service/src/common/model"

	"github.com/stretchr/testify/mock"
)

type Order struct {
	mock.Mock
}

func NewOrder() *Order {
	return &Order{
		Mock: mock.Mock{},
	}
}

func (m *Order) CreateOrder(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *Order) Update(ctx context.Context, whereClause string, args []any, data map[string]any) error {
	arguments := m.Called(ctx, whereClause, args, data)
	return arguments.Error(0)
}

func (m *Order) Find(ctx context.Context, whereClause string, args []any) (*model.Order, error) {
	arguments := m.Called(ctx, whereClause, args)
	if arguments.Get(0) != nil {
		return arguments.Get(0).(*model.Order), arguments.Error(1)
	}
	return nil, arguments.Error(1)
}
