package interceptor

import (
	"context"
	"encoding/base64"
	"fmt"
	"user-service/env"
	"user-service/src/common/dto/request"
	"user-service/src/common/errors"
	"user-service/src/common/log"
	"user-service/src/factory"
	"user-service/src/publisher"
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

func NewUnary(f *factory.Factory) *Unary {
	return &Unary{
		grpcLogPublisher: f.GrpcLogPublisher,
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
