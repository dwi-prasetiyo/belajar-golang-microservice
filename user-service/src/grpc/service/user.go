package service

import (
	"context"
	"fmt"
	"time"
	"user-service/src/cache"
	"user-service/src/common/dto/response"
	"user-service/src/factory"
	"user-service/src/repository"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/user"
)

type User interface {
	CheckUserBlock(ctx context.Context, req *pb.CheckUserBlockReq) (*pb.CheckUserBlockRes, error)
}

type userImpl struct {
	userBlockLogRepository repository.UserBlockLog
	cacheRepository        cache.Cache
}

func NewUser(f *factory.Factory) User {
	return &userImpl{
		userBlockLogRepository: f.UserBlockLogRepository,
		cacheRepository:        f.CacheRepository,
	}
}

func (u *userImpl) CheckUserBlock(ctx context.Context, req *pb.CheckUserBlockReq) (*pb.CheckUserBlockRes, error) {

	userBlockInfo, err := u.cacheRepository.GetUserBlockInfo(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	if userBlockInfo != nil {
		return &pb.CheckUserBlockRes{IsBlocked: userBlockInfo.IsBlocked, Reason: userBlockInfo.Reason}, nil
	}

	userBlockLog, err := u.userBlockLogRepository.Find(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var isBlocked bool
	var reason string

	if userBlockLog != nil {
		isBlocked = true
		reason = userBlockLog.Reason
	}

	err = u.cacheRepository.Set(ctx, fmt.Sprintf("user_block_info:%s", req.UserId), &response.UserBlockInfo{IsBlocked: isBlocked, Reason: reason}, 30*time.Minute)
	if err != nil {
		return nil, err
	}

	return &pb.CheckUserBlockRes{IsBlocked: isBlocked, Reason: reason}, nil
}
