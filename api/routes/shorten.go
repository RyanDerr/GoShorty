package routes

import (
	"log"
	"net/http"

	"github.com/RyanDerr/GoShorty/api/database"
	models "github.com/RyanDerr/GoShorty/api/modules"
	"github.com/RyanDerr/GoShorty/api/repositories"
	"github.com/RyanDerr/GoShorty/api/services"
	"github.com/gin-gonic/gin"
)

// ShortenURL godoc
// @Summary Shorten a URL
// @Description Shorten a given URL with an optional custom short and expiration time
// @Tags URL
// @Accept json
// @Produce json
// @Param request body models.ShortenRequest true "URL Shorten Request"
// @Success 201 {object} models.ShortenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 429 {object} RateLimitExceededResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1 [post]
func ShortenURL(ctx *gin.Context) {
	log.Printf("Received request to shorten URL from %v\n", ctx.ClientIP())
	body := new(models.ShortenRequest)
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	client := database.CreateRedisClient(0)
	defer client.Close()

	repo := repositories.NewURLRepository(client)
	service := services.NewURLService(ctx, repo)

	resp, err := service.ShortenURL(body)
	if err != nil {
		log.Printf("Error shortening URL: %v\n", err)
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, resp)
}
