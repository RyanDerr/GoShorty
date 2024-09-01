package routes

import (
	"ryan-golang-url-shortener/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
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
	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	//Rate limits

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "You do not have permission to shorten this URL"})
	}

	body.URL = helpers.EnforceHTTPS(body.URL)
}
