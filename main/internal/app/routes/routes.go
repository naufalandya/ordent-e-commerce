package routes

import (
	"commerce/internal/app/middlewares"
	"commerce/internal/features/order"

	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(app *fiber.App) {
	// app.Use(middlewares.LoggingMiddleware)
	// app.Get("/orders", order.GetOrdersHandler)
	// app.Get("/orders", order.GetOrdersHandler)
	// app.Get("/orders", order.GetOrdersHandler)
	// app.Get("/orders", order.GetOrdersHandler)
	// app.Get("/orders", order.GetOrdersHandler)
	// app.Get("/orders", order.GetOrdersHandler)
	// app.Get("/orders", order.GetOrdersHandler)
	app.Get("/orders", middlewares.LoggingMiddleware, middlewares.BearerMiddleware, order.GetOrdersHandler)

}
