package services

import (
	"context"
	"pkg/outbox/entities"
)

// OutboxService interface para operações de outbox
type OutboxService interface {
	CreateMessage(ctx context.Context, aggregate, eventType string, payload interface{}) (*entities.OutboxMessage, error)
	CreateMessageInTransaction(ctx context.Context, tx interface{}, aggregate, eventType string, payload interface{}) (*entities.OutboxMessage, error)
	GetPendingMessages(ctx context.Context, limit int) ([]entities.OutboxMessage, error)
	MarkMessageAsProcessed(ctx context.Context, id uint) error
	GetMessageByID(ctx context.Context, id uint) (*entities.OutboxMessage, error)
	GetPendingCount(ctx context.Context) (int64, error)
}
