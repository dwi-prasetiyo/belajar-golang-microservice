package factory

import (
	"user-service/src/cache"
	"user-service/src/publisher"
	"user-service/src/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Factory struct {
	UserRepository         repository.User
	CredentialRepository   repository.Credential
	RefreshTokenRepository repository.RefreshToken
	UserBlockLogRepository repository.UserBlockLog
	CacheRepository        cache.Cache
	RabbitMQPublisher      *publisher.RabbitMQ
	RestfulLogPublisher    *publisher.Kafka
	GrpcLogPublisher       *publisher.Kafka
}

func New(db *gorm.DB, rc *redis.ClusterClient, pr *publisher.RabbitMQ, rl *publisher.Kafka, gl *publisher.Kafka) *Factory {
	return &Factory{
		UserRepository:         repository.NewUser(db),
		CredentialRepository:   repository.NewCredential(db),
		RefreshTokenRepository: repository.NewRefreshToken(db),
		UserBlockLogRepository: repository.NewUserBlockLog(db),
		CacheRepository:        cache.NewCache(rc),
		RabbitMQPublisher:      pr,
		RestfulLogPublisher:    rl,
		GrpcLogPublisher:       gl,
	}
}
