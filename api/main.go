package main

import (
	"log"
	"os"

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

	db.AutoMigrate(&entity.User{})

	log.Println("Database connection successful")
	return db, nil
}

func main() {

	db, err := loadDatabase()
	if err != nil {
		log.Fatalf("Error loading database: %v", err)
	}

	redis, err := cache.CreateRedisClient(0)

	if err != nil {
		log.Fatalf("Error creating Redis client: %v", err)
	}

	port := os.Getenv("PORT")

	app := routes.SetupRouter(redis, db)
	app.Run(":" + port)
}
