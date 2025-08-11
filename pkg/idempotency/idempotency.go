package idempotency

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ProcessedEvent representa um evento já processado
type ProcessedEvent struct {
	EventID      string    `gorm:"primaryKey;column:event_id"`
	ServiceName  string    `gorm:"not null;column:service_name"`
	ProcessedAt  time.Time `gorm:"not null;column:processed_at"`
}

// TableName especifica o nome da tabela
func (ProcessedEvent) TableName() string {
	return "processed_events"
}

// Repository interface para controle de idempotência
type Repository interface {
	IsProcessed(ctx context.Context, eventID, serviceName string) (bool, error)
	MarkAsProcessed(ctx context.Context, eventID, serviceName string) error
}

// GormRepository implementação usando GORM
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository cria um novo repositório GORM
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// IsProcessed verifica se um evento já foi processado
func (r *GormRepository) IsProcessed(ctx context.Context, eventID, serviceName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&ProcessedEvent{}).
		Where("event_id = ? AND service_name = ?", eventID, serviceName).
		Count(&count).Error
	
	if err != nil {
		return false, fmt.Errorf("erro ao verificar evento processado: %w", err)
	}
	
	return count > 0, nil
}

// MarkAsProcessed marca um evento como processado
func (r *GormRepository) MarkAsProcessed(ctx context.Context, eventID, serviceName string) error {
	processedEvent := &ProcessedEvent{
		EventID:     eventID,
		ServiceName: serviceName,
		ProcessedAt: time.Now(),
	}
	
	return r.db.WithContext(ctx).Create(processedEvent).Error
}

// Handler wrapper para garantir idempotência
type Handler struct {
	repo        Repository
	serviceName string
}

// NewHandler cria um novo handler de idempotência
func NewHandler(repo Repository, serviceName string) *Handler {
	return &Handler{
		repo:        repo,
		serviceName: serviceName,
	}
}

// ProcessWithIdempotency processa um evento garantindo idempotência
func (h *Handler) ProcessWithIdempotency(ctx context.Context, eventID string, processor func() error) error {
	// Verifica se já foi processado
	processed, err := h.repo.IsProcessed(ctx, eventID, h.serviceName)
	if err != nil {
		return fmt.Errorf("erro ao verificar idempotência: %w", err)
	}
	
	if processed {
		// Evento já processado, ignora
		return nil
	}
	
	// Processa o evento
	if err := processor(); err != nil {
		return fmt.Errorf("erro ao processar evento: %w", err)
	}
	
	// Marca como processado
	if err := h.repo.MarkAsProcessed(ctx, eventID, h.serviceName); err != nil {
		return fmt.Errorf("erro ao marcar evento como processado: %w", err)
	}
	
	return nil
}
