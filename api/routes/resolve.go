package routes

import (
	"log"
	"net/http"

	"github.com/RyanDerr/GoShorty/api/database"
	"github.com/RyanDerr/GoShorty/api/repositories"
	"github.com/RyanDerr/GoShorty/api/services"
	"github.com/gin-gonic/gin"
)

func ResolveURL(ctx *gin.Context) {
	short := ctx.Param("url")
	client := database.CreateRedisClient(0)
	defer client.Close()
	repo := repositories.NewURLRepository(client)
	service := services.NewURLService(ctx, repo)

	url, err := service.ResolveURL(short)
	if err != nil {
		if err.Error() == "URL not found" {
			log.Printf("URL short not found for short %v \n", short)
			ctx.IndentedJSON(http.StatusNotFound, ErrorResponse{Error: "URL not found"})
		} else {
			log.Printf("Error retrieving URL from Redis: %v \n", err)
			ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal Server Error"})
		}
		return
	}

	log.Println("Redirecting to:", url)
	ctx.Redirect(http.StatusMovedPermanently, url)
}
