package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRESS"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       dbNo,
	})

	return rdb
}
