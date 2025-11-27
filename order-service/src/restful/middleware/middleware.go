package middleware

import (
	"order-service/src/broker/publisher"
	"order-service/src/factory"
	"order-service/src/grpc/client"
)

type Middleware struct {
	restfulLogPublisher *publisher.Kafka
	userClient          client.User
}

func New(f *factory.Factory) *Middleware {
	return &Middleware{
		restfulLogPublisher: f.RestfulLogPublisher,
		userClient:          f.UserClient,
	}
}
