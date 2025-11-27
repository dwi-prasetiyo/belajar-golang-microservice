package service

import (
	"context"
	"fmt"
	"time"
	"user-service/src/cache"
	"user-service/src/common/log"
	"user-service/src/common/model"
	"user-service/src/factory"
	"user-service/src/repository"
)

type Profile interface{
	GetProfile(ctx context.Context, userID string) (*model.User, error) 
}

type profileImpl struct {
	cacheRepository cache.Cache
	userRepository  repository.User
}

func NewProfile(f *factory.Factory) Profile {
	return &profileImpl{
		cacheRepository: f.CacheRepository,
		userRepository:  f.UserRepository,
	}
}

func (s *profileImpl) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.cacheRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	user, err = s.userRepository.Find(ctx, "id = ?", []any{userID})
	if err != nil {
		return nil, err
	}

	if user != nil {
		if err := s.cacheRepository.Set(ctx, fmt.Sprintf("user:%s", userID), user, 30*time.Minute); err != nil {
			log.Logger.Error(err.Error())
		}
	}

	return user, nil
}
