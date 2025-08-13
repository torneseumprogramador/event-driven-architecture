package repository

import (
	"context"
	"pkg/outbox/entities"
)

// OutboxRepository interface para operações de outbox
type OutboxRepository interface {
	Save(ctx context.Context, message *entities.OutboxMessage) error
	GetPending(ctx context.Context, limit int) ([]entities.OutboxMessage, error)
	MarkAsProcessed(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*entities.OutboxMessage, error)
}
