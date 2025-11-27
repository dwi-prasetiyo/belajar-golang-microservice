package database

import (
	"user-service/env"
	"user-service/src/common/log"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.ClusterClient {
	redisCluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			env.Conf.Redis.AddrNode1,
			env.Conf.Redis.AddrNode2,
			env.Conf.Redis.AddrNode3,
			env.Conf.Redis.AddrNode4,
			env.Conf.Redis.AddrNode5,
			env.Conf.Redis.AddrNode6,
		},
		Password: env.Conf.Redis.Password,
	})

	return redisCluster
}

func CloseRedis(rc *redis.ClusterClient) {
	if err := rc.Close(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
