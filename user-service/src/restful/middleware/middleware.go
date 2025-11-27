package middleware

import (
	"user-service/src/cache"
	"user-service/src/factory"
	"user-service/src/publisher"
	"user-service/src/repository"
)

type Middleware struct {
	restfulLogPublisher    *publisher.Kafka
	cacheRepository        cache.Cache
	userBlockLogRepository repository.UserBlockLog
}

func New(f *factory.Factory) *Middleware {
	return &Middleware{
		restfulLogPublisher:    f.RestfulLogPublisher,
		cacheRepository:        f.CacheRepository,
		userBlockLogRepository: f.UserBlockLogRepository,
	}
}
