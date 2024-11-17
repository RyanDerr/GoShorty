package main

import (
	"log"
	"os"

	"github.com/RyanDerr/GoShorty/api/routes"
	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/pkg/cache"
	"github.com/RyanDerr/GoShorty/pkg/database"
)

func loadDatabase() {
	db, err := database.CreateDatabaseConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Entry{})

	log.Println("Database connection successful")
}

func main() {

	loadDatabase()
	redis, err := cache.CreateRedisClient(0)

	if err != nil {
		log.Fatalf("Error creating Redis client: %v", err)
	}

	port := os.Getenv("PORT")

	app := routes.SetupRouter(redis)
	app.Run(":" + port)
}
