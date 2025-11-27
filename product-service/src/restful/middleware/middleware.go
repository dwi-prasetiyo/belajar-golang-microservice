package middleware

import (
	"product-service/src/factory"
	"product-service/src/grpc/client"
	"product-service/src/publisher"
)

type Middleware struct {
	userClient          client.User
	restfulLogPublisher *publisher.Kafka
}

func New(f *factory.Factory) *Middleware {
	return &Middleware{
		userClient:          f.UserClient,
		restfulLogPublisher: f.RestfulLogPublisher,
	}
}
