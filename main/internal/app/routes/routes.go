package routes

import (
	"commerce/internal/app/middlewares"
	"commerce/internal/features/auth"

	"github.com/gofiber/fiber/v2"
)

func ApiV1Routes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	authRoutes(v1)
}

func authRoutes(v1 fiber.Router) {

	authRoute := v1.Group("/auth")

	authRoute.Post("/tes", middlewares.LoggingMiddleware, auth.LoginHandler)
	authRoute.Post("/signin", auth.LoginHandler)
}
