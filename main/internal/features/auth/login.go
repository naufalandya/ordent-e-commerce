package auth

import (
	"commerce/internal/features/auth/models"
	"commerce/internal/features/auth/services"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *fiber.Ctx) error {
	var requestBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Invalid request body, unable to parse the request.",
		})
	}

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		var errorMessages []string
		for _, e := range err.(validator.ValidationErrors) {
			switch e.Field() {
			case "Email":
				if e.Tag() == "required" {
					errorMessages = append(errorMessages, "Email is required.")
				} else if e.Tag() == "email" {
					errorMessages = append(errorMessages, "Please provide a valid email address.")
				}
			case "Password":
				if e.Tag() == "required" {
					errorMessages = append(errorMessages, "Password is required.")
				} else if e.Tag() == "min" {
					errorMessages = append(errorMessages, "Password must be at least 6 characters long.")
				}
			}
		}

		firstError := "Validation failed"
		if len(errorMessages) > 0 {
			firstError = errorMessages[0]
		}

		log.Println("Validation errors:", errorMessages)
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: firstError,
		})
	}

	user, err := services.GetUserByEmail(requestBody.Email)
	if err != nil {
		log.Println("Error fetching user by email:", err)
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseFailed{
				Status:  "error",
				Message: "Invalid credentials, user not found.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Internal server error, please try again later.",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		log.Println("Password mismatch:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Invalid credentials, incorrect password.",
		})
	}

	response := models.ApiResponseSuccess{
		Status:  "success",
		Message: "Login successful",
		Data:    user,
	}

	return c.JSON(response)
}
