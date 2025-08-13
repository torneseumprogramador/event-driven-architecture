package repository

import (
	"context"
	"time"
	"pkg/outbox/entities"

	"gorm.io/gorm"
)

// GormOutboxRepository implementação usando GORM para outbox
type GormOutboxRepository struct {
	db *gorm.DB
}

// NewGormOutboxRepository cria um novo repositório GORM para outbox
func NewGormOutboxRepository(db *gorm.DB) OutboxRepository {
	return &GormOutboxRepository{db: db}
}

// Save salva uma mensagem na outbox
func (r *GormOutboxRepository) Save(ctx context.Context, message *entities.OutboxMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetPending retorna mensagens pendentes de processamento
func (r *GormOutboxRepository) GetPending(ctx context.Context, limit int) ([]entities.OutboxMessage, error) {
	var messages []entities.OutboxMessage
	err := r.db.WithContext(ctx).
		Where("processed_at IS NULL").
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

// MarkAsProcessed marca uma mensagem como processada
func (r *GormOutboxRepository) MarkAsProcessed(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entities.OutboxMessage{}).
		Where("id = ?", id).
		Update("processed_at", now).Error
}

// GetByID busca uma mensagem por ID
func (r *GormOutboxRepository) GetByID(ctx context.Context, id uint) (*entities.OutboxMessage, error) {
	var message entities.OutboxMessage
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}
