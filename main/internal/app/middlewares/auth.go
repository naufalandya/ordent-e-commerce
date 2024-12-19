package middlewares

import (
	"fmt"
	"log"
	"os"
	"strings"

	"commerce/internal/features/auth/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func BearerTokenAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Authorization header is missing.",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Bearer token missing.",
		})
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Println("JWT_SECRET_KEY is not set in .env file")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Server configuration error. Please contact support.",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println("Error parsing JWT:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Invalid or expired token.",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("user_id", claims["user_id"])
		c.Locals("roles", claims["roles"])

		return c.Next()
	}

	log.Println("Invalid token")
	return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseFailed{
		Status:  "error",
		Message: "Invalid or expired token.",
	})
}

func RoleCheck(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles := c.Locals("roles").([]interface{})

		for _, role := range roles {
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(models.ApiResponseFailed{
			Status:  "error",
			Message: "Access forbidden: insufficient permissions.",
		})
	}
}
