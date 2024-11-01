package services

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (s *URLService) ResolveURL(short string) (string, error) {
	url, err := s.repo.GetURL(s.ctx, short)
	if err == redis.Nil {
		return "", fmt.Errorf("URL not found")
	} else if err != nil {
		return "", err
	}

	_ = s.repo.IncrementVisitCount(s.ctx, s.ctx.ClientIP())
	return url, nil
}
