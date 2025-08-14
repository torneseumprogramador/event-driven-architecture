package routes

import (
	"order-consumer/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes configura as rotas do order-consumer
func SetupOrderRoutes(r *gin.Engine, orderController *controllers.OrderController) {
	// Rota home
	r.GET("/", func(c *gin.Context) {
		endpoints := map[string]string{
			"GET    /orders":     "Listar todos os pedidos",
			"POST   /orders":     "Criar novo pedido",
			"GET    /orders/:id": "Obter pedido por ID",
			"PUT    /orders/:id": "Atualizar pedido",
			"DELETE /orders/:id": "Remover pedido",
			"POST   /orders/:id/pay": "Pagar pedido",
			"POST   /orders/:id/cancel": "Cancelar pedido",
			"GET    /healthz":    "Health check",
		}
		
		c.JSON(200, gin.H{
			"service": "order-consumer",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": endpoints,
			"docs": map[string]string{
				"health": "/healthz",
				"home":   "/",
			},
		})
	})

	// Grupo de rotas para pedidos
	orders := r.Group("/orders")
	{
		orders.GET("", orderController.ListOrders)           // Listar todos
		orders.POST("", orderController.CreateOrder)         // Criar
		orders.GET("/:id", orderController.GetOrder)         // Obter por ID
		orders.PUT("/:id", orderController.UpdateOrder)      // Atualizar
		orders.DELETE("/:id", orderController.DeleteOrder)   // Remover
		orders.POST("/:id/pay", orderController.PayOrder)    // Pagar
		orders.POST("/:id/cancel", orderController.CancelOrder) // Cancelar
	}
}
