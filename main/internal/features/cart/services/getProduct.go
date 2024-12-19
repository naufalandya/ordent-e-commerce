package services

import (
	"commerce/internal/features/product/models"
	"commerce/internal/repositories"
	"context"
	"fmt"
	"log"
)

func GetProductByID(productID int) (*models.Product, error) {
	productQuery := `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	err := repositories.DB.QueryRow(context.Background(), productQuery, productID).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("product not found")
		}
		log.Println("Error fetching product:", err)
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	return &product, nil
}
