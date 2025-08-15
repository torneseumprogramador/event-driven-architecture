package routes

import (
	"query-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupQueryRoutes configura as rotas da API de consultas
func SetupQueryRoutes(router *gin.Engine, userController *controllers.UserController, productController *controllers.ProductController, orderController *controllers.OrderController) {
	// Grupo de rotas para consultas
	queryGroup := router.Group("/q")
	{
		// Rotas de usuários
		queryGroup.GET("/users", userController.GetUsers)
		queryGroup.GET("/users/:id", userController.GetUserByID)

		// Rotas de produtos
		queryGroup.GET("/products", productController.GetProducts)
		queryGroup.GET("/products/:id", productController.GetProductByID)

		// Rotas de pedidos
		queryGroup.GET("/orders", orderController.GetOrders)
		queryGroup.GET("/orders/:id", orderController.GetOrder)
	}
}
