package routes

import (
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL short not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	rInf := database.CreateClient(1)
	defer rInf.Close()

	_ = rInf.Incr(database.Ctx, "counter")

	return c.Redirect(res, fiber.StatusMovedPermanently)
}
