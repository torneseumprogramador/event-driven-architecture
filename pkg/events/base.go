package events

import (
	"time"

	"github.com/google/uuid"
)

// BaseEvent estrutura base para todos os eventos
type BaseEvent struct {
	EventID    string    `json:"event_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

// NewBaseEvent cria um novo evento base
func NewBaseEvent() BaseEvent {
	return BaseEvent{
		EventID:    uuid.New().String(),
		OccurredAt: time.Now(),
	}
}
