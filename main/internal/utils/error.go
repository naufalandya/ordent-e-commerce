package utils

import (
	"commerce/internal/features/product/models"

	"github.com/gofiber/fiber/v2"
)

func CreateErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(models.ApiResponseProductFailed{
		Status:  "error",
		Message: message,
	})
}
