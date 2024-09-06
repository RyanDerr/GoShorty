package routes

import (
	"log"
	"net/http"

	"github.com/RyanDerr/GoShorty/api/database"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func ResolveURL(ctx *gin.Context) {
	short := ctx.Param("url")
	r := database.CreateClient(0)
	defer r.Close()

	log.Printf("Attempting to retrieve entry for %v", short)
	res, err := r.Get(ctx, short).Result()

	if err == redis.Nil {
		log.Printf("URL short not found for short %v \n", short)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	} else if err != nil {
		log.Printf("Error retrieving URL from Redis: %v \n", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	rInf := database.CreateClient(1)
	defer rInf.Close()

	_ = rInf.Incr(ctx, ctx.ClientIP())

	log.Println("Redirecting to:", res)
	ctx.Redirect(http.StatusMovedPermanently, res)
}
