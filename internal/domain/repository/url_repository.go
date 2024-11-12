package repository

import (
	"log"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UrlRepository struct {
	client *redis.Client
}

type IUrlRepository interface {
	CheckShortInUse(*gin.Context, string) (bool, error)
	SaveUrl(*gin.Context, *entity.ShortenUrl) (*entity.ShortenUrl, error)
	GetUrl(*gin.Context, string) (string, error)
}

func NewUrlRepository(redisClient *redis.Client) *UrlRepository {
	return &UrlRepository{
		client: redisClient,
	}
}

func (r *UrlRepository) CheckShortInUse(ctx *gin.Context, short string) (bool, error) {
	val, err := r.client.Get(ctx, short).Result()

	if err == redis.Nil {
		log.Printf("Key %s does not exist", short)
		return false, nil
	}

	if err != nil {
		log.Printf("Error checking if key %s exists: %v", short, err)
		return false, err
	}

	if val != "" {
		return true, nil
	}

	return false, nil
}

func (r *UrlRepository) SaveUrl(ctx *gin.Context, short *entity.ShortenUrl) (*entity.ShortenUrl, error) {
	err := r.client.Set(ctx, short.Short, short.BaseUrl, short.Expiration).Err()

	if err != nil {
		log.Printf("Error saving URL: %v", err)
		return nil, err
	}

	return short, nil
}

func (r *UrlRepository) GetUrl(ctx *gin.Context, short string) (string, error) {
	return r.client.Get(ctx, short).Result()
}
