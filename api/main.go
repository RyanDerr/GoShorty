package main

import (
	"log"
	"os"
	"ryan-golang-url-shortener/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func routeSetup(app *fiber.App) {
	log.Println("Setting up routes")
	app.Post("/api/v1", routes.ShortenURL)
	app.Get("/:url", routes.ResolveURL)
}

func main() {
	log.Println("Starting server")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Creating new fiber app")
	app := fiber.New()
	routeSetup(app)

	log.Println("Setting up logger middleware")
	app.Use(logger.New())

	//Start the server on port and log any errors
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
