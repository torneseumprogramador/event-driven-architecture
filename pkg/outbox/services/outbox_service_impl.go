package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"pkg/outbox/entities"
	"pkg/outbox/repository"

	"github.com/rs/zerolog/log"
)

// OutboxServiceImpl implementação do serviço de outbox
type OutboxServiceImpl struct {
	outboxRepo repository.OutboxRepository
}

// NewOutboxService cria um novo serviço de outbox
func NewOutboxService(outboxRepo repository.OutboxRepository) OutboxService {
	return &OutboxServiceImpl{
		outboxRepo: outboxRepo,
	}
}

// CreateMessage cria uma nova mensagem de outbox
func (s *OutboxServiceImpl) CreateMessage(ctx context.Context, aggregate, eventType string, payload interface{}) (*entities.OutboxMessage, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	message := &entities.OutboxMessage{
		Aggregate: aggregate,
		EventType: eventType,
		Payload:   string(payloadBytes),
		CreatedAt: time.Now(),
	}

	if err := s.outboxRepo.Save(ctx, message); err != nil {
		return nil, fmt.Errorf("erro ao salvar mensagem na outbox: %w", err)
	}

	log.Info().
		Uint("message_id", message.ID).
		Str("aggregate", aggregate).
		Str("event_type", eventType).
		Msg("mensagem criada na outbox")

	return message, nil
}

// GetPendingMessages retorna mensagens pendentes de processamento
func (s *OutboxServiceImpl) GetPendingMessages(ctx context.Context, limit int) ([]entities.OutboxMessage, error) {
	return s.outboxRepo.GetPending(ctx, limit)
}

// MarkMessageAsProcessed marca uma mensagem como processada
func (s *OutboxServiceImpl) MarkMessageAsProcessed(ctx context.Context, id uint) error {
	return s.outboxRepo.MarkAsProcessed(ctx, id)
}

// GetMessageByID busca uma mensagem por ID
func (s *OutboxServiceImpl) GetMessageByID(ctx context.Context, id uint) (*entities.OutboxMessage, error) {
	return s.outboxRepo.GetByID(ctx, id)
}

// GetPendingCount retorna o número de mensagens pendentes
func (s *OutboxServiceImpl) GetPendingCount(ctx context.Context) (int64, error) {
	// Para implementar isso, precisaríamos adicionar um método no repositório
	// Por enquanto, vamos buscar todas e contar
	messages, err := s.outboxRepo.GetPending(ctx, 10000) // Número alto para pegar todas
	if err != nil {
		return 0, err
	}
	return int64(len(messages)), nil
}
