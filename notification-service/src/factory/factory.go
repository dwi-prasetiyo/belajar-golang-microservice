package factory

import "notification-service/src/publisher"

type Factory struct {
	RabbitMQPublisher *publisher.RabbitMQ
}

func New(pr *publisher.RabbitMQ) *Factory {
	return &Factory{
		RabbitMQPublisher: pr,
	}
}
