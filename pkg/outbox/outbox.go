package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"pkg/outbox/dispatcher"
	"pkg/outbox/entities"
	"pkg/outbox/repository"
	"pkg/outbox/services"
)

// Producer interface para publicação no Kafka
type Producer interface {
	PublishEvent(ctx context.Context, topic string, event interface{}) error
}

// CreateOutboxMessage cria uma nova mensagem de outbox (função de conveniência)
func CreateOutboxMessage(aggregate, eventType string, payload interface{}) (*entities.OutboxMessage, error) {
	// Esta função é mantida para compatibilidade
	// Em novos códigos, use o OutboxService diretamente
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}
	
	return &entities.OutboxMessage{
		Aggregate: aggregate,
		EventType: eventType,
		Payload:   string(payloadBytes),
		CreatedAt: time.Now(),
	}, nil
}

// NewGormRepository cria um novo repositório GORM (função de conveniência)
func NewGormRepository(db interface{}) repository.Repository {
	// Esta função é mantida para compatibilidade
	// Em novos códigos, use repository.NewGormRepository diretamente
	return repository.NewGormRepository(db.(*gorm.DB))
}

// NewDispatcher cria um novo dispatcher (função de conveniência)
func NewDispatcher(repo repository.Repository, producer Producer, interval time.Duration) *dispatcher.OutboxDispatcher {
	// Esta função é mantida para compatibilidade
	// Em novos códigos, use dispatcher.NewOutboxDispatcher diretamente
	outboxService := services.NewOutboxService(repo)
	return dispatcher.NewOutboxDispatcher(outboxService, producer, interval)
}
