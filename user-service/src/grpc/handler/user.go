package handler

import (
	"context"
	"user-service/src/factory"
	"user-service/src/grpc/service"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/user"
)

type userImpl struct {
	service service.User
	pb.UnimplementedUserServiceServer
}

func NewUser(f *factory.Factory) pb.UserServiceServer {
	return &userImpl{
		service: service.NewUser(f),
	}
}

func (u *userImpl) CheckUserBlock(ctx context.Context, req *pb.CheckUserBlockReq) (*pb.CheckUserBlockRes, error) {
	return u.service.CheckUserBlock(ctx, req)
}
