package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"product-service/internal/repo"
	pkgkafka "pkg/kafka"
	pkgevents "pkg/events"
	pkgidempotency "pkg/idempotency"

	"github.com/rs/zerolog/log"
)

// OrderConsumer consumidor de eventos de pedido
type OrderConsumer struct {
	productRepo    repo.ProductRepository
	kafkaProducer  *pkgkafka.Producer
	idempotencyHandler *pkgidempotency.Handler
}

// NewOrderConsumer cria um novo consumidor de pedidos
func NewOrderConsumer(productRepo repo.ProductRepository, kafkaProducer *pkgkafka.Producer, idempotencyHandler *pkgidempotency.Handler) *OrderConsumer {
	return &OrderConsumer{
		productRepo:    productRepo,
		kafkaProducer:  kafkaProducer,
		idempotencyHandler: idempotencyHandler,
	}
}

// HandleOrderCreated processa evento de pedido criado
func (c *OrderConsumer) HandleOrderCreated(ctx context.Context, message []byte) error {
	var event pkgevents.OrderCreated
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("order_id", event.Order.ID).
			Msg("processando evento order.created")
		
		// Processa cada item do pedido
		for _, item := range event.Order.Items {
			// Tenta reservar estoque
			if err := c.productRepo.ReserveStock(ctx, item.ProductID, item.Quantity); err != nil {
				log.Error().
					Err(err).
					Uint("product_id", item.ProductID).
					Int("quantity", item.Quantity).
					Msg("erro ao reservar estoque")
				
				// Publica evento de cancelamento
				cancelEvent := pkgevents.OrderCanceled{
					BaseEvent: pkgevents.NewBaseEvent(),
					OrderID:   event.Order.ID,
					Reason:    fmt.Sprintf("Estoque insuficiente para produto %d", item.ProductID),
				}
				
				if err := c.kafkaProducer.PublishEvent(ctx, "order.canceled", cancelEvent); err != nil {
					log.Error().Err(err).Msg("erro ao publicar evento de cancelamento")
				}
				
				return err
			}
			
			// Publica evento de estoque reservado
			stockEvent := pkgevents.StockReserved{
				BaseEvent: pkgevents.NewBaseEvent(),
				OrderID:   event.Order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}
			
			if err := c.kafkaProducer.PublishEvent(ctx, "stock.reserved", stockEvent); err != nil {
				log.Error().Err(err).Msg("erro ao publicar evento de estoque reservado")
				return err
			}
			
			log.Info().
				Uint("product_id", item.ProductID).
				Int("quantity", item.Quantity).
				Msg("estoque reservado com sucesso")
		}
		
		return nil
	})
}

// HandleOrderCanceled processa evento de pedido cancelado
func (c *OrderConsumer) HandleOrderCanceled(ctx context.Context, message []byte) error {
	var event pkgevents.OrderCanceled
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("order_id", event.OrderID).
			Msg("processando evento order.canceled")
		
		// Nota: Em um cenário real, você precisaria buscar os itens do pedido
		// para liberar o estoque. Por simplicidade, vamos apenas logar o evento.
		// Em uma implementação completa, você faria uma chamada para o order-service
		// ou teria acesso aos dados do pedido através de uma projeção.
		
		return nil
	})
}
