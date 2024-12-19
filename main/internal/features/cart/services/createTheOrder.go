package services

import (
	"commerce/internal/features/cart/models"
	"commerce/internal/repositories"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

func CreateOrder(userID, productID, quantity int) (*models.Order, error) {

	ctx := context.Background()

	tx, err := repositories.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Println("Error starting transaction:", err)
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	var product models.Product
	err = tx.QueryRow(context.Background(), `SELECT id, name, stock, price FROM products WHERE id = $1 FOR UPDATE`, productID).
		Scan(&product.ID, &product.Name, &product.Stock, &product.Price)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("product not found")
		}
		log.Println("Error fetching product:", err)
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	if product.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock for the product")
	}

	_, err = tx.Exec(context.Background(), `UPDATE products SET stock = stock - $1 WHERE id = $2`, quantity, productID)
	if err != nil {
		log.Println("Error updating stock:", err)
		return nil, fmt.Errorf("failed to update product stock: %w", err)
	}

	orderQuery := `
		INSERT INTO orders (user_id, product_id, status_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, now(), now())
		RETURNING id, user_id, product_id, status_id, quantity, created_at, updated_at
	`
	var order models.Order
	err = tx.QueryRow(context.Background(), orderQuery, userID, productID, 1, quantity).
		Scan(&order.ID, &order.UserID, &order.ProductID, &order.StatusID, &order.Quantity, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		log.Println("Error creating order:", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Println("Error committing transaction:", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
}
