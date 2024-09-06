package main

import (
	"log"
	"os"

	"github.com/RyanDerr/GoShorty/api/helpers"
	"github.com/RyanDerr/GoShorty/api/routes"
	"github.com/gin-gonic/gin"
)

func routeSetup(app *gin.Engine) {
	log.Println("Setting up routes")
	app.POST("/api/v1", routes.ShortenURL)
	app.GET(":url", routes.ResolveURL)
}

// @title GoShorty API
// @version 1.0
// @description This is a sample URL shortener server.
// @host localhost:3000
// @BasePath /
func main() {
	log.Println("Starting server")
	err := helpers.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config")
	}
	log.Println("Creating new gin app")
	router := gin.Default()
	routeSetup(router)

	log.Println("Setting up logger middleware")
	router.Use(gin.Logger())

	//Start the server on port and log any errors
	log.Fatal(router.Run(os.Getenv("APP_PORT")))
}
