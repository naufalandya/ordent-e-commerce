package main

import (
	"commerce/internal/app/middlewares"
	"commerce/internal/app/routes"
	"commerce/internal/repositories"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	repositories.InitDB()

	app := fiber.New()

	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("ENV")
	if env == "development" {
		middlewares.SetupLogger(app)
	}

	routes.ApiV1Routes(app)

	log.Fatal(app.Listen(":8080"))
}
