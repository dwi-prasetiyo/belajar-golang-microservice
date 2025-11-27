package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-service/src/common/dto/request"
	"user-service/src/common/dto/response"
	"user-service/src/common/model"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, expired time.Duration) error
	GetSendOtp(ctx context.Context, key string) (*request.SendOtp, error)
	GetRegister(ctx context.Context, key string) (*request.Register, error)
	GetUser(ctx context.Context, userID string) (*model.User, error)
	GetUserBlockInfo(ctx context.Context, userID string) (*response.UserBlockInfo, error)
	Delete(ctx context.Context, key string) error
}

type cacheImpl struct {
	redis *redis.ClusterClient
}

func NewCache(r *redis.ClusterClient) Cache {
	return &cacheImpl{redis: r}
}

func (c *cacheImpl) Set(ctx context.Context, key string, value any, expired time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.redis.SetEx(ctx, key, jsonData, expired).Err()
}

func (c *cacheImpl) GetSendOtp(ctx context.Context, key string) (*request.SendOtp, error) {
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if data == "" {
		return nil, nil
	}

	var otp request.SendOtp
	if err := json.Unmarshal([]byte(data), &otp); err != nil {
		return nil, err
	}

	return &otp, nil
}

func (c *cacheImpl) GetRegister(ctx context.Context, key string) (*request.Register, error) {
	data, err := c.redis.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if data == "" {
		return nil, nil
	}

	var register request.Register
	if err := json.Unmarshal([]byte(data), &register); err != nil {
		return nil, err
	}

	return &register, nil
}

func (c *cacheImpl) GetUser(ctx context.Context, userID string) (*model.User, error) {
	data, err := c.redis.Get(ctx, fmt.Sprintf("user:%s", userID)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if data == "" {
		return nil, nil
	}

	var user model.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *cacheImpl) GetUserBlockInfo(ctx context.Context, userID string) (*response.UserBlockInfo, error) {
	data, err := c.redis.Get(ctx, fmt.Sprintf("user_block_info:%s", userID)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	
	if data == "" {
		return nil, nil
	}

	var userBlockInfo response.UserBlockInfo
	if err := json.Unmarshal([]byte(data), &userBlockInfo); err != nil {
		return nil, err
	}

	return &userBlockInfo, nil
}

func (c *cacheImpl) Delete(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}
