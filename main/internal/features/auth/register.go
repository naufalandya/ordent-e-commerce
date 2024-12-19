package auth

import (
	"commerce/internal/features/auth/models"
	"commerce/internal/features/auth/services"
	"commerce/internal/utils"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func SignupHandler(c *fiber.Ctx) error {
	var requestBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		Name     string `json:"name" validate:"required"`
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
			case "Name":
				if e.Tag() == "required" {
					errorMessages = append(errorMessages, "Name is required.")
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

	_, err := services.GetUserByEmail(requestBody.Email)
	if err == nil {
		log.Println("Email already registered")
		return c.Status(fiber.StatusConflict).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Email is already registered. Please use another email.",
		})
	}

	hashedPassword, err := utils.HashPassword(requestBody.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Internal server error while hashing password.",
		})
	}

	newUser := models.User{
		Email:    requestBody.Email,
		Password: hashedPassword,
		Name:     requestBody.Name,
	}

	if err := services.CreateUser(&newUser, 1); err != nil {
		log.Println("Error creating user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Failed to create user. Please try again later.",
		})
	}

	response := models.ApiResponseSuccess{
		Status:  "success",
		Message: "Signup successful. Welcome to the platform!",
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"id":    newUser.ID,
				"name":  newUser.Name,
				"email": newUser.Email,
			},
		},
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
