package routes

import (
	"product-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes configura as rotas do product-api
func SetupProductRoutes(r *gin.Engine, productController *controllers.ProductController) {
	// Rota home
	r.GET("/", func(c *gin.Context) {
		endpoints := map[string]string{
			"GET    /products":     "Listar todos os produtos",
			"POST   /products":     "Criar novo produto",
			"GET    /products/:id": "Obter produto por ID",
			"PUT    /products/:id": "Atualizar produto",
			"DELETE /products/:id": "Remover produto",
			"GET    /healthz":      "Health check",
		}
		
		c.JSON(200, gin.H{
			"service": "product-api",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": endpoints,
			"docs": map[string]string{
				"health": "/healthz",
				"home":   "/",
			},
		})
	})

	// Grupo de rotas para produtos
	products := r.Group("/products")
	{
		products.GET("", productController.ListProducts)           // Listar todos
		products.POST("", productController.CreateProduct)         // Criar
		products.GET("/:id", productController.GetProduct)         // Obter por ID
		products.PUT("/:id", productController.UpdateProduct)      // Atualizar
		products.DELETE("/:id", productController.DeleteProduct)   // Remover
	}
}
