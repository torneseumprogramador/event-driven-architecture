package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

// Producer wrapper para o produtor Kafka
type Producer struct {
	writer *kafka.Writer
}

// NewProducer cria um novo produtor Kafka
func NewProducer(brokers []string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false, // Síncrono para garantir entrega
		Logger:       kafka.LoggerFunc(log.Printf),
	}
	
	return &Producer{writer: writer}
}

// PublishEvent publica um evento no tópico especificado
func (p *Producer) PublishEvent(ctx context.Context, topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}
	
	message := kafka.Message{
		Topic: topic,
		Key:   []byte(fmt.Sprintf("%d", time.Now().UnixNano())),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "Content-Type", Value: []byte("application/json")},
			{Key: "Timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		},
	}
	
	log.Info().
		Str("topic", topic).
		Str("event_type", fmt.Sprintf("%T", event)).
		Msg("publicando evento no Kafka")
	
	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("erro ao publicar evento no tópico %s: %w", topic, err)
	}
	
	log.Info().
		Str("topic", topic).
		Msg("evento publicado com sucesso")
	
	return nil
}

// PublishToDLQ publica uma mensagem na DLQ
func (p *Producer) PublishToDLQ(ctx context.Context, originalTopic string, event interface{}, errorMsg string) error {
	dlqTopic := originalTopic + ".dlq"
	
	dlqEvent := map[string]interface{}{
		"original_topic": originalTopic,
		"error_message":  errorMsg,
		"timestamp":      time.Now().Format(time.RFC3339),
		"original_event": event,
	}
	
	return p.PublishEvent(ctx, dlqTopic, dlqEvent)
}

// Close fecha o produtor
func (p *Producer) Close() error {
	return p.writer.Close()
}
