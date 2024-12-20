package cart

import (
	"commerce/generated"
	"commerce/internal/features/cart/models"
	"commerce/internal/features/cart/services"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	fmt.Println(userID)

	order, err := services.CreateOrder(userID, requestBody.ProductID, requestBody.Quantity)
	if err != nil {
		log.Println("Error creating order:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Failed to create order. Please try again later.",
		})
	}

	transactionID, total, err := services.CreateTransactionHistory(order.ID, product.Price, order.Quantity)
	if err != nil {
		log.Println("Error creating transaction history:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Failed to record transaction. Please try again later.",
		})
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Failed to connect to gRPC server.",
		})
	}
	defer conn.Close()

	totalFloat64, exact := total.Float64()
	if !exact {
		// Handle the case where the conversion was not exact
		// For example, you might want to log a warning or handle the precision loss
		fmt.Println("Warning: Conversion from decimal to float64 was not exact")
	}

	client := generated.NewTransactionServiceClient(conn)

	req := &generated.TransactionRequest{
		TransactionHistoryId: int32(transactionID),
		Total:                totalFloat64,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := client.HandleTransaction(ctx, req)
	if err != nil {
		log.Printf("Error calling gRPC service: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseCartFailed{
			Status:  "error",
			Message: "Failed to send transactionHistoryId to gRPC server.",
		})
	}

	log.Printf("Response from gRPC server: %s", res.Message)

	return c.Status(fiber.StatusCreated).JSON(models.ApiResponseCartSuccess{
		Status:  "success",
		Message: "Order created and transaction history sent successfully.",
		Data:    order,
	})
}
