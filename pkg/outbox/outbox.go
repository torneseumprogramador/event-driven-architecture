package outbox

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"pkg/outbox/dispatcher"
	"pkg/outbox/entities"
	"pkg/outbox/repository"
	"pkg/outbox/services"
)

// Producer interface para publicação no Kafka (re-exportada para compatibilidade)
type Producer = dispatcher.Producer

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
func NewGormRepository(db interface{}) repository.OutboxRepository {
	// Esta função é mantida para compatibilidade
	// Em novos códigos, use repository.NewGormOutboxRepository diretamente
	return repository.NewGormOutboxRepository(db.(*gorm.DB))
}

// NewDispatcher cria um novo dispatcher (função de conveniência)
func NewDispatcher(repo repository.OutboxRepository, producer Producer, interval time.Duration) dispatcher.OutboxDispatcher {
	// Esta função é mantida para compatibilidade
	// Em novos códigos, use dispatcher.NewOutboxDispatcher diretamente
	outboxService := services.NewOutboxService(repo)
	return dispatcher.NewOutboxDispatcher(outboxService, producer, interval)
}
