package product

import (
	"commerce/internal/features/product/models"
	"commerce/internal/features/product/services"
	"commerce/internal/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateProductHandler(c *fiber.Ctx) error {
	var requestBody struct {
		Name        string  `json:"name" validate:"required,min=3,max=100"`
		Description string  `json:"description"`
		Price       float64 `json:"price" validate:"required,gt=0"`
		Stock       int     `json:"stock" validate:"required,min=0"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiResponseProductFailed{
			Status:  "error",
			Message: "Invalid request body, unable to parse the request.",
		})
	}

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		var errorMessages []string
		for _, e := range err.(validator.ValidationErrors) {
			switch e.Field() {
			case "Name":
				errorMessages = append(errorMessages, "Product name is required and must be between 3 and 100 characters.")
			case "Price":
				errorMessages = append(errorMessages, "Price must be greater than 0.")
			case "Stock":
				errorMessages = append(errorMessages, "Stock must be a non-negative integer.")
			case "UserID":
				errorMessages = append(errorMessages, "User ID is required.")
			}
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiResponseProductFailed{
			Status:  "error",
			Message: errorMessages[0],
		})
	}

	fmt.Println(c.Locals("user_id"))

	userIDInterface := c.Locals("user_id")
	log.Printf("user_id: %v, type: %T", userIDInterface, userIDInterface)

	var userIDStr string

	switch v := userIDInterface.(type) {
	case string:
		userIDStr = v
	case float64:
		userIDStr = fmt.Sprintf("%v", v)
	default:
		log.Println("Error: user_id is not a string or float64")
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseProductFailed{
			Status:  "error",
			Message: "Unauthorized. User ID not found or invalid.",
		})
	}

	UserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("Error converting user_id to integer:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(models.ApiResponseProductFailed{
			Status:  "error",
			Message: "Invalid user ID format.",
		})
	}

	product := models.Product{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		UserID:      UserID,
		Price:       requestBody.Price,
		Stock:       requestBody.Stock,
	}

	if err := services.CreateProduct(&product); err != nil {
		log.Println("Error creating product:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseProductFailed{
			Status:  "error",
			Message: "Failed to create product. Please try again later.",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println("Error parsing multipart form:", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ApiResponseProductFailed{
			Status:  "error",
			Message: "Invalid form data.",
		})
	}

	fmt.Print("checkpoint 1")

	files := form.File["images"]
	var uploadResults []string

	if len(files) > 0 {
		if len(files) > 10 {

			fmt.Print("checkpoint 2")

			return c.Status(fiber.StatusBadRequest).JSON(models.ApiResponseProductFailed{
				Status:  "error",
				Message: "You can upload a maximum of 10 files.",
			})

		}

		uploadResults, err = utils.UploadImagesToImageKit(files)

		fmt.Print("checkpoint 3")

		if err != nil {
			log.Println("Error uploading images:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseProductFailed{
				Status:  "error",
				Message: "Failed to upload images. Please try again later.",
			})
		}

		fmt.Print("checkpoint 4")

	}

	fmt.Print("checkpoint 4")

	if len(uploadResults) > 0 {
		err = services.CreateProductAttachments(product.ID, uploadResults)

		fmt.Print("checkpoint 5")

		if err != nil {
			log.Println("Error creating product attachments:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ApiResponseProductFailed{
				Status:  "error",
				Message: "Failed to save product attachments. Please try again later.",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(models.ApiResponseProductSuccess{
		Status:  "success",
		Message: "Product created successfully.",
		Data:    product,
	})
}
