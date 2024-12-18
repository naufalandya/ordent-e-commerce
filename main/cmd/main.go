package main

import (
	"commerce/internal/app/routes"
	repositories "commerce/internal/repositories"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	repositories.InitDB()

	app := fiber.New()

	routes.OrderRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
