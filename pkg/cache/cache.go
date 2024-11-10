package cache

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	redisAddr = "REDIS_ADDRESS"
	redisPass = "REDIS_PASSWORD"
)

func CreateRedisClient(ctx *gin.Context, dbNo int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv(redisAddr),
		DB:   dbNo,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		return nil, err
	}

	return client, nil
}
