package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// OutboxMessage representa uma mensagem na tabela outbox
type OutboxMessage struct {
	ID          uint           `gorm:"primaryKey"`
	Aggregate   string         `gorm:"not null"`
	EventType   string         `gorm:"not null"`
	Payload     string         `gorm:"type:json;not null"`
	Headers     sql.NullString `gorm:"type:json"`
	CreatedAt   time.Time      `gorm:"not null"`
	ProcessedAt *time.Time     `gorm:"null"`
}

// TableName especifica o nome da tabela
func (OutboxMessage) TableName() string {
	return "outbox"
}

// Repository interface para operações de outbox
type Repository interface {
	Save(ctx context.Context, message *OutboxMessage) error
	GetPending(ctx context.Context, limit int) ([]OutboxMessage, error)
	MarkAsProcessed(ctx context.Context, id uint) error
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
func (r *GormRepository) Save(ctx context.Context, message *OutboxMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetPending retorna mensagens pendentes de processamento
func (r *GormRepository) GetPending(ctx context.Context, limit int) ([]OutboxMessage, error) {
	var messages []OutboxMessage
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
		Model(&OutboxMessage{}).
		Where("id = ?", id).
		Update("processed_at", now).Error
}

// Dispatcher processa mensagens da outbox e as publica no Kafka
type Dispatcher struct {
	repo     Repository
	producer Producer
	interval time.Duration
}

// Producer interface para publicação no Kafka
type Producer interface {
	PublishEvent(ctx context.Context, topic string, event interface{}) error
}

// NewDispatcher cria um novo dispatcher
func NewDispatcher(repo Repository, producer Producer, interval time.Duration) *Dispatcher {
	return &Dispatcher{
		repo:     repo,
		producer: producer,
		interval: interval,
	}
}

// Start inicia o dispatcher em background
func (d *Dispatcher) Start(ctx context.Context) {
	log.Info().
		Dur("interval", d.interval).
		Msg("iniciando outbox dispatcher")
	
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("outbox dispatcher interrompido")
			return
		case <-ticker.C:
			if err := d.processPending(ctx); err != nil {
				log.Error().Err(err).Msg("erro ao processar mensagens pendentes")
			}
		}
	}
}

// processPending processa mensagens pendentes da outbox
func (d *Dispatcher) processPending(ctx context.Context) error {
	messages, err := d.repo.GetPending(ctx, 100) // Processa até 100 mensagens por vez
	if err != nil {
		return fmt.Errorf("erro ao buscar mensagens pendentes: %w", err)
	}
	
	if len(messages) == 0 {
		return nil
	}
	
	log.Info().
		Int("count", len(messages)).
		Msg("processando mensagens da outbox")
	
	for _, message := range messages {
		if err := d.processMessage(ctx, message); err != nil {
			log.Error().
				Err(err).
				Uint("message_id", message.ID).
				Str("event_type", message.EventType).
				Msg("erro ao processar mensagem da outbox")
			continue
		}
		
		// Marca como processada
		if err := d.repo.MarkAsProcessed(ctx, message.ID); err != nil {
			log.Error().
				Err(err).
				Uint("message_id", message.ID).
				Msg("erro ao marcar mensagem como processada")
		}
	}
	
	return nil
}

// processMessage processa uma mensagem individual
func (d *Dispatcher) processMessage(ctx context.Context, message OutboxMessage) error {
	// Determina o tópico baseado no tipo de evento
	topic := d.getTopicForEvent(message.EventType)
	
	// Deserializa o payload
	var event interface{}
	if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
		return fmt.Errorf("erro ao deserializar payload: %w", err)
	}
	
	// Publica no Kafka
	if err := d.producer.PublishEvent(ctx, topic, event); err != nil {
		return fmt.Errorf("erro ao publicar evento: %w", err)
	}
	
	log.Info().
		Uint("message_id", message.ID).
		Str("topic", topic).
		Str("event_type", message.EventType).
		Msg("mensagem da outbox publicada com sucesso")
	
	return nil
}

// getTopicForEvent mapeia tipos de evento para tópicos
func (d *Dispatcher) getTopicForEvent(eventType string) string {
	// Mapeamento direto: user.created -> user.created
	return eventType
}

// CreateOutboxMessage cria uma nova mensagem de outbox
func CreateOutboxMessage(aggregate, eventType string, payload interface{}) (*OutboxMessage, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %w", err)
	}
	
	return &OutboxMessage{
		Aggregate: aggregate,
		EventType: eventType,
		Payload:   string(payloadBytes),
		CreatedAt: time.Now(),
	}, nil
}
