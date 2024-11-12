package api

import (
	"fmt"

	"github.com/RyanDerr/GoShorty/internal/domain/handler"
	"github.com/RyanDerr/GoShorty/internal/domain/repository"
	"github.com/RyanDerr/GoShorty/internal/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	BaseRoute = "/api/v1"
	UrlRoute  = "/url"
	Shorten   = "/shorten"
	Resolve   = "/:short"
)

//	@title			GoShorty API
//	@version		0.1.0
//	@host			localhost
//	@BasePath		/api/v1
//	@schemes		http
//	@schemes		https
//	@description	This is the API for GoShorty, a URL shortening service
func SetupRouter(redis *redis.Client) *gin.Engine {
	router := gin.Default()
	urlRepo := repository.NewUrlRepository(redis)
	urlService := service.NewShortenUrlService(urlRepo)
	urlHandler := handler.NewUrlHandler(urlService)

	urlRoute := GetUrlRoute()

	url := router.Group(urlRoute)
	{
		url.POST(Shorten, urlHandler.ShortenUrl)
		url.GET(Resolve, urlHandler.ResolveUrl)
	}

	return router
}

func GetUrlRoute() string {
	return fmt.Sprintf("%s%s", BaseRoute, UrlRoute)
}

func GetShortenRoute() string {
	return fmt.Sprintf("%s%s", GetUrlRoute(), Shorten)
}

func GetResolveRoute() string {
	return fmt.Sprintf("%s%s", GetUrlRoute(), Resolve)
}
