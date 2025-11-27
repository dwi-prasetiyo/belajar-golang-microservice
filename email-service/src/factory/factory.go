package factory

import (
	"email-service/src/broker/publisher"
	"email-service/src/common/pkg/oauth"

	"google.golang.org/api/gmail/v1"
)

type Factory struct {
	GmailService         *gmail.Service
	RabbitMQLogPublisher *publisher.Kafka
}

func NewFactory(rl *publisher.Kafka) *Factory {
	return &Factory{
		GmailService:         oauth.NewGmailService(),
		RabbitMQLogPublisher: rl,
	}
}
