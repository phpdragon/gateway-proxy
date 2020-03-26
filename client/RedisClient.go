package client

import (
	"flag"
	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
	redisAddr = flag.String("redis_addr", "172.16.1.120:6379", "redis addr")
	redisPwd  = flag.String("redis_pwd", "", "redis pwd")
	redisDb   = flag.Int("redis_db", 0, "redis db")
)

func init()  {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: *redisPwd, // no password set
		DB:       *redisDb,  // use default DB
	})
}

func Redis() *redis.Client {
	return redisClient
}
