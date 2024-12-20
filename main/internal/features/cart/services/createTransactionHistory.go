package services

import (
	"commerce/internal/repositories"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

func CreateTransactionHistory(orderID, price, quantity int) (int, decimal.Decimal, error) {
	priceDecimal := decimal.NewFromInt(int64(price))
	quantityDecimal := decimal.NewFromInt(int64(quantity))

	total := priceDecimal.Mul(quantityDecimal)

	transactionQuery := `
		INSERT INTO transaction_history (order_id, midtrans_order_id, payment, total, created_at, updated_at)
		VALUES ($1, NULL, 'N/A', $2, now(), now())
		RETURNING id;  -- Mengambil id dari transaksi yang baru dibuat
	`

	var transactionID int
	err := repositories.DB.QueryRow(context.Background(), transactionQuery, orderID, total).Scan(&transactionID)
	if err != nil {
		return 0, decimal.Zero, fmt.Errorf("failed to create transaction history: %w", err)
	}

	return transactionID, total, nil
}
