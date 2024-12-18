package middlewares

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupLogger(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		elapsed := time.Since(start).Milliseconds()
		log.Printf("Request took %dms", elapsed)
		return err
	})
}

func LoggingMiddleware(c *fiber.Ctx) error {
	fmt.Printf("Request: %s %s\n", c.Method(), c.Path())
	return c.Next()
}
