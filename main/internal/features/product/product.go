package product

import (
	"commerce/internal/features/product/models"
	"commerce/internal/features/product/services"
	"commerce/internal/features/product/utils"
	"commerce/internal/repositories"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateProductHandler(c *fiber.Ctx) error {
	var requestBody struct {
		Name        string `json:"name" validate:"required,min=3,max=100"`
		Description string `json:"description"`
		Price       int    `json:"price" validate:"required,gt=0"`
		Stock       int    `json:"stock" validate:"required,min=0"`
		Category    string `json:"category"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return utils.CreateErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		errorMessages := utils.HandleValidationErrors(err)
		return utils.CreateErrorResponse(c, fiber.StatusBadRequest, errorMessages[0])
	}

	userID, err := utils.ParseUserID(c)
	if err != nil {
		return utils.CreateErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized. User ID not found or invalid.")
	}

	product := models.Product{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		UserID:      userID,
		Price:       requestBody.Price,
		Stock:       requestBody.Stock,
		Category:    requestBody.Category,
	}

	if err := services.CreateProduct(&product); err != nil {
		return utils.CreateErrorResponse(c, fiber.StatusInternalServerError, "Failed to create product")
	}

	uploadResults, err := utils.HandleFileUploads(c)
	if err != nil {
		return utils.CreateErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if len(uploadResults) > 0 {
		if err := services.CreateProductAttachments(product.ID, uploadResults); err != nil {
			return utils.CreateErrorResponse(c, fiber.StatusInternalServerError, "Failed to save product attachments")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(models.ApiResponseProductSuccess{
		Status:  "success",
		Message: "Product created successfully.",
		Data:    product,
	})
}

func DeleteProductHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Product ID is required.",
		})
	}

	productID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid product ID format.",
		})
	}

	deleteQuery := `
		DELETE FROM products
		WHERE id = $1
		RETURNING id
	`

	var deletedProductID int
	err = repositories.DB.QueryRow(context.Background(), deleteQuery, productID).Scan(&deletedProductID)
	if err != nil {
		log.Println("Error deleting product:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": fmt.Sprintf("Failed to delete product with ID %d.", productID),
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": fmt.Sprintf("Product with ID %d deleted successfully.", deletedProductID),
	})
}

func GetUserProductsHandler(c *fiber.Ctx) error {
	userIDInterface := c.Locals("user_id")
	if userIDInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized. User ID not found.",
		})
	}

	var userID int
	switch v := userIDInterface.(type) {
	case float64:
		userID = int(v)
	case string:
		var err error
		userID, err = strconv.Atoi(v)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Invalid user ID format.",
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid user ID type.",
		})
	}

	getProductsQuery := `
		SELECT 
			p.id, p.name, p.description, p.price, p.stock, p.category, p.created_at, p.updated_at,
			pa.id AS attachment_id, pa.link AS attachment_link
		FROM products p
		LEFT JOIN product_attachments pa ON p.id = pa.product_id
		WHERE p.user_id = $1
	`

	rows, err := repositories.DB.Query(context.Background(), getProductsQuery, userID)
	if err != nil {
		log.Println("Error fetching user products:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to retrieve products.",
		})
	}
	defer rows.Close()

	var products []Product
	var currentProduct *Product

	for rows.Next() {
		var product Product
		var attachmentID int
		var attachmentLink string

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Category,
			&product.CreatedAt,
			&product.UpdatedAt,
			&attachmentID,
			&attachmentLink,
		)
		if err != nil {
			log.Println("Error scanning product row:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "Error parsing product data.",
			})
		}

		if currentProduct == nil || currentProduct.ID != product.ID {
			if currentProduct != nil {
				products = append(products, *currentProduct)
			}
			currentProduct = &product
			currentProduct.Attachments = []ProductAttachment{}
		}

		if attachmentID != 0 {
			attachment := ProductAttachment{
				ID:   attachmentID,
				Link: attachmentLink,
			}
			currentProduct.Attachments = append(currentProduct.Attachments, attachment)
		}
	}

	if currentProduct != nil {
		products = append(products, *currentProduct)
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Products retrieved successfully.",
		"data":    products,
	})
}

type Product struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Price       float64             `json:"price"`
	Stock       int                 `json:"stock"`
	Category    string              `json:"category"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Attachments []ProductAttachment `json:"attachments"`
}

type ProductAttachment struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
}

func GetTransactionHistoryHandler(c *fiber.Ctx) error {
	userIDInterface := c.Locals("user_id")
	if userIDInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized. User ID not found.",
		})
	}

	var userID int
	switch v := userIDInterface.(type) {
	case float64:
		userID = int(v)
	case string:
		var err error
		userID, err = strconv.Atoi(v)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Invalid user ID format.",
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid user ID type.",
		})
	}

	getTransactionQuery := `
		SELECT 
			th.id AS transaction_id,
			th.midtrans_order_id,
			th.payment,
			th.total,
			th.status AS transaction_status,
			o.id AS order_id,
			p.name AS product_name,
			pa.link AS attachment_link,
			o.created_at AS order_created_at
		FROM transaction_history th
		JOIN orders o ON th.order_id = o.id
		JOIN products p ON o.product_id = p.id
		LEFT JOIN product_attachments pa ON p.id = pa.product_id
		WHERE o.user_id = $1
		ORDER BY th.created_at DESC
		LIMIT 1
	`

	// Execute the query
	rows, err := repositories.DB.Query(context.Background(), getTransactionQuery, userID)
	if err != nil {
		log.Println("Error fetching transaction history:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to retrieve transaction history.",
		})
	}
	defer rows.Close()

	var transactionHistory []TransactionHistoryResponse

	for rows.Next() {
		var th TransactionHistoryResponse
		var attachmentLink string

		err = rows.Scan(
			&th.TransactionID,
			&th.MidtransOrderID,
			&th.Payment,
			&th.Total,
			&th.TransactionStatus,
			&th.OrderID,
			&th.ProductName,
			&attachmentLink,
			&th.OrderCreatedAt,
		)
		if err != nil {
			log.Println("Error scanning transaction history row:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "Error parsing transaction history data.",
			})
		}

		if attachmentLink != "" {
			th.AttachmentLink = attachmentLink
		}

		transactionHistory = append(transactionHistory, th)
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Transaction history retrieved successfully.",
		"data":    transactionHistory,
	})
}

type TransactionHistoryResponse struct {
	TransactionID     int       `json:"transaction_id"`
	MidtransOrderID   string    `json:"midtrans_order_id"`
	Payment           string    `json:"payment"`
	Total             float64   `json:"total"`
	TransactionStatus bool      `json:"transaction_status"`
	OrderID           int       `json:"order_id"`
	ProductName       string    `json:"product_name"`
	AttachmentLink    string    `json:"attachment_link"`
	OrderCreatedAt    time.Time `json:"order_created_at"`
}
