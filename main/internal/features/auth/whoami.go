package auth

import "github.com/gofiber/fiber/v2"

func ProtectedHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	roles := c.Locals("roles")

	return c.JSON(fiber.Map{
		"user_id": userID,
		"roles":   roles,
	})
}
