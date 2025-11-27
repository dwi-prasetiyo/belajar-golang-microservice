package main

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"
	"product-service/database"
	"product-service/env"
	"product-service/src/common/log"
	"product-service/src/common/util"
	"product-service/src/factory"
	"product-service/src/grpc/client"
	"product-service/src/grpc/interceptor"
	"product-service/src/publisher"
	"product-service/src/server"
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

	serverCert, err := tls.LoadX509KeyPair("../certs/server.crt", "../certs/server.key")
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	clientCert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	serverTlsConf := util.CreateServerTlsConf(caCert, serverCert)
	clientTlsConf := util.CreateClientTlsConf(caCert, clientCert)

	db := database.NewPostgreSQL()
	defer database.ClosePostgreSQL(db)

	gl := publisher.NewKafka("grpc-log")
	defer gl.Close()

	rl := publisher.NewKafka("restful-log")
	defer rl.Close()

	i := interceptor.NewUnary(gl)

	uc := client.NewUser(i, clientTlsConf)

	f := factory.New(db, gl, rl, uc)

	restfulServer := server.NewRestful(f)
	defer restfulServer.Stop()

	go restfulServer.Start()

	grpcServer := server.NewGrpc(f, i, serverTlsConf)
	defer grpcServer.Stop()

	go grpcServer.Run()

	<-ctx.Done()
}
