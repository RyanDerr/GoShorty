package routes

import (
	"log"
	"ryan-golang-url-shortener/database"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()

	res, err := r.Get(database.Ctx, url).Result()

	if err != redis.Nil {
		log.Println("URL short not found:", url)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL short not found"})
	} else if err != nil {
		log.Println("Error retrieving URL from Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	rInf := database.CreateClient(1)
	defer rInf.Close()

	_ = rInf.Incr(database.Ctx, "counter")

	log.Println("Redirecting to:", res)
	return c.Redirect(res, fiber.StatusMovedPermanently)
}
