package model

import (
	"log"

	"git.oschina.net/cnjack/novel-spider/config"
	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
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
