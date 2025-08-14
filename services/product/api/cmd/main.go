package main

import (
	"fmt"
	"product-api/internal/api/controllers"
	"product-api/internal/api/routes"
	"product-api/internal/domain/entities"
	"product-api/internal/services"
	"product-api/internal/repo"
	pkgconfig "pkg/config"
	pkglog "pkg/log"
	pkgoutboxrepo "pkg/outbox/repository"
	pkgoutboxservices "pkg/outbox/services"
	pkghttp "pkg/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		Msg("iniciando product-api")
	
	// Conecta ao MySQL
	db, err := gorm.Open(mysql.Open(config.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("erro ao conectar ao MySQL")
	}
	
	// Auto-migra tabelas
	if err := db.AutoMigrate(&entities.Product{}); err != nil {
		log.Fatal().Err(err).Msg("erro ao migrar tabelas")
	}
	
	// Inicializa repositórios
	productRepo := repo.NewGormProductRepository(db)
	
	// Inicializa outbox
	outboxRepo := pkgoutboxrepo.NewGormOutboxRepository(db)
	outboxService := pkgoutboxservices.NewOutboxService(outboxRepo)
	
	// Inicializa serviços
	productService := services.NewProductService(productRepo, outboxService, db)
	
	// Configura Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Middlewares
	router.Use(pkghttp.Logger())
	router.Use(pkghttp.Recovery())
	
	// Healthcheck
	router.GET("/healthz", pkghttp.HealthCheck())
	
	// Configura controllers
	productController := controllers.NewProductController(productService)
	
	// Configura rotas
	routes.SetupProductRoutes(router, productController)
	
	// Inicia servidor
	log.Info().Msg("servidor iniciado")
	if err := router.Run(fmt.Sprintf(":%d", config.Port)); err != nil {
		log.Fatal().Err(err).Msg("erro ao iniciar servidor")
	}
}
