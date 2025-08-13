package entities

import (
	"database/sql"
	"time"
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

// IsProcessed verifica se a mensagem já foi processada
func (m *OutboxMessage) IsProcessed() bool {
	return m.ProcessedAt != nil
}

// MarkAsProcessed marca a mensagem como processada
func (m *OutboxMessage) MarkAsProcessed() {
	now := time.Now()
	m.ProcessedAt = &now
}
