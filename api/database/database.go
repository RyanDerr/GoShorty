package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient(dbNo int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNo,
	})
}
