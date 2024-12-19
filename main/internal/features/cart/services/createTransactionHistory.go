package services

import (
	"commerce/internal/repositories"
	"context"
	"fmt"
)

func CreateTransactionHistory(orderID, price, quantity int) error {
	transactionQuery := `
		INSERT INTO transaction_history (order_id, price, order_id, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now());
	`

	_, err := repositories.DB.Exec(context.Background(), transactionQuery,
		orderID,
		price,
		quantity,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction history: %w", err)
	}

	return nil
}
