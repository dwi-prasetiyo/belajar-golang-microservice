package main

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"
	"syscall"
	"user-service/database"
	"user-service/env"
	"user-service/src/common/log"
	"user-service/src/common/util"
	"user-service/src/factory"
	"user-service/src/publisher"
	"user-service/src/server"

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

	serverCert, err := tls.LoadX509KeyPair("../certs/server.crt", "../certs/server.key")
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	serverTlsConf := util.CreateServerTlsConf(caCert, serverCert)

	db := database.NewPostgreSQL()
	defer database.ClosePostgreSQL(db)

	rc := database.NewRedis()
	defer database.CloseRedis(rc)

	pr := publisher.NewRabbitMQ()
	defer pr.Close()

	rl := publisher.NewKafka("restful-log")
	defer rl.Close()

	gl := publisher.NewKafka("grpc-log")
	defer gl.Close()

	f := factory.New(db, rc, pr, rl, gl)

	restfulServer := server.NewRestful(f)
	defer restfulServer.Stop()

	go restfulServer.Start()

	grpcServer := server.NewGrpc(f, serverTlsConf)
	defer grpcServer.Stop()

	go grpcServer.Run()

	<-ctx.Done()
}
