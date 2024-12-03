package main

import (
	"log"
	"os"
	"time"

	"github.com/RyanDerr/GoShorty/api/middleware"
	"github.com/RyanDerr/GoShorty/api/routes"
	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/pkg/cache"
	"github.com/RyanDerr/GoShorty/pkg/database"
	"gorm.io/gorm"
)

func loadDatabase() (*gorm.DB, error) {
	db, err := database.CreateDatabaseConnection()
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return nil, err
	}

	log.Println("Database connection successful")
	return db, nil
}

func setupRateLimits() (*middleware.RateLimiter, *middleware.RateLimiter) {
	// Setup rate limits for shortening and resolving URLs with a refill rate of 1 token per 6 seconds, with a max of 10 tokens
	shortUrlRateLimit := middleware.NewRateLimiter(10, 1, time.Second*6)
	go shortUrlRateLimit.RefillTokens()

	// Setup rate limits for getting short URLs with a refill rate of 1 token per second, with a max of 60 tokens
	getShortRateLimit := middleware.NewRateLimiter(60, 1, time.Second)
	go getShortRateLimit.RefillTokens()

	return shortUrlRateLimit, getShortRateLimit
}

func main() {
	// Setup user database and Redis cache
	db, err := loadDatabase()
	if err != nil {
		log.Fatalf("Error loading database: %v", err)
	}

	redis, err := cache.CreateRedisClient(0)
	if err != nil {
		log.Fatalf("Error creating Redis client: %v", err)
	}

	// Setup rate limits for shortening and resolving URLs
	shortUrlRateLimit, getShortRateLimit := setupRateLimits()

	// Setup and run the Gin router
	port := os.Getenv("PORT")
	app := routes.SetupRouter(redis, db, shortUrlRateLimit, getShortRateLimit)

	if err := app.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
