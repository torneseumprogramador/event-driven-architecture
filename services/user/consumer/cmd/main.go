package main

import (
	"context"
	"time"
	"user-consumer/internal/domain/entities"
	"user-consumer/internal/repo"
	pkgconfig "pkg/config"
	pkgkafka "pkg/kafka"
	pkglog "pkg/log"
	pkgoutboxdispatcher "pkg/outbox/dispatcher"
	pkgoutboxrepo "pkg/outbox/repository"
	pkgoutboxservices "pkg/outbox/services"

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
	pkglog.Setup("user-consumer")
	
	log.Info().
		Str("service", "user-consumer").
		Msg("iniciando user-consumer")
	
	// Conecta ao MySQL
	db, err := gorm.Open(mysql.Open(config.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("erro ao conectar ao MySQL")
	}
	
	// Auto-migra tabelas
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		log.Fatal().Err(err).Msg("erro ao migrar tabelas")
	}
	
	// Inicializa repositórios (para futuras funcionalidades)
	_ = repo.NewGormUserRepository(db)
	
	// Inicializa outbox
	outboxRepo := pkgoutboxrepo.NewGormOutboxRepository(db)
	outboxService := pkgoutboxservices.NewOutboxService(outboxRepo)
	
	// Inicializa Kafka producer
	kafkaProducer := pkgkafka.NewProducer(config.GetKafkaBrokers())
	defer kafkaProducer.Close()
	
	// Inicializa outbox dispatcher
	outboxDispatcher := pkgoutboxdispatcher.NewOutboxDispatcher(outboxService, kafkaProducer, time.Second)
	
	// Inicia dispatcher em background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go outboxDispatcher.Start(ctx)
	
	// Mantém o consumer rodando
	log.Info().Msg("user-consumer iniciado")
	select {} // Bloqueia indefinidamente
}
