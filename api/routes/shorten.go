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
	URL                string        `json:"url"`
	CustomShort        string        `json:"short"`
	Expiration         time.Duration `json:"expiration"`
	RateLimitRemaining int           `json:"rate_limit"`
	RateLimitReset     time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	log.Printf("Received request to shorten URL from %v\n", c.IP())
	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	//Rate limit with redis in memory database based on IP address
	r2 := database.CreateClient(1)
	defer r2.Close()

	if err := handleRateLimiting(c, r2); err != nil {
		log.Printf("Rate limiting error: %v", err)
		return err
	}

	if err := validateURL(body.URL); err != nil {
		log.Printf("Invalid URL: %v", err)
		return err
	}

	body.URL = helpers.EnforceHTTPS(body.URL)

	id := generateID(body.CustomShort)
	log.Printf("Generated ID: %s", id)

	r := database.CreateClient(0)
	defer r.Close()

	if err := checkURLInUse(r, id); err != nil {
		return err
	}

	if body.Expiration == 0 {
		body.Expiration = 24 * time.Hour
	}

	if err := saveURL(r, id, body.URL, body.Expiration); err != nil {
		log.Printf("URL already in use: %v", err)
		return err
	}

	resp := response{
		URL:                body.URL,
		CustomShort:        id,
		Expiration:         body.Expiration,
		RateLimitRemaining: 10,
		RateLimitReset:     30 * time.Minute,
	}

	//Decrement rate limit
	r2.Decr(database.Ctx, c.IP())

	populateResponse(c, &resp, r2)
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

func handleRateLimiting(c *fiber.Ctx, rc *redis.Client) error {
	val, err := rc.Get(database.Ctx, c.IP()).Result()

	//Check to see if the IP address is in the database
	if err == redis.Nil {
		rc.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute).Err()
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	} else {
		// Find the current rate limit and determine if limit is exceeded
		valInt, _ := strconv.Atoi(val)

		if valInt <= 0 {
			limit, _ := rc.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded!",
				"rate_limit_reset": limit.Minutes(),
			})
		}
	}

	return nil
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

func populateResponse(ctx *fiber.Ctx, resp *response, r *redis.Client) {
	remRate, _ := r.Get(database.Ctx, ctx.IP()).Result()
	resp.RateLimitRemaining, _ = strconv.Atoi(remRate)

	ttl, _ := r.TTL(database.Ctx, ctx.IP()).Result()
	resp.RateLimitReset = ttl * time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + resp.CustomShort
}
