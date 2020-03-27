package client

import (
	"../core"
	"github.com/go-redis/redis"
)

//初始化
var redisClient *redis.Client

func init()  {
	redisConfig := core.GetRedisConfig()
	redisClient = redis.NewClient(&redis.Options{
		Addr:    redisConfig.Host,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.Db,  // use default DB
	})
}

func Redis() *redis.Client {
	return redisClient
}
