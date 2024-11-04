package database

import (
	"os"

	"github.com/RyanDerr/GoShorty/api/config"
	"github.com/redis/go-redis/v9"
)

func CreateRedisClient(dbNo int) *redis.Client {
	var rdb *redis.Client
	if config.IsLocal() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDRESS"),
			Password: os.Getenv("DB_PASSWORD"),
			DB:       dbNo,
		})
	} else {
		// TODO: Implement Elasticache connection
		rdb = nil
	}

	return rdb
}
