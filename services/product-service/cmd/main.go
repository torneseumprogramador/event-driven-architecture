package main

import (
	"context"
	"fmt"
	"time"
	"product-service/internal/api/controllers"
	"product-service/internal/api/routes"
	"product-service/internal/consumer"
	"product-service/internal/domain"
	"product-service/internal/services"
	"product-service/internal/repo"
	pkgconfig "pkg/config"
	pkgkafka "pkg/kafka"
	pkglog "pkg/log"
	pkgoutbox "pkg/outbox"
	pkghttp "pkg/http"
	pkgidempotency "pkg/idempotency"

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
		Msg("iniciando product-service")
	
	// Conecta ao MySQL
	db, err := gorm.Open(mysql.Open(config.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("erro ao conectar ao MySQL")
	}
	
	// Auto-migra tabelas
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		log.Fatal().Err(err).Msg("erro ao migrar tabelas")
	}
	
	// Inicializa repositórios
	productRepo := repo.NewGormProductRepository(db)
	
	// Inicializa serviços
	productService := services.NewProductService(db)
	
	// Inicializa Kafka producer
	kafkaProducer := pkgkafka.NewProducer(config.GetKafkaBrokers())
	defer kafkaProducer.Close()
	
	// Inicializa outbox dispatcher
	outboxRepo := pkgoutbox.NewGormRepository(db)
	outboxDispatcher := pkgoutbox.NewDispatcher(outboxRepo, kafkaProducer, time.Second)
	
	// Inicializa idempotência
	idempotencyRepo := pkgidempotency.NewGormRepository(db)
	idempotencyHandler := pkgidempotency.NewHandler(idempotencyRepo, config.ServiceName)
	
	// Inicializa consumidores
	orderConsumer := consumer.NewOrderConsumer(productRepo, kafkaProducer, idempotencyHandler)
	
	// Inicia consumidores em background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Consumidor de order.created
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.created", "product-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, orderConsumer.HandleOrderCreated)
	}()
	
	// Consumidor de order.canceled
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.canceled", "product-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, orderConsumer.HandleOrderCanceled)
	}()
	
	// Inicia dispatcher em background
	go outboxDispatcher.Start(ctx)
	
	// Configura Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Middlewares
	router.Use(pkghttp.Logger())
	router.Use(pkghttp.Recovery())
	
	// Healthcheck
	router.GET("/healthz", pkghttp.HealthCheck())
	
	// Configura controllers
	productController := controllers.NewProductController(productRepo, productService)
	
	// Configura rotas
	routes.SetupProductRoutes(router, productController)
	
	// Inicia servidor
	log.Info().Msg("servidor iniciado")
	if err := router.Run(fmt.Sprintf(":%d", config.Port)); err != nil {
		log.Fatal().Err(err).Msg("erro ao iniciar servidor")
	}
}
