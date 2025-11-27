package main

import (
	"context"
	"crypto/tls"
	"order-service/database"
	"order-service/env"
	"order-service/src/broker/consumer"
	"order-service/src/broker/publisher"
	"order-service/src/common/log"
	"order-service/src/common/util"
	"order-service/src/factory"
	"order-service/src/grpc/client"
	"order-service/src/grpc/interceptor"
	"order-service/src/server"
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

	caCert, err := os.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	clientCert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	clientTlsConf := util.CreateClientTlsConf(caCert, clientCert)

	db := database.NewPostgreSQL()
	defer database.ClosePostgreSQL(db)

	i := interceptor.NewUnary()

	productClient := client.NewProduct(i, clientTlsConf)
	defer productClient.Close()

	rl := publisher.NewKafka("restful-log")
	defer rl.Close()

	uc := client.NewUser(i, clientTlsConf)
	defer uc.Close()

	f := factory.New(db, productClient, rl, uc)

	restfulServer := server.NewRestful(f)
	defer restfulServer.Stop()

	go restfulServer.Start()

	rabbitMQConsumer := consumer.NewRabbitMQ(ctx, f)
	defer rabbitMQConsumer.Close()

	go rabbitMQConsumer.MidtransTx()

	<-ctx.Done()
}
