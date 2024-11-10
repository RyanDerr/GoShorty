package service

import (
	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/internal/domain/repository"
	"github.com/gin-gonic/gin"
)

type ShortenUrlService struct {
	urlRepo repository.IUrlRepository
}

type IUrlService interface {
	ShortenUrl(*gin.Context, *entity.ShortenUrl) (*entity.ShortenUrl, error)
	ResolveUrl(*gin.Context, string) (string, error)
}

func NewShortenUrlService(urlRepo repository.IUrlRepository) *ShortenUrlService {
	return &ShortenUrlService{
		urlRepo: urlRepo,
	}
}

func (s *ShortenUrlService) ShortenUrl(ctx *gin.Context, short *entity.ShortenUrl) (*entity.ShortenUrl, error) {
	return s.urlRepo.SaveUrl(ctx, short)
}

func (s *ShortenUrlService) ResolveUrl(ctx *gin.Context, short string) (string, error) {
	return s.urlRepo.GetUrl(ctx, short)
}
