package factory

import (
	"order-service/src/broker/publisher"
	grpcclient "order-service/src/grpc/client"
	"order-service/src/repository"
	"order-service/src/restful/client"

	"gorm.io/gorm"
)

type Factory struct {
	MidtransClient         client.Midtrans
	OrderRepository        repository.Order
	ProductClient          grpcclient.Product
	TxRepository           repository.TxBeginner
	ProductOrderRepository repository.ProductOrder
	RestfulLogPublisher    *publisher.Kafka
	UserClient             grpcclient.User
}

func New(db *gorm.DB, pc grpcclient.Product, rl *publisher.Kafka, uc grpcclient.User) *Factory {
	return &Factory{
		MidtransClient:         client.NewMidtrans(),
		OrderRepository:        repository.NewOrder(db),
		ProductClient:          pc,
		TxRepository:           repository.NewTxBeginner(db),
		ProductOrderRepository: repository.NewProductOrder(db),
		RestfulLogPublisher:    rl,
		UserClient:             uc,
	}
}
