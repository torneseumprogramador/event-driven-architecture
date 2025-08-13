package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pkg/outbox/entities"
	"pkg/outbox/services"

	"github.com/rs/zerolog/log"
)

// Producer interface para publicação no Kafka
type Producer interface {
	PublishEvent(ctx context.Context, topic string, event interface{}) error
}

// OutboxDispatcher processa mensagens da outbox e as publica no Kafka
type OutboxDispatcher struct {
	outboxService *services.OutboxService
	producer      Producer
	interval      time.Duration
	batchSize     int
}

// NewOutboxDispatcher cria um novo dispatcher
func NewOutboxDispatcher(outboxService *services.OutboxService, producer Producer, interval time.Duration) *OutboxDispatcher {
	return &OutboxDispatcher{
		outboxService: outboxService,
		producer:      producer,
		interval:      interval,
		batchSize:     100, // Processa até 100 mensagens por vez
	}
}

// Start inicia o dispatcher em background
func (d *OutboxDispatcher) Start(ctx context.Context) {
	log.Info().
		Dur("interval", d.interval).
		Int("batch_size", d.batchSize).
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
func (d *OutboxDispatcher) processPending(ctx context.Context) error {
	messages, err := d.outboxService.GetPendingMessages(ctx, d.batchSize)
	if err != nil {
		return fmt.Errorf("erro ao buscar mensagens pendentes: %w", err)
	}

	if len(messages) == 0 {
		log.Debug().Msg("nenhuma mensagem pendente encontrada na outbox")
		return nil
	}

	log.Info().
		Int("count", len(messages)).
		Msg("processando mensagens da outbox")

	processedCount := 0
	failedCount := 0

	for _, message := range messages {
		if err := d.processMessage(ctx, message); err != nil {
			log.Error().
				Err(err).
				Uint("message_id", message.ID).
				Str("event_type", message.EventType).
				Msg("erro ao processar mensagem da outbox")
			failedCount++
			continue
		}

		// Marca como processada
		if err := d.outboxService.MarkMessageAsProcessed(ctx, message.ID); err != nil {
			log.Error().
				Err(err).
				Uint("message_id", message.ID).
				Msg("erro ao marcar mensagem como processada")
			failedCount++
			continue
		}

		processedCount++
	}

	log.Info().
		Int("processed", processedCount).
		Int("failed", failedCount).
		Int("total", len(messages)).
		Msg("processamento de mensagens concluído")

	return nil
}

// processMessage processa uma mensagem individual
func (d *OutboxDispatcher) processMessage(ctx context.Context, message entities.OutboxMessage) error {
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
func (d *OutboxDispatcher) getTopicForEvent(eventType string) string {
	// Mapeamento direto: user.created -> user.created
	return eventType
}

// SetBatchSize define o tamanho do lote de processamento
func (d *OutboxDispatcher) SetBatchSize(batchSize int) {
	d.batchSize = batchSize
}

// GetStats retorna estatísticas do dispatcher
func (d *OutboxDispatcher) GetStats(ctx context.Context) (map[string]interface{}, error) {
	pendingCount, err := d.outboxService.GetPendingCount(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"batch_size":    d.batchSize,
		"interval":      d.interval.String(),
		"pending_count": pendingCount,
	}, nil
}
