package repository

import (
	"context"
	"order-service/src/common/model"

	"gorm.io/gorm"
)

type Order interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, whereClause string, args []any, data map[string]any) error
	Find(ctx context.Context, whereClause string, args []any) (*model.Order, error)
}

type orderImpl struct {
	db *gorm.DB
}

func NewOrder(db *gorm.DB) Order {
	return &orderImpl{db: db}
}

func (r *orderImpl) CreateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderImpl) Update(ctx context.Context, whereClause string, args []any, data map[string]any) error {
	return r.db.WithContext(ctx).Model(model.Order{}).Where(whereClause, args...).Updates(data).Error
}

func (r *orderImpl) Find(ctx context.Context, whereClause string, args []any) (*model.Order, error) {
	var order model.Order

	err := r.db.WithContext(ctx).Model(model.Order{}).Where(whereClause, args...).First(&order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if order.ID == "" {
		return nil, nil
	}

	return &order, nil
}
