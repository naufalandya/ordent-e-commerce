package services

import (
	"commerce/internal/features/product/models"
	"commerce/internal/repositories"
	"context"
	"fmt"
	"log"
	"time"
)

func handleDBError(query string, args []interface{}, err error) error {
	if err != nil {
		log.Printf("Error executing query: %s with params: %v - %v", query, args, err)
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func CreateProduct(product *models.Product) error {
	productQuery := `
		INSERT INTO products (name, description, user_id, price, stock, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now(), now())
		RETURNING id, description, category, created_at, updated_at;
	`

	var newProductID int
	var description, category string
	var createdAt, updatedAt time.Time

	err := repositories.DB.QueryRow(context.Background(), productQuery,
		product.Name,
		product.Description,
		product.UserID,
		product.Price,
		product.Stock,
		product.Category,
	).Scan(&newProductID, &description, &category, &createdAt, &updatedAt)

	if err != nil {
		handleDBError(productQuery, []interface{}{product.Name, product.Description, product.UserID, product.Price, product.Stock, product.Category}, err)
		return err
	}

	product.ID = newProductID
	product.Description = description
	product.Category = category
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

		if err := handleDBError(attachmentsQuery, []interface{}{link, productID}, err); err != nil {
			return err
		}
	}

	return nil
}
