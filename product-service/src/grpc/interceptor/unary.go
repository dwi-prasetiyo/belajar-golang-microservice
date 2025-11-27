package interceptor

import (
	"context"
	"encoding/base64"
	"fmt"
	"product-service/env"
	"product-service/src/common/constant"
	"product-service/src/common/dto/request"
	"product-service/src/common/errors"
	"product-service/src/common/log"
	"product-service/src/publisher"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Unary struct {
	grpcLogPublisher *publisher.Kafka
}

func NewUnary(gl *publisher.Kafka) *Unary {
	return &Unary{
		grpcLogPublisher: gl,
	}
}

func (u *Unary) Error(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	res, err := handler(ctx, req)

	if err != nil {
		log.Logger.Error(err.Error())

		if e, ok := err.(*errors.Response); ok {
			return nil, status.Errorf(e.GrpcCode, "%s", e.Message)
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return res, nil
}

func (u *Unary) Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error(fmt.Sprintf("%v", r))

			resp = nil
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()

	res, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *Unary) BasicAuthValidate(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("Authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Basic ") {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization type")
	}

	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid base64 authorization")
	}

	if string(decoded) != env.Conf.CurrentApp.GrpcBasicAuth {
		return nil, status.Error(codes.Unauthenticated, "invalid basic auth")
	}

	return handler(ctx, req)
}

func (u *Unary) Log(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	now := time.Now()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing metadata")
	}

	res, err := handler(ctx, req)

	var errMessage string
	if err != nil {
		errMessage = err.Error()
	}

	var requestID *string
	if v := md.Get("request_id"); len(v) > 0 {
		requestID = &v[0]
	}

	var userID *string
	if v := md.Get("user_id"); len(v) > 0 {
		userID = &v[0]
	}

	var appID *string
	if v := md.Get("app_id"); len(v) > 0 {
		appID = &v[0]
	}

	st, _ := status.FromError(err)

	u.grpcLogPublisher.Publish(request.GrpcLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		RequestID:  requestID,
		UserID:     userID,
		AppID:      appID,
		Method:     info.FullMethod,
		Body:       req,
		StatusCode: st.Code(),
		Latency:    float64(time.Since(now).Nanoseconds()) / 1e6,
		Error:      errMessage,
	})

	return res, err
}

func (u *Unary) AddBasicAuth(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}

	auth := base64.StdEncoding.EncodeToString([]byte(env.Conf.UserService.GrpcAuth))
	md.Append("Authorization", "Basic "+auth)

	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}

func (u *Unary) AddMetadata(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}

	if requestID, ok := ctx.Value(constant.RequestID).(string); ok {
		md.Append("request_id", requestID)
	}

	if userID, ok := ctx.Value(constant.UserID).(string); ok {
		md.Append("user_id", userID)
	}

	md.Append("app_id", "product-service")

	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}
