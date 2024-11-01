package services

import (
	"github.com/RyanDerr/GoShorty/api/repositories"
	"github.com/gin-gonic/gin"
)

type URLService struct {
	ctx  *gin.Context
	repo *repositories.URLRepository
}

func NewURLService(ctx *gin.Context, repo *repositories.URLRepository) *URLService {
	return &URLService{
		ctx:  ctx,
		repo: repo,
	}
}
