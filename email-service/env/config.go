package env

import (
	"os"
)

type rabbitmq struct {
	DSN string
}

type oauth struct {
	GmailRefreshToken string
	GmailClientId     string
	GmailClientSecret string
}

type kafka struct {
	Addr1 string
	Addr2 string
	Addr3 string
}

type Config struct {
	RabbitMQ *rabbitmq
	Oauth    *oauth
	Kafka    *kafka
}

var Conf *Config

func Load() {
	if os.Getenv("MODE") == "NON-DEV" {
		LoadFromVault()
	}

	rabbitmqConf := new(rabbitmq)
	rabbitmqConf.DSN = os.Getenv("RABBITMQ_DSN")

	oauthConf := new(oauth)
	oauthConf.GmailRefreshToken = os.Getenv("OAUTH_GMAIL_REFRESH_TOKEN")
	oauthConf.GmailClientId = os.Getenv("OAUTH_GMAIL_CLIENT_ID")
	oauthConf.GmailClientSecret = os.Getenv("OAUTH_GMAIL_CLIENT_SECRET")

	kafkaConf := new(kafka)
	kafkaConf.Addr1 = os.Getenv("KAFKA_ADDR_1")
	kafkaConf.Addr2 = os.Getenv("KAFKA_ADDR_2")
	kafkaConf.Addr3 = os.Getenv("KAFKA_ADDR_3")

	Conf = &Config{
		RabbitMQ: rabbitmqConf,
		Oauth:    oauthConf,
		Kafka:    kafkaConf,
	}
}
