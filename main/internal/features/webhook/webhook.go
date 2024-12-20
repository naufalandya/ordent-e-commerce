package webhook

import (
	"commerce/internal/repositories"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HandleMidtransWebhook(c *fiber.Ctx) error {
	var requestBody struct {
		OrderID           string `json:"order_id"`
		TransactionStatus string `json:"transaction_status"`
		PaymentType       string `json:"payment_type"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		log.Println("Error parsing webhook request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request payload.",
		})
	}

	if requestBody.OrderID == "" || requestBody.TransactionStatus == "" || requestBody.PaymentType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Missing required fields: order_id, transaction_status, or payment_type.",
		})
	}

	if requestBody.TransactionStatus == "pending" {
		return c.JSON(fiber.Map{"status": true})
	}

	transactionQuery := `
		SELECT id, status
		FROM transaction_history
		WHERE midtrans_order_id = $1
	`
	var transactionID int
	var existingStatus bool

	err := repositories.DB.QueryRow(context.Background(), transactionQuery, requestBody.OrderID).
		Scan(&transactionID, &existingStatus)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Transaction with order ID %s not found.", requestBody.OrderID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": "Transaction not found.",
			})
		}
		log.Println("Error fetching transaction:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Internal server error.",
		})
	}

	if existingStatus {
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "Payment has already been completed!",
		})
	}

	var newStatus bool
	switch requestBody.TransactionStatus {
	case "capture", "settlement":
		newStatus = true
	case "deny", "cancel", "expire":
		newStatus = false
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid transaction status.",
		})
	}

	updateTransactionQuery := `
		UPDATE transaction_history
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err = repositories.DB.Exec(context.Background(), updateTransactionQuery, newStatus, time.Now(), transactionID)
	if err != nil {
		log.Println("Error updating transaction status:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to update transaction status.",
		})
	}

	insertPaymentQuery := `
    INSERT INTO payments (transaction_id, payment_type, payment_status, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (transaction_id) DO NOTHING
`
	_, err = repositories.DB.Exec(context.Background(), insertPaymentQuery, transactionID, requestBody.PaymentType, requestBody.TransactionStatus, time.Now(), time.Now())
	if err != nil {
		log.Println("Error inserting payment record:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to create payment record.",
		})
	}

	checkPaymentQuery := `
		SELECT id
		FROM payments
		WHERE transaction_id = $1
	`
	var paymentID int
	err = repositories.DB.QueryRow(context.Background(), checkPaymentQuery, transactionID).Scan(&paymentID)
	if err != nil {
		log.Println("Error fetching payment record:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to verify payment record.",
		})
	}

	updatePaymentQuery := `
		UPDATE payments
		SET payment_type = $1, payment_status = $2, updated_at = $3
		WHERE id = $4
	`
	_, err = repositories.DB.Exec(context.Background(), updatePaymentQuery, requestBody.PaymentType, requestBody.TransactionStatus, time.Now(), paymentID)
	if err != nil {
		log.Println("Error updating payment record:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to update payment record.",
		})
	}

	updateOrderStatusQuery := `
		UPDATE orders
		SET status_id = $1, updated_at = $2
		WHERE id = (
			SELECT order_id
			FROM transaction_history
			WHERE id = $3
		)
	`
	_, err = repositories.DB.Exec(context.Background(), updateOrderStatusQuery, 3, time.Now(), transactionID)
	if err != nil {
		log.Println("Error updating order status:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to update order status.",
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Transaction updated successfully.",
	})
}
