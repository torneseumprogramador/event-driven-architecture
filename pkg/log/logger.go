package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup configura o logger global
func Setup(serviceName string) {
	// Configuração para desenvolvimento
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	
	// Log estruturado em JSON para produção
	if os.Getenv("ENV") != "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	
	// Adiciona campos padrão
	log.Logger = log.With().
		Str("service", serviceName).
		Timestamp().
		Logger()
	
	// Configura o logger global
	zerolog.DefaultContextLogger = &log.Logger
}

// GetLogger retorna o logger configurado
func GetLogger() zerolog.Logger {
	return log.Logger
}
