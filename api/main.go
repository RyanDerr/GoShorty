package main

import (
	"log"

	"github.com/RyanDerr/GoShorty/api/config"
	"github.com/RyanDerr/GoShorty/api/routes"
	"github.com/gin-gonic/gin"
)

const (
	port = ":8080"
)

func routeSetup(app *gin.Engine) {
	log.Println("Setting up routes")
	app.POST("/api/v1", routes.ShortenURL)
	app.GET(":url", routes.ResolveURL)
}

// @title GoShorty API
// @version 1.0
// @description URL shortener server.
// @host localhost:3000
// @BasePath /
func main() {
	log.Println("Starting server")
	config.LoadConfig()
	log.Println("Creating new gin app")
	router := gin.Default()
	routeSetup(router)

	log.Println("Setting up logger middleware")
	router.Use(gin.Logger())

	//Start the server on port and log any errors
	log.Fatal(router.Run(port))
}
