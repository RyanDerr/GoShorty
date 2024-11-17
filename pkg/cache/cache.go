package cache

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

const (
	redisAddr = "REDIS_ADDRESS"
)

func CreateRedisClient(dbNo int) (*redis.Client, error) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv(redisAddr),
		DB:   dbNo,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Error pinging Redis: %v", err)
		return nil, err
	}

	return client, nil
}
