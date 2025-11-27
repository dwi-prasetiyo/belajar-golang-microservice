package main

import (
	"context"
	"email-service/env"
	"email-service/src/broker/consumer"
	"email-service/src/broker/publisher"
	"email-service/src/common/log"
	"email-service/src/factory"

	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func handleCloseApp(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	handleCloseApp(cancel)

	defer log.Logger.Sync()

	godotenv.Load()
	env.Load()

	rl := publisher.NewKafka("rabbitmq-log")
	defer rl.Close()

	f := factory.NewFactory(rl)
	rabbitMQConsumer := consumer.NewRabbitMQ(ctx, f)
	defer rabbitMQConsumer.Close()

	go rabbitMQConsumer.Otp()

	<-ctx.Done()
}
