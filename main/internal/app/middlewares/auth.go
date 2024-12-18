package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func BearerMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid Authorization header",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if !validateToken(token) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	c.Locals("user", "example_user_id")

	return c.Next()
}

func validateToken(token string) bool {
	return token == "valid_token_example"
}
