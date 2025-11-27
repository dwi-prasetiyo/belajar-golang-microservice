package repository

import (
	"context"
	"order-service/src/common/model"

	"gorm.io/gorm"
)

type ProductOrder interface {
	Create(ctx context.Context, productOrder []*model.ProductOrder) error
	Find(ctx context.Context, orderID string) ([]*model.ProductOrder, error) 
}

type productOrderImpl struct {
	db *gorm.DB
}

func NewProductOrder(db *gorm.DB) ProductOrder {
	return &productOrderImpl{db: db}
}

func (p *productOrderImpl) Create(ctx context.Context, productOrder []*model.ProductOrder) error {
	return p.db.WithContext(ctx).Create(productOrder).Error
}

func (p *productOrderImpl) Find(ctx context.Context, orderID string) ([]*model.ProductOrder, error) {
	var res []*model.ProductOrder

	err := p.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&res).Error
	return res, err
}
