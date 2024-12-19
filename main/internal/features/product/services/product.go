package services

import (
	"commerce/internal/features/product/models"
	"commerce/internal/repositories"
	"context"
	"fmt"
	"log"
	"time"
)

func CreateProduct(product *models.Product) error {
	productQuery := `
		INSERT INTO products (name, description, user_id, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, now(), now())
		RETURNING id, description, created_at, updated_at;
	`

	var newProductID int
	var description string
	var createdAt, updatedAt time.Time // Change to time.Time

	// Execute the query
	err := repositories.DB.QueryRow(context.Background(), productQuery,
		product.Name,
		product.Description,
		product.UserID,
		product.Price,
		product.Stock,
	).Scan(&newProductID, &description, &createdAt, &updatedAt)

	if err != nil {
		log.Println("Error inserting product into database:", err)
		return fmt.Errorf("failed to create product: %w", err)
	}

	if description == "" {
		description = "No description available"
	}

	product.ID = newProductID
	product.Description = description
	product.CreatedAt = createdAt
	product.UpdatedAt = updatedAt

	return nil
}

func CreateProductAttachments(productID int, links []string) error {
	attachmentsQuery := `
		INSERT INTO product_attachments (link, product_id, created_at, updated_at)
		VALUES ($1, $2, now(), now());
	`

	for _, link := range links {
		_, err := repositories.DB.Exec(context.Background(), attachmentsQuery, link, productID)
		if err != nil {
			log.Println("Error inserting product attachment into database:", err)
			return fmt.Errorf("failed to create product attachment: %w", err)
		}
	}

	return nil
}
