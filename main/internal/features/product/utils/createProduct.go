package utils

import (
	"commerce/internal/features/product/models"
	"commerce/internal/utils"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const maxUploadFiles = 10

func CreateErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(models.ApiResponseProductFailed{
		Status:  "error",
		Message: message,
	})
}

func ParseUserID(c *fiber.Ctx) (int, error) {
	userIDInterface := c.Locals("user_id")

	fmt.Printf("user_id from context: %v, type: %T\n", userIDInterface, userIDInterface)

	switch v := userIDInterface.(type) {
	case string:
		userID, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid user ID format: %w", err)
		}
		return userID, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("user ID not found or invalid type")
	}
}

func HandleValidationErrors(err error) []string {
	var errorMessages []string
	for _, e := range err.(validator.ValidationErrors) {
		switch e.Field() {
		case "Name":
			errorMessages = append(errorMessages, "Product name is required and must be between 3 and 100 characters.")
		case "Price":
			errorMessages = append(errorMessages, "Price must be greater than 0.")
		case "Stock":
			errorMessages = append(errorMessages, "Stock must be a non-negative integer.")
		}
	}
	return errorMessages
}

func HandleFileUploads(c *fiber.Ctx) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %v", err)
	}
	files := form.File["images"]
	if len(files) > maxUploadFiles {
		return nil, fmt.Errorf("you can upload a maximum of %d files", maxUploadFiles)
	}
	return utils.UploadImagesToImageKit(files)
}
