package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config representa a configuração da aplicação
type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	Port        int    `mapstructure:"PORT"`
	
	// MySQL (Write Model)
	MySQLDSN string `mapstructure:"MYSQL_DSN"`
	
	// MongoDB (Read Model)
	MongoURI string `mapstructure:"MONGO_URI"`
	
	// Kafka
	KafkaBrokers string `mapstructure:"KAFKA_BROKERS"`
	
	// Outbox
	OutboxPollInterval string `mapstructure:"OUTBOX_POLL_INTERVAL"`
}

// Load carrega a configuração do ambiente
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	
	// Configurações padrão
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("SERVICE_NAME", "unknown-service")
	viper.SetDefault("MYSQL_DSN", "ecommerce:ecommerce@tcp(mysql:3306)/ecommerce?parseTime=true")
	viper.SetDefault("MONGO_URI", "mongodb://admin:admin@mongo:27017/ecommerce?authSource=admin")
	viper.SetDefault("KAFKA_BROKERS", "kafka:9092")
	viper.SetDefault("OUTBOX_POLL_INTERVAL", "1s")
	
	// Lê variáveis de ambiente
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// Tenta ler arquivo .env se existir
	viper.ReadInConfig()
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("erro ao deserializar configuração: %w", err)
	}
	
	return &config, nil
}

// GetKafkaBrokers retorna os brokers do Kafka como slice
func (c *Config) GetKafkaBrokers() []string {
	return strings.Split(c.KafkaBrokers, ",")
}
