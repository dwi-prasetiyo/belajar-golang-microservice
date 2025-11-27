package factory

import (
	grpcclient "product-service/src/grpc/client"
	"product-service/src/publisher"
	"product-service/src/repository"

	"product-service/src/restful/client"

	"gorm.io/gorm"
)

type Factory struct {
	ProductRepository   repository.Product
	ImageKitClient      *client.ImageKit
	UserClient          grpcclient.User
	GrpcLogPublisher    *publisher.Kafka
	RestfulLogPublisher *publisher.Kafka
}

func New(db *gorm.DB, gl, rl *publisher.Kafka, uc grpcclient.User) *Factory {
	return &Factory{
		ProductRepository:   repository.NewProduct(db),
		ImageKitClient:      client.NewImageKit(),
		UserClient:          uc,
		GrpcLogPublisher:    gl,
		RestfulLogPublisher: rl,
	}
}
