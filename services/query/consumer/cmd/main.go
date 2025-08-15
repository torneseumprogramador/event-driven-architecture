package main

import (
	"context"
	"query-consumer/internal/consumer"
	"query-consumer/internal/repository"
	"query-consumer/internal/services"
	pkgconfig "pkg/config"
	pkgkafka "pkg/kafka"
	pkglog "pkg/log"
	pkgidempotency "pkg/idempotency"

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
		Msg("iniciando query-consumer")
	
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
	
	// Inicializa repositories
	orderRepository := repository.NewMongoOrderRepository(db)
	userRepository := repository.NewMongoUserRepository(db)
	productRepository := repository.NewMongoProductRepository(db)
	
	// Inicializa services
	orderService := services.NewOrderService(orderRepository)
	userService := services.NewUserService(userRepository)
	productService := services.NewProductService(productRepository)
	
	// Inicializa Kafka producer
	kafkaProducer := pkgkafka.NewProducer(config.GetKafkaBrokers())
	defer kafkaProducer.Close()
	
	// Inicializa idempotência usando MongoDB
	idempotencyRepo := pkgidempotency.NewMongoRepository(db)
	idempotencyHandler := pkgidempotency.NewHandler(idempotencyRepo, config.ServiceName)
	
	// Inicializa consumidor de eventos
	eventConsumer := consumer.NewEventConsumer(
		orderService,
		productService,
		userService,
		userRepository,
		productRepository,
		kafkaProducer,
		idempotencyHandler,
	)
	
	// Inicia consumidores em background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Consumidor de user.created
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "user.created", "query-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleUserCreated)
	}()
	
	// Consumidor de product.created
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "product.created", "query-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleProductCreated)
	}()
	
	// Consumidor de product.updated
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "product.updated", "query-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleProductUpdated)
	}()
	
	// Consumidor de order.created
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.created", "query-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleOrderCreated)
	}()
	
	// Consumidor de order.paid
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.paid", "query-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleOrderPaid)
	}()
	
	// Consumidor de order.canceled
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.canceled", "query-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleOrderCanceled)
	}()
	
	// Consumidor de stock.reserved
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "stock.reserved", "query-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleStockReserved)
	}()
	
	// Consumidor de stock.released
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "stock.released", "query-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, eventConsumer.HandleStockReleased)
	}()
	
	log.Info().Msg("query-consumer iniciado")
	
	// Aguarda indefinidamente
	select {}
}
