package main

import (
	"context"
	"dummy/producer"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func addBasicAuth(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}

	auth := base64.StdEncoding.EncodeToString([]byte("auth:rahasia"))
	md.Append("Authorization", "Basic "+auth)

	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}

func deadlockSimmulation() {
	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(addBasicAuth),
	)

	conn, err := grpc.NewClient("localhost:9091", opts...)
	if err != nil {
		panic(err.Error())
	}

	client := pb.NewProductServiceClient(conn)

	wg := sync.WaitGroup{}
	wg.Add(4)

	data := []*pb.ReduceStocksReq{
		{
			ProductOrders: []*pb.ProductOrder{
				{
					ProductId: 1,
					Quantity:  1,
				},
				{
					ProductId: 2,
					Quantity:  1,
				},
				{
					ProductId: 3,
					Quantity:  1,
				},
				{
					ProductId: 4,
					Quantity:  1,
				},
			},
		},
		{
			ProductOrders: []*pb.ProductOrder{
				{
					ProductId: 4,
					Quantity:  1,
				},
				{
					ProductId: 3,
					Quantity:  1,
				},
				{
					ProductId: 2,
					Quantity:  1,
				},
				{
					ProductId: 1,
					Quantity:  1,
				},
			},
		},
		{
			ProductOrders: []*pb.ProductOrder{
				{
					ProductId: 2,
					Quantity:  1,
				},
				{
					ProductId: 3,
					Quantity:  1,
				},
				{
					ProductId: 4,
					Quantity:  1,
				},
				{
					ProductId: 1,
					Quantity:  1,
				},
			},
		},
		{
			ProductOrders: []*pb.ProductOrder{
				{
					ProductId: 3,
					Quantity:  1,
				},
				{
					ProductId: 1,
					Quantity:  1,
				},
				{
					ProductId: 2,
					Quantity:  1,
				},
				{
					ProductId: 4,
					Quantity:  1,
				},
			},
		},
	}

	for i := 0; i < 4; i++ {
		go func(i int) {
			defer wg.Done()
			_, err := client.ReduceStocks(context.Background(), data[i])
			if err != nil {
				panic(err.Error())
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("done")
}

func singglePointOfFailureSimulation() {
	ctx := context.Background()

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// })

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:5371",
			"localhost:5372",
			"localhost:5373",
			"localhost:5374",
			"localhost:5375",
			"localhost:5376",
		},
		Password: "rahasia",
	})

	defer rdb.Close()

	for {
		val, err := rdb.Get(ctx, "session:user123").Result()
		if err == redis.Nil {
			fmt.Println("session belum ada")
		} else if err != nil {
			log.Println("Redis error:", err)
		} else {
			fmt.Println("session:", val)
		}
		time.Sleep(2 * time.Second)
	}
}

func flushall() {

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:5371",
			"localhost:5372",
			"localhost:5373",
			"localhost:5374",
			"localhost:5375",
			"localhost:5376",
		},
		Password: "rahasia",
	})

	defer rdb.Close()

	err := rdb.ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
		return client.FlushAll(ctx).Err()
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	// deadlockSimmulation()
	// singglePointOfFailureSimulation()
	// flushall()
	producer.SendRealLogTraffic()
}
