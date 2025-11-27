package client

import (
	"context"
	"crypto/tls"
	"order-service/env"
	"order-service/src/common/log"
	"order-service/src/common/pkg/cbreaker"
	"order-service/src/grpc/interceptor"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/user"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type User interface {
	CheckUserBlock(ctx context.Context, req *pb.CheckUserBlockReq) (*pb.CheckUserBlockRes, error)
	Close()
}

type userImpl struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
	cb     *gobreaker.CircuitBreaker[any]
}

func NewUser(i *interceptor.Unary, tlsConf *tls.Config) User {
	cb := cbreaker.NewGrpc("user-grpc")

	var opts []grpc.DialOption
	opts = append(
		opts,
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)),
		grpc.WithChainUnaryInterceptor(
			i.AddMetadata,
			i.AddBasicAuth,
		),
	)

	conn, err := grpc.NewClient(env.Conf.UserService.GrpcAddr, opts...)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	return &userImpl{
		client: pb.NewUserServiceClient(conn),
		conn:   conn,
		cb:     cb,
	}
}

func (u *userImpl) CheckUserBlock(ctx context.Context, req *pb.CheckUserBlockReq) (*pb.CheckUserBlockRes, error) {
	resp, err := u.cb.Execute(func() (any, error) {
		return u.client.CheckUserBlock(ctx, req)
	})

	if err != nil {
		return nil, err
	}

	result, ok := resp.(*pb.CheckUserBlockRes)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid response type")
	}

	return result, nil
}

func (u *userImpl) Close() {
	if err := u.conn.Close(); err != nil {
		log.Logger.Error(err.Error())
	}
}
