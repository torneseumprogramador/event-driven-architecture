package repository

import (
	"context"
	"pkg/outbox/entities"
	"time"

	"gorm.io/gorm"
)

// Repository interface para operações de outbox
type Repository interface {
	Save(ctx context.Context, message *entities.OutboxMessage) error
	GetPending(ctx context.Context, limit int) ([]entities.OutboxMessage, error)
	MarkAsProcessed(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*entities.OutboxMessage, error)
}

// GormRepository implementação usando GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository cria um novo repositório GORM
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Save salva uma mensagem na outbox
func (r *GormRepository) Save(ctx context.Context, message *entities.OutboxMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetPending retorna mensagens pendentes de processamento
func (r *GormRepository) GetPending(ctx context.Context, limit int) ([]entities.OutboxMessage, error) {
	var messages []entities.OutboxMessage
	err := r.db.WithContext(ctx).
		Where("processed_at IS NULL").
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

// MarkAsProcessed marca uma mensagem como processada
func (r *GormRepository) MarkAsProcessed(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entities.OutboxMessage{}).
		Where("id = ?", id).
		Update("processed_at", now).Error
}

// GetByID busca uma mensagem por ID
func (r *GormRepository) GetByID(ctx context.Context, id uint) (*entities.OutboxMessage, error) {
	var message entities.OutboxMessage
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}
