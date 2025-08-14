package routes

import (
	"user-consumer/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura as rotas do user-consumer
func SetupUserRoutes(r *gin.Engine, userController *controllers.UserController) {
	// Rota home
	r.GET("/", func(c *gin.Context) {
		endpoints := map[string]string{
			"GET    /users":     "Listar todos os usuários",
			"POST   /users":     "Criar novo usuário",
			"GET    /users/:id": "Obter usuário por ID",
			"PUT    /users/:id": "Atualizar usuário",
			"DELETE /users/:id": "Remover usuário",
			"GET    /healthz":   "Health check",
		}
		
		c.JSON(200, gin.H{
			"service": "user-consumer",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": endpoints,
			"docs": map[string]string{
				"health": "/healthz",
				"home":   "/",
			},
		})
	})

	// Grupo de rotas para usuários
	users := r.Group("/users")
	{
		users.GET("", userController.ListUsers)           // Listar todos
		users.POST("", userController.CreateUser)         // Criar
		users.GET("/:id", userController.GetUser)         // Obter por ID
		users.PUT("/:id", userController.UpdateUser)      // Atualizar
		users.DELETE("/:id", userController.DeleteUser)   // Remover
	}
}
