package cart

import (
	"commerce/internal/features/cart/models"
	"commerce/internal/features/cart/services"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateOrderHandler(c *fiber.Ctx) error {
	var requestBody struct {
		ProductID int `json:"product_id" validate:"required"`
		Quantity  int `json:"quantity" validate:"required,min=1"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Invalid request body, unable to parse the request.",
		})
	}

	userIDInterface := c.Locals("user_id")

	var userIDStr string
	switch v := userIDInterface.(type) {
	case string:
		userIDStr = v
	case float64:
		userIDStr = fmt.Sprintf("%v", v)
	default:
		log.Println("Error: user_id is not a string or float64")
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Unauthorized. User ID not found or invalid.",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("Error converting user_id to integer:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Invalid user ID format.",
		})
	}

	product, err := services.GetProductByID(requestBody.ProductID)
	if err != nil {
		log.Println("Error retrieving product:", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Product not found.",
		})
	}

	if product.Stock < requestBody.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Insufficient stock for the product.",
		})
	}

	order, err := services.CreateOrder(userID, requestBody.ProductID, requestBody.Quantity)
	if err != nil {
		log.Println("Error creating order:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Failed to create order. Please try again later.",
		})
	}

	err = services.CreateTransactionHistory(order.ID, product.Price, order.Quantity)
	if err != nil {
		log.Println("Error creating transaction history:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Failed to record transaction. Please try again later.",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.ApiResponseCartSuccess{
		Status:  "success",
		Message: "Order created successfully.",
		Data:    order,
	})
}
