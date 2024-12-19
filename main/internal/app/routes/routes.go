package routes

import (
	"commerce/internal/app/middlewares"
	"commerce/internal/features/auth"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ApiV1Routes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	authRoutes(v1)
}

func authRoutes(v1 fiber.Router) {

	authRoute := v1.Group("/auth")
	// belum implementasi refresh
	authRoute.Get("/am-i-user", middlewares.BearerTokenAuth, middlewares.RoleCheck([]string{"USER"}), auth.ProtectedHandler)
	authRoute.Get("/am-i-admin", middlewares.BearerTokenAuth, middlewares.RoleCheck([]string{"ADMIN"}), auth.ProtectedHandler)
	authRoute.Post("/signin", auth.LoginHandler)
	authRoute.Post("/signup", auth.SignupHandler)
}

func productRoutes(v1 fiber.Router) {
	productRoute := v1.Group("/product")

	fmt.Println(productRoute)
}
