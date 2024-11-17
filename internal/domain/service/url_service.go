package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/internal/domain/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ShortenUrlService struct {
	urlRepo repository.IUrlRepository
}

type IUrlService interface {
	ShortenUrl(*gin.Context, *entity.ShortenUrl) (*entity.ShortenUrl, int, error)
	ResolveUrl(*gin.Context, string) (string, int, error)
}

func NewShortenUrlService(urlRepo repository.IUrlRepository) *ShortenUrlService {
	return &ShortenUrlService{
		urlRepo: urlRepo,
	}
}

func (s *ShortenUrlService) ShortenUrl(ctx *gin.Context, short *entity.ShortenUrl) (*entity.ShortenUrl, int, error) {
	res, err := s.urlRepo.CheckShortInUse(ctx, short.Short)

	if err != nil {
		log.Printf("Error checking if short is in use: %s", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	if res {
		log.Printf("Short is already in use: %s", short.Short)
		return nil, http.StatusConflict, fmt.Errorf("short is already in use: %s", short.Short)
	}

	response, err := s.urlRepo.SaveUrl(ctx, short)

	if err != nil {
		log.Printf("Error saving URL: %s", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return response, http.StatusCreated, nil

}

func (s *ShortenUrlService) ResolveUrl(ctx *gin.Context, shortUrl string) (string, int, error) {
	res, err := s.urlRepo.GetUrl(ctx, shortUrl)

	if err == redis.Nil {
		log.Printf("Key %s does not exist", shortUrl)
		return "", http.StatusNotFound, fmt.Errorf("short URL not found: %s", shortUrl)
	}

	if err != nil {
		log.Printf("Error getting URL: %s", err.Error())
		return "", http.StatusInternalServerError, err
	}

	return res, http.StatusPermanentRedirect, nil
}
