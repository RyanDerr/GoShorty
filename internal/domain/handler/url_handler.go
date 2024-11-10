package handler

import (
	"net/http"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type UrlHandler struct {
	urlService service.IUrlService
}

func NewUrlHandler(urlService service.IUrlService) *UrlHandler {
	return &UrlHandler{
		urlService: urlService,
	}
}

func (h *UrlHandler) ShortenUrl(ctx *gin.Context) {
	var short entity.ShortenUrl
	err := ctx.BindJSON(&short)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	short, err = h.urlService.ShortenUrl(ctx, &short)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, short)

}
