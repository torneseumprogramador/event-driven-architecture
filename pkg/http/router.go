package http

import (
	"github.com/gin-gonic/gin"
)

// Router é a interface para configurar rotas
type Router interface {
	SetupRoutes(r *gin.Engine)
}

// SetupRouter configura o router com middlewares e rotas
func SetupRouter(r *gin.Engine) {
	// Middlewares globais
	r.Use(Logger())
	r.Use(Recovery())
	
	// Health check
	r.GET("/healthz", HealthCheck())
}

// HomeHandler retorna informações sobre os endpoints disponíveis
func HomeHandler(serviceName string, endpoints map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": serviceName,
			"version": "1.0.0",
			"status":  "running",
			"endpoints": endpoints,
			"docs": map[string]string{
				"health": "/healthz",
				"home":   "/",
			},
		})
	}
}
