package product

import (
	"commerce/internal/features/product/models"
	"commerce/internal/features/product/services"
	"commerce/internal/features/product/utils"

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
