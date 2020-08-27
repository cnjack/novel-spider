package repository

import (
	"log"

	"spider/internal/config"

	"github.com/go-redis/redis"
)

var client *redis.Client

func InitRedis() {
	redisConf := config.GetConfig().RedisConfig
	client = redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password,
		DB:       redisConf.DB,
	})
	if err := client.Ping().Err(); err != nil {
		log.Panic(err)
	}
}

func MustGetRedisClient() *redis.Client {
	return client
}
