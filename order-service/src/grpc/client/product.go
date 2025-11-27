package client

import (
	"context"
	"crypto/tls"
	"order-service/env"
	"order-service/src/common/log"
	"order-service/src/common/pkg/cbreaker"
	"order-service/src/grpc/interceptor"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Product interface {
	ReduceStocks(ctx context.Context, data *pb.ReduceStocksReq) error
	RollbackStocks(ctx context.Context, data *pb.RollbackStocksReq) error
	Close()
}

type productImpl struct {
	client pb.ProductServiceClient
	conn   *grpc.ClientConn
	cb     *gobreaker.CircuitBreaker[any]
}

func NewProduct(i *interceptor.Unary, tlsConf *tls.Config) Product {
	cb := cbreaker.NewGrpc("product-grpc")

	var opts []grpc.DialOption
	opts = append(
		opts,
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)),
		grpc.WithChainUnaryInterceptor(
			i.AddMetadata,
			// i.AddBasicAuth,
		),
	)

	conn, err := grpc.NewClient(env.Conf.ProductService.GrpcAddr, opts...)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	return &productImpl{
		client: pb.NewProductServiceClient(conn),
		conn:   conn,
		cb:     cb,
	}
}

func (p *productImpl) ReduceStocks(ctx context.Context, data *pb.ReduceStocksReq) error {
	_, err := p.cb.Execute(func() (any, error) {
		_, err := p.client.ReduceStocks(ctx, data)
		return nil, err
	})

	return err
}

func (p *productImpl) RollbackStocks(ctx context.Context, data *pb.RollbackStocksReq) error {
	_, err := p.cb.Execute(func() (any, error) {
		_, err := p.client.RollbackStocks(ctx, data)
		return nil, err
	})

	return err
}

func (p *productImpl) Close() {
	if err := p.conn.Close(); err != nil {
		log.Logger.Error(err.Error())
	}
}
