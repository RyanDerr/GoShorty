package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

const (
	redisAddr = "REDIS_ADDRESS"
	redisPass = "REDIS_PASSWORD"
)

func CreateRedisClient(dbNo int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv(redisAddr),
		Password: os.Getenv(redisPass),
		DB:       dbNo,
	})
}
