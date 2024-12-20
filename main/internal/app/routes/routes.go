package routes

import (
	"commerce/internal/app/middlewares"
	"commerce/internal/features/auth"
	"commerce/internal/features/cart"
	"commerce/internal/features/product"
	"commerce/internal/features/webhook"

	"github.com/gofiber/fiber/v2"
)

func ApiV1Routes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	authRoutes(v1)
	productOperationRoutes(v1)
	cartOperationRoutes(v1)
	webhookOperationRoutes(v1)
	transactionOperationRoutes(v1)
}

func authRoutes(v1 fiber.Router) {

	authRoute := v1.Group("/auth")
	// belum implementasi refresh
	// belum implementasi RBAC
	// belum implementasi rate limiter
	authRoute.Get("/am-i-user", middlewares.BearerTokenAuth, middlewares.RoleCheck([]string{"USER"}), auth.ProtectedHandler)
	authRoute.Get("/am-i-admin", middlewares.BearerTokenAuth, middlewares.RoleCheck([]string{"ADMIN"}), auth.ProtectedHandler)
	authRoute.Post("/signin", auth.LoginHandler)
	authRoute.Post("/signup", auth.SignupHandler)
}

func productOperationRoutes(v1 fiber.Router) {
	productRoute := v1.Group("/product")
	productRoute.Post("/", middlewares.BearerTokenAuth, middlewares.RoleCheck([]string{"USER"}), product.CreateProductHandler)
	productRoute.Delete("/:id", middlewares.BearerTokenAuth, middlewares.RoleCheck([]string{"USER", "ADMIN"}), product.DeleteProductHandler)
	productRoute.Get("/", middlewares.BearerTokenAuth, middlewares.RoleCheck([]string{"USER", "ADMIN"}), product.GetUserProductsHandler)
}

func cartOperationRoutes(v1 fiber.Router) {
	productRoute := v1.Group("/carts")
	productRoute.Post("/", middlewares.BearerTokenAuth, middlewares.RoleCheck([]string{"USER"}), cart.CreateOrderHandler)
}

func webhookOperationRoutes(v1 fiber.Router) {
	productRoute := v1.Group("/webhook")
	productRoute.Post("/", webhook.HandleMidtransWebhook)
}

func transactionOperationRoutes(v1 fiber.Router) {
	productRoute := v1.Group("/transactions")
	productRoute.Get("/history", middlewares.BearerTokenAuth, middlewares.RoleCheck([]string{"USER", "ADMIN"}), product.GetTransactionHistoryHandler)
}
