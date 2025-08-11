package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

// MessageHandler função para processar mensagens
type MessageHandler func(ctx context.Context, message []byte) error

// Consumer wrapper para o consumidor Kafka
type Consumer struct {
	reader *kafka.Reader
	producer *Producer
	maxRetries int
}

// NewConsumer cria um novo consumidor Kafka
func NewConsumer(brokers []string, topic, groupID string, producer *Producer) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		Logger:   kafka.LoggerFunc(log.Printf),
	})
	
	return &Consumer{
		reader:     reader,
		producer:   producer,
		maxRetries: 5,
	}
}

// Consume inicia o consumo de mensagens com retry e DLQ
func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) error {
	log.Info().
		Str("topic", c.reader.Config().Topic).
		Str("group_id", c.reader.Config().GroupID).
		Msg("iniciando consumo de mensagens")
	
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("consumo interrompido")
			return ctx.Err()
		default:
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msg("erro ao ler mensagem")
				continue
			}
			
			// Processa mensagem com retry
			if err := c.processWithRetry(ctx, message, handler); err != nil {
				log.Error().
					Err(err).
					Str("topic", message.Topic).
					Int("partition", message.Partition).
					Int64("offset", message.Offset).
					Msg("falha ao processar mensagem após retries")
				
				// Publica na DLQ
				if err := c.producer.PublishToDLQ(ctx, message.Topic, string(message.Value), err.Error()); err != nil {
					log.Error().Err(err).Msg("erro ao publicar na DLQ")
				}
			}
		}
	}
}

// processWithRetry processa uma mensagem com retry exponencial
func (c *Consumer) processWithRetry(ctx context.Context, message kafka.Message, handler MessageHandler) error {
	var lastErr error
	
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Backoff exponencial: 1s, 2s, 4s, 8s, 16s
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			log.Info().
				Int("attempt", attempt).
				Dur("backoff", backoff).
				Msg("tentativa de retry")
			
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}
		
		if err := handler(ctx, message.Value); err != nil {
			lastErr = err
			log.Error().
				Err(err).
				Int("attempt", attempt+1).
				Msg("erro ao processar mensagem")
			continue
		}
		
		// Sucesso
		log.Info().
			Int("attempt", attempt+1).
			Msg("mensagem processada com sucesso")
		return nil
	}
	
	return fmt.Errorf("falha após %d tentativas: %w", c.maxRetries+1, lastErr)
}

// Close fecha o consumidor
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// UnmarshalMessage deserializa uma mensagem JSON
func UnmarshalMessage(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
