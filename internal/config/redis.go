package config

import (
	"github.com/go-redis/redis"
	"runtime"
)

// 初始化
var redisClient *redis.Client

func init() {
	redisConfig := GetRedisConfig()
	poolSize := redisConfig.PoolSize
	if poolSize <= 0 {
		poolSize = runtime.NumCPU()
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.Db,       // use default DB
		PoolSize: poolSize,
	})
}

func Redis() *redis.Client {
	return redisClient
}
