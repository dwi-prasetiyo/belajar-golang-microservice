package main

import (
	"context"
	"log-service/env"
	"log-service/src/common/log"
	"log-service/src/consumer"
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

	rlc := consumer.NewKafka(ctx, env.Conf.RestfulLogConsumer)
	defer rlc.Close()

	go rlc.Run()

	rmqc := consumer.NewKafka(ctx, env.Conf.RabbitMQLogConsumer)
	defer rmqc.Close()

	go rmqc.Run()

	gqc := consumer.NewKafka(ctx, env.Conf.GrpcLogConsumer)
	defer gqc.Close()

	go gqc.Run()

	<-ctx.Done()
}
