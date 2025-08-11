package main

import (
	"context"
	"fmt"
	"query-service/internal/api"
	"query-service/internal/consumer"
	"query-service/internal/projections"
	pkgconfig "pkg/config"
	pkgkafka "pkg/kafka"
	pkglog "pkg/log"
	pkghttp "pkg/http"
	pkgidempotency "pkg/idempotency"

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
		Msg("iniciando query-service")
	
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
	
	// Inicializa projeções
	orderProjection := projections.NewOrderProjection(db)
	productProjection := projections.NewProductProjection(db)
	userProjection := projections.NewUserProjection(db)
	
	// Inicializa Kafka producer
	kafkaProducer := pkgkafka.NewProducer(config.GetKafkaBrokers())
	defer kafkaProducer.Close()
	
	// Inicializa idempotência usando MongoDB
	idempotencyRepo := pkgidempotency.NewMongoRepository(db)
	idempotencyHandler := pkgidempotency.NewHandler(idempotencyRepo, config.ServiceName)
	
	// Inicializa consumidor de eventos
	eventConsumer := consumer.NewEventConsumer(
		orderProjection,
		productProjection,
		userProjection,
		kafkaProducer,
		idempotencyHandler,
	)
	
	// Inicia consumidores em background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Consumidor de user.created
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "user.created", "query-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleUserCreated)
	}()
	
	// Consumidor de product.created
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "product.created", "query-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleProductCreated)
	}()
	
	// Consumidor de product.updated
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "product.updated", "query-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleProductUpdated)
	}()
	
	// Consumidor de order.created
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.created", "query-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleOrderCreated)
	}()
	
	// Consumidor de order.paid
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.paid", "query-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleOrderPaid)
	}()
	
	// Consumidor de order.canceled
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.canceled", "query-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleOrderCanceled)
	}()
	
	// Consumidor de stock.reserved
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "stock.reserved", "query-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleStockReserved)
	}()
	
	// Consumidor de stock.released
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "stock.released", "query-service", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleStockReleased)
	}()
	
	// Configura Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Middlewares
	router.Use(pkghttp.Logger())
	router.Use(pkghttp.Recovery())
	
	// Healthcheck
	router.GET("/healthz", pkghttp.HealthCheck())
	
	// Handlers
	queryHandler := api.NewQueryHandler(orderProjection, productProjection, userProjection)
	
	// Rotas
	api := router.Group("/q")
	{
		orders := api.Group("/orders")
		{
			orders.GET("/:id", queryHandler.GetOrder)
			orders.GET("", queryHandler.GetOrders)
		}
		products := api.Group("/products")
		{
			products.GET("", queryHandler.GetProducts)
		}
		users := api.Group("/users")
		{
			users.GET("", queryHandler.GetUsers)
		}
	}
	
	// Inicia servidor
	log.Info().Msg("servidor iniciado")
	if err := router.Run(fmt.Sprintf(":%d", config.Port)); err != nil {
		log.Fatal().Err(err).Msg("erro ao iniciar servidor")
	}
}
