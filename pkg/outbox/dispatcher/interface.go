package dispatcher

import (
	"context"
)

// Producer interface para publicação no Kafka
type Producer interface {
	PublishEvent(ctx context.Context, topic string, event interface{}) error
}

// OutboxDispatcher interface para processamento de mensagens da outbox
type OutboxDispatcher interface {
	Start(ctx context.Context)
	SetBatchSize(batchSize int)
	GetStats(ctx context.Context) (map[string]interface{}, error)
}
