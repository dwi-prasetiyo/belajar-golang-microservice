package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"product-service/env"
	"product-service/src/common/log"
	"product-service/src/factory"
	"product-service/src/grpc/handler"
	"product-service/src/grpc/interceptor"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Grpc struct {
	server         *grpc.Server
	productHandler pb.ProductServiceServer
}

func NewGrpc(f *factory.Factory, i *interceptor.Unary, tlsConf *tls.Config) *Grpc {
	productHandler := handler.NewProduct(f)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			i.Log,
			// i.BasicAuthValidate,
			i.Recovery,
			i.Error,
		),
		grpc.Creds(credentials.NewTLS(tlsConf)),
	)

	pb.RegisterProductServiceServer(grpcServer, productHandler)

	return &Grpc{
		server:         grpcServer,
		productHandler: productHandler,
	}
}

func (g *Grpc) Run() {
	listener, err := net.Listen("tcp", env.Conf.CurrentApp.GrpcAddr)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	log.Logger.Info(fmt.Sprintf("product grpc server run in %s", env.Conf.CurrentApp.GrpcAddr))

	if err := g.server.Serve(listener); err != nil {
		log.Logger.Fatal(err.Error())
	}
}

func (g *Grpc) Stop() {
	g.server.Stop()
}
