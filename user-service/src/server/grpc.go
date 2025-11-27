package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"user-service/env"
	"user-service/src/common/log"
	"user-service/src/factory"
	"user-service/src/grpc/handler"
	"user-service/src/grpc/interceptor"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Grpc struct {
	server      *grpc.Server
	userHandler pb.UserServiceServer
}

func NewGrpc(f *factory.Factory, tlsConf *tls.Config) *Grpc {
	userHandler := handler.NewUser(f)

	i := interceptor.NewUnary(f)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			i.Log,
			// i.BasicAuthValidate,
			i.Recovery,
			i.Error,
		),
		grpc.Creds(credentials.NewTLS(tlsConf)),
	)

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	reflection.Register(grpcServer)

	return &Grpc{
		server:      grpcServer,
		userHandler: userHandler,
	}
}

func (g *Grpc) Run() {
	listener, err := net.Listen("tcp", env.Conf.CurrentApp.GrpcAddr)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	log.Logger.Info(fmt.Sprintf("user grpc server run in %s", env.Conf.CurrentApp.GrpcAddr))

	if err := g.server.Serve(listener); err != nil {
		log.Logger.Fatal(err.Error())
	}
}

func (g *Grpc) Stop() {
	g.server.Stop()
}
