package routes

import (
	"fmt"

	"github.com/RyanDerr/GoShorty/api/middleware"
	"github.com/RyanDerr/GoShorty/internal/domain/handler"
	"github.com/RyanDerr/GoShorty/internal/domain/repository"
	"github.com/RyanDerr/GoShorty/internal/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	BaseRoute = "/api/v1"
	UrlRoute  = "/url"
	Shorten   = "/shorten"
	Resolve   = "/:short"
	AuthRoute = "/auth"
	Register  = "/register"
	Login     = "/login"
)

//	@title			GoShorty API
//	@version		0.1.0
//	@host			localhost
//	@BasePath		/api/v1
//	@schemes		http
//	@schemes		https
//	@description	This is the API for GoShorty, a URL shortening service
func SetupRouter(redis *redis.Client, db *gorm.DB, shortenUrlRateLimit, getShortUrlRateLimit *middleware.RateLimiter) *gin.Engine {
	router := gin.Default()

	userHandler, urlHandler := setupHandlers(redis, db)

	urlRoute := GetUrlRoute()
	authRoute := GetAuthRoute()

	authenticate := router.Group(authRoute)
	{
		authenticate.POST(Register, userHandler.RegisterUser)
		authenticate.POST(Login, userHandler.LoginUser)
	}

	shorten := router.Group(urlRoute).Use(middleware.JWTAuthMiddleware())
	{
		shorten.POST(Shorten, shortenUrlRateLimit.IsRateLimited(), urlHandler.ShortenUrl)
		shorten.GET(Resolve, getShortUrlRateLimit.IsRateLimited(), urlHandler.ResolveUrl)
	}

	return router
}

func GetUrlRoute() string {
	return fmt.Sprintf("%s%s", BaseRoute, UrlRoute)
}

func GetAuthRoute() string {
	return fmt.Sprintf("%s%s", BaseRoute, AuthRoute)
}

func GetShortenRoute() string {
	return fmt.Sprintf("%s%s", GetUrlRoute(), Shorten)
}

func GetResolveRoute() string {
	return fmt.Sprintf("%s%s", GetUrlRoute(), Resolve)
}

func GetRegisterRoute() string {
	return fmt.Sprintf("%s%s", GetAuthRoute(), Register)
}

func GetLoginRoute() string {
	return fmt.Sprintf("%s%s", GetAuthRoute(), Login)
}

func setupHandlers(redis *redis.Client, db *gorm.DB) (*handler.UserHandler, *handler.UrlHandler) {
	userService, urlService := setupServices(redis, db)

	userHandler := handler.NewUserHandler(userService)
	urlHandler := handler.NewUrlHandler(urlService)

	return userHandler, urlHandler
}

func setupServices(redis *redis.Client, db *gorm.DB) (*service.UserService, *service.ShortenUrlService) {
	userRepo, urlRepo := setupRepositories(redis, db)

	userService := service.NewUserService(userRepo)
	urlService := service.NewShortenUrlService(urlRepo)

	return userService, urlService
}

func setupRepositories(redis *redis.Client, db *gorm.DB) (*repository.UserRepository, *repository.UrlRepository) {
	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewUrlRepository(redis)

	return userRepo, urlRepo
}
