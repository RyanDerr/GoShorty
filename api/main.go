package main

import (
	"log"
	"os"

	api "github.com/RyanDerr/GoShorty/api/routes"
	"github.com/RyanDerr/GoShorty/pkg/cache"
)

func main() {

	redis, err := cache.CreateRedisClient(0)

	if err != nil {
		log.Fatalf("Error creating Redis client: %v", err)
	}

	port := os.Getenv("PORT")

	app := api.SetupRouter(redis)
	app.Run(":" + port)
}
