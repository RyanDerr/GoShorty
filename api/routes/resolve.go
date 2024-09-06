package routes

import (
	"log"
	"ryan-golang-url-shortener/database"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func ResolveURL(c *fiber.Ctx) error {
	short := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()

	log.Printf("Attempting to retrieve entry for %v", short)
	res, err := r.Get(database.Ctx, short).Result()

	if err == redis.Nil {
		log.Printf("URL short not found for short %v \n", short)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL short not found"})
	} else if err != nil {
		log.Printf("Error retrieving URL from Redis: %v \n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	rInf := database.CreateClient(1)
	defer rInf.Close()

	_ = rInf.Incr(database.Ctx, c.IP())

	log.Println("Redirecting to:", res)
	return c.Redirect(res, fiber.StatusMovedPermanently)
}
