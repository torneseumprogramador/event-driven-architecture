package main

import (
	"context"
	"fmt"
	"query-api/internal/api/routes"
	"query-api/internal/controllers"
	"query-api/internal/repo"
	"query-api/internal/services"
	pkgconfig "pkg/config"
	pkglog "pkg/log"
	pkghttp "pkg/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Carrega configuração
	config, err := pkgconfig.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("erro ao carregar configuração")
	}
	
	// Configura logger
	pkglog.Setup(config.ServiceName)
	
	log.Info().
		Str("service", config.ServiceName).
		Int("port", config.Port).
		Msg("iniciando query-api")
	
	// Conecta ao MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal().Err(err).Msg("erro ao conectar ao MongoDB")
	}
	defer client.Disconnect(context.Background())
	
	// Verifica conexão
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal().Err(err).Msg("erro ao verificar conexão com MongoDB")
	}
	
	db := client.Database("ecommerce")
	
	// Inicializa repositórios
	userRepo := repo.NewMongoUserRepository(db)
	productRepo := repo.NewMongoProductRepository(db)
	orderRepo := repo.NewMongoOrderRepository(db)
	
	// Inicializa serviços
	queryService := services.NewQueryService(userRepo, productRepo, orderRepo)
	
	// Configura Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Middlewares
	router.Use(pkghttp.Logger())
	router.Use(pkghttp.Recovery())
	
	// Healthcheck
	router.GET("/healthz", pkghttp.HealthCheck())
	
	// Configura controllers
	queryController := controllers.NewQueryController(queryService)
	
	// Configura rotas
	routes.SetupQueryRoutes(router, queryController)
	
	// Inicia servidor
	log.Info().Msg("servidor iniciado")
	if err := router.Run(fmt.Sprintf(":%d", config.Port)); err != nil {
		log.Fatal().Err(err).Msg("erro ao iniciar servidor")
	}
}
