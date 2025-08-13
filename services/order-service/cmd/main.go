package main

import (
	"context"
	"fmt"
	"time"
	"order-service/internal/api/controllers"
	"order-service/internal/api/routes"
	"order-service/internal/domain/entities"
	"order-service/internal/services"
	"order-service/internal/repo"
	pkgconfig "pkg/config"
	pkgkafka "pkg/kafka"
	pkglog "pkg/log"
	pkgoutboxdispatcher "pkg/outbox/dispatcher"
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
		Msg("iniciando order-service")
	
	// Conecta ao MySQL
	db, err := gorm.Open(mysql.Open(config.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("erro ao conectar ao MySQL")
	}
	
	// Auto-migra tabelas
	if err := db.AutoMigrate(&entities.Order{}, &entities.OrderProduct{}); err != nil {
		log.Fatal().Err(err).Msg("erro ao migrar tabelas")
	}
	
	// Inicializa repositórios
	orderRepo := repo.NewGormOrderRepository(db)
	
	// Inicializa outbox
	outboxRepo := pkgoutboxrepo.NewGormRepository(db)
	outboxService := pkgoutboxservices.NewOutboxService(outboxRepo)
	
	// Inicializa serviços
	orderService := services.NewOrderService(orderRepo, outboxService, db)
	
	// Inicializa Kafka producer
	kafkaProducer := pkgkafka.NewProducer(config.GetKafkaBrokers())
	defer kafkaProducer.Close()
	
	// Inicializa outbox dispatcher
	outboxDispatcher := pkgoutboxdispatcher.NewOutboxDispatcher(outboxService, kafkaProducer, time.Second)
	
	// Inicia dispatcher em background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	orderController := controllers.NewOrderController(orderService)
	
	// Configura rotas
	routes.SetupOrderRoutes(router, orderController)
	
	// Inicia servidor
	log.Info().Msg("servidor iniciado")
	if err := router.Run(fmt.Sprintf(":%d", config.Port)); err != nil {
		log.Fatal().Err(err).Msg("erro ao iniciar servidor")
	}
}
