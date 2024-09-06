package routes

import (
	"log"
	"os"
	"ryan-golang-url-shortener/database"
	"ryan-golang-url-shortener/helpers"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiration  time.Duration `json:"expiration"`
}

type response struct {
	URL                string `json:"url"`
	CustomShort        string `json:"short"`
	Expiration         string `json:"expiration"`
	RateLimitRemaining int    `json:"rate_limit"`
	RateLimitReset     string `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	log.Printf("Received request to shorten URL from %v\n", c.IP())
	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		return err
	}

	//Rate limit with redis in memory database based on IP address
	rateLimitDB := database.CreateClient(1)
	defer rateLimitDB.Close()

	if limit, err := handleRateLimiting(c, rateLimitDB); err != nil {
		log.Printf("Rate limiting error: %v", err)
		return err
	} else if limit > 0 {
		log.Printf("Rate limit exceeded for %v", c.IP())
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":            "Rate limit exceeded",
			"rate_limit_reset": limit.String(),
		})
	}

	if err := validateURL(body.URL); err != nil {
		log.Printf("Invalid URL: %v", err)
		return err
	}

	body.URL = helpers.EnforceHTTPS(body.URL)

	id := generateID(body.CustomShort)
	log.Printf("Generated ID: %s", id)

	customShortDB := database.CreateClient(0)
	defer customShortDB.Close()

	if err := checkURLInUse(customShortDB, id); err != nil {
		return err
	}

	if body.Expiration == 0 {
		body.Expiration = 24 * time.Hour
	}

	if err := saveURL(customShortDB, id, body.URL, body.Expiration); err != nil {
		log.Printf("URL already in use: %v", err)
		return err
	}

	//Decrement rate limit
	rateLimitDB.Decr(database.Ctx, c.IP())

	resp := populateResponse(c, rateLimitDB, body, id)
	log.Printf("Generated url %v for origin %v", resp.CustomShort, resp.URL)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func validateURL(url string) error {
	if !govalidator.IsURL(url) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid URL")
	}

	if !helpers.RemoveDomainError(url) {
		return fiber.NewError(fiber.StatusServiceUnavailable, "You cannot shorten this URL")
	}

	return nil
}

func handleRateLimiting(c *fiber.Ctx, rc *redis.Client) (time.Duration, error) {
	val, err := rc.Get(database.Ctx, c.IP()).Result()

	//Check to see if the IP address is in the database
	if err == redis.Nil {
		rc.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute).Err()
		return 0, nil
	} else if err != nil {
		return 0, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	// Find the current rate limit and determine if limit is exceeded
	valInt, _ := strconv.Atoi(val)

	if valInt <= 0 {
		limit, _ := rc.TTL(database.Ctx, c.IP()).Result()
		return limit, nil
	}

	return 0, nil
}

func generateID(customShort string) string {
	if customShort == "" {
		return uuid.New().String()[:6]
	}
	return customShort
}

func checkURLInUse(r *redis.Client, id string) error {
	val, _ := r.Get(database.Ctx, id).Result()
	if val != "" {
		return fiber.NewError(fiber.StatusForbidden, "URL short is already in use.")
	}
	return nil
}

func saveURL(r *redis.Client, id string, url string, expiration time.Duration) error {
	err := r.Set(database.Ctx, id, url, expiration).Err()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
	}
	return nil
}

func populateResponse(ctx *fiber.Ctx, rateLimitDB *redis.Client, req *request, id string) *response {
	resp := new(response)
	resp.URL = req.URL
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	resp.Expiration = req.Expiration.String()

	remRate, _ := rateLimitDB.Get(database.Ctx, ctx.IP()).Result()
	resp.RateLimitRemaining, _ = strconv.Atoi(remRate)

	ttl, _ := rateLimitDB.TTL(database.Ctx, ctx.IP()).Result()
	resp.RateLimitReset = ttl.String()

	return resp
}
