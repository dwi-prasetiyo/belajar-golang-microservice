package mockrepository

import (
	"order-service/src/repository"

	"github.com/stretchr/testify/mock"
)

type Transactioner struct {
	Mock         mock.Mock
	Order        *Order
	ProductOrder *ProductOrder
}

func NewTransactioner() *Transactioner {
	return &Transactioner{
		Mock:         mock.Mock{},
		Order:        NewOrder(),
		ProductOrder: NewProductOrder(),
	}
}

func (m *Transactioner) Commit() error {
	args := m.Mock.Called()
	return args.Error(0)
}

func (m *Transactioner) Rollback() error {
	args := m.Mock.Called()
	return args.Error(0)
}

func (m *Transactioner) OrderRepository() repository.Order {
	return m.Order
}

func (m *Transactioner) ProductOrderRepository() repository.ProductOrder {
	return m.ProductOrder
}
