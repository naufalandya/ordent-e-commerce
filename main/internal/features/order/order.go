package order

import (
	repositories "commerce/internal/repositories"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Order struct {
	ID       int     `json:"id"`
	Customer string  `json:"customer"`
	Amount   float64 `json:"amount"`
}

type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetOrdersHandler(c *fiber.Ctx) error {
	orders, err := GetOrders()
	if err != nil {
		response := ApiResponse{
			Status:  "error",
			Message: "Error fetching orders",
			Data:    nil,
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := ApiResponse{
		Status:  "success",
		Message: "Orders fetched successfully",
		Data:    orders,
	}
	return c.JSON(response)
}

func GetOrders() ([]Order, error) {
	rows, err := repositories.DB.Query(context.Background(), "SELECT id, customer, amount FROM orders")
	if err != nil {
		log.Println("Error fetching orders:", err)
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.Customer, &order.Amount); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, err
	}
	return orders, nil
}
