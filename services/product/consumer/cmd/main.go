package main

import (
	"context"
	"time"
	"product-consumer/internal/consumer"
	"product-consumer/internal/domain/entities"
	"product-consumer/internal/repo"
	pkgconfig "pkg/config"
	pkgkafka "pkg/kafka"
	pkglog "pkg/log"
	pkgoutboxdispatcher "pkg/outbox/dispatcher"
	pkgoutboxrepo "pkg/outbox/repository"
	pkgoutboxservices "pkg/outbox/services"
	pkgidempotency "pkg/idempotency"

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
	pkglog.Setup("product-consumer")
	
	log.Info().
		Str("service", "product-consumer").
		Msg("iniciando product-consumer")
	
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
	
	// Inicializa Kafka producer
	kafkaProducer := pkgkafka.NewProducer(config.GetKafkaBrokers())
	defer kafkaProducer.Close()
	
	// Inicializa outbox dispatcher
	outboxDispatcher := pkgoutboxdispatcher.NewOutboxDispatcher(outboxService, kafkaProducer, time.Second)
	
	// Inicializa idempotência
	idempotencyRepo := pkgidempotency.NewGormRepository(db)
	idempotencyHandler := pkgidempotency.NewHandler(idempotencyRepo, "product-consumer")
	
	// Inicializa consumidores
	orderConsumer := consumer.NewOrderConsumer(productRepo, kafkaProducer, idempotencyHandler)
	
	// Inicia consumidores em background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Consumidor de order.created
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.created", "product-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, orderConsumer.HandleOrderCreated)
	}()
	
	// Consumidor de order.canceled
	go func() {
		consumer := pkgkafka.NewConsumer(config.GetKafkaBrokers(), "order.canceled", "product-consumer", kafkaProducer)
		defer consumer.Close()
		consumer.Consume(ctx, orderConsumer.HandleOrderCanceled)
	}()
	
	// Inicia dispatcher em background
	go outboxDispatcher.Start(ctx)
	
	// Mantém o consumer rodando
	log.Info().Msg("product-consumer iniciado")
	select {} // Bloqueia indefinidamente
}
