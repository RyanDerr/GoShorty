package repositories

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type URLRepository struct {
	client *redis.Client
}

func NewURLRepository(client *redis.Client) *URLRepository {
	return &URLRepository{client: client}
}

func (r *URLRepository) CheckURLInUse(ctx *gin.Context, id string) error {
	val, _ := r.client.Get(ctx, id).Result()
	if val != "" {
		return fmt.Errorf("url already in use")
	}
	return nil
}

func (r *URLRepository) SaveURL(ctx *gin.Context, id string, url string, expiration time.Duration) error {
	return r.client.Set(ctx, id, url, expiration).Err()
}

func (r *URLRepository) GetURL(ctx *gin.Context, short string) (string, error) {
	return r.client.Get(ctx, short).Result()
}

func (r *URLRepository) IncrementVisitCount(ctx *gin.Context, clientIP string) error {
	return r.client.Incr(ctx, clientIP).Err()
}
