package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/RyanDerr/GoShorty/api/database"
	"github.com/RyanDerr/GoShorty/api/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL         string `json:"url"`
	CustomShort string `json:"short"`
	Expiration  string `json:"expiration"`
}

type response struct {
	URL                string `json:"url"`
	CustomShort        string `json:"short"`
	Expiration         string `json:"expiration"`
	RateLimitRemaining int    `json:"rate_limit"`
	RateLimitReset     string `json:"rate_limit_reset"`
}

// ShortenURL godoc
// @Summary Shorten a URL
// @Description Shorten a given URL with an optional custom short and expiration time
// @Tags URL
// @Accept json
// @Produce json
// @Param request body request true "URL Shorten Request"
// @Success 201 {object} response
// @Failure 400 {object} ErrorResponse
// @Failure 429 {object} RateLimitExceededResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1 [post]
func ShortenURL(ctx *gin.Context) {
	log.Printf("Received request to shorten URL from %v\n", ctx.ClientIP())
	body := new(request)
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	rateLimitDB := database.CreateClient(1)
	defer rateLimitDB.Close()
	if limit, err := handleRateLimiting(ctx, rateLimitDB); err != nil {
		log.Printf("Rate limiting error: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal Server Error"})
		return
	} else if limit > 0 {
		log.Printf("Rate limit exceeded for %v", ctx.ClientIP())
		ctx.IndentedJSON(http.StatusTooManyRequests, RateLimitExceededResponse{
			Error:          "Rate limit exceeded",
			RateLimitReset: limit.String(),
		})
		return
	}

	if err := validateURL(body.URL); err != nil {
		log.Printf("Invalid URL: %v", err)
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid URL"})
		return
	}

	body.URL = helpers.EnforceHTTPS(body.URL)
	id := generateID(body.CustomShort)
	customShortDB := database.CreateClient(0)
	defer customShortDB.Close()
	if err := checkURLInUse(ctx, customShortDB, id); err != nil {
		ctx.IndentedJSON(http.StatusForbidden, ErrorResponse{Error: "URL short is already in use"})
		return
	}

	var expirationTime time.Duration
	if body.Expiration == "" {
		body.Expiration = "24h"
		expirationTime = 24 * time.Hour
	} else {
		var err error
		expirationTime, err = time.ParseDuration(body.Expiration)
		if err != nil {
			log.Printf("Invalid expiration time: %v", err)
			ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid expiration time"})
			return
		}
	}

	if err := saveURL(ctx, customShortDB, id, body.URL, expirationTime); err != nil {
		log.Printf("URL already in use: %v", err)
		ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal Server Error"})
		return
	}

	rateLimitDB.Decr(ctx, ctx.ClientIP())
	resp := populateResponse(ctx, rateLimitDB, body, id, expirationTime)
	log.Printf("Generated url %v for origin %v", resp.CustomShort, resp.URL)
	ctx.IndentedJSON(http.StatusCreated, resp)
}

func handleRateLimiting(ctx *gin.Context, rc *redis.Client) (time.Duration, error) {
	clientIP := ctx.ClientIP()
	val, err := rc.Get(ctx, clientIP).Result()

	//Check to see if the IP address is in the database
	if err == redis.Nil {
		rc.Set(ctx, clientIP, os.Getenv("API_QUOTA"), 30*time.Minute).Err()
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	// Find the current rate limit and determine if limit is exceeded
	valInt, _ := strconv.Atoi(val)

	if valInt <= 0 {
		limit, _ := rc.TTL(ctx, clientIP).Result()
		return limit, nil
	}

	return 0, nil
}

func validateURL(url string) error {
	if !govalidator.IsURL(url) {
		return fmt.Errorf("invalid url")
	}

	if !helpers.RemoveDomainError(url) {
		return fmt.Errorf("url is already shortened")
	}

	return nil
}

func generateID(customShort string) string {
	if customShort == "" {
		return uuid.New().String()[:6]
	}
	return customShort
}

func checkURLInUse(ctx *gin.Context, r *redis.Client, id string) error {
	val, _ := r.Get(ctx, id).Result()
	if val != "" {
		return fmt.Errorf("url already in use")
	}
	return nil
}

func saveURL(ctx *gin.Context, r *redis.Client, id string, url string, expiration time.Duration) error {
	err := r.Set(ctx, id, url, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func populateResponse(ctx *gin.Context, rateLimitDB *redis.Client, req *request, id string, exp time.Duration) *response {
	resp := new(response)
	resp.URL = req.URL
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	resp.Expiration = exp.String()

	clientIP := ctx.ClientIP()
	remRate, _ := rateLimitDB.Get(ctx, clientIP).Result()
	resp.RateLimitRemaining, _ = strconv.Atoi(remRate)

	ttl, _ := rateLimitDB.TTL(ctx, clientIP).Result()
	resp.RateLimitReset = ttl.String()

	return resp
}
