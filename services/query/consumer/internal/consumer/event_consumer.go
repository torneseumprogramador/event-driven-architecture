package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"query-consumer/internal/projections"
	pkgkafka "pkg/kafka"
	pkgevents "pkg/events"
	pkgidempotency "pkg/idempotency"

	"github.com/rs/zerolog/log"
)

// EventConsumer consumidor de eventos para o query-service
type EventConsumer struct {
	orderProjection   *projections.OrderProjection
	productProjection *projections.ProductProjection
	userProjection    *projections.UserProjection
	kafkaProducer     *pkgkafka.Producer
	idempotencyHandler *pkgidempotency.Handler
}

// NewEventConsumer cria um novo consumidor de eventos
func NewEventConsumer(
	orderProjection *projections.OrderProjection,
	productProjection *projections.ProductProjection,
	userProjection *projections.UserProjection,
	kafkaProducer *pkgkafka.Producer,
	idempotencyHandler *pkgidempotency.Handler,
) *EventConsumer {
	return &EventConsumer{
		orderProjection:   orderProjection,
		productProjection: productProjection,
		userProjection:    userProjection,
		kafkaProducer:     kafkaProducer,
		idempotencyHandler: idempotencyHandler,
	}
}

// HandleUserCreated processa evento de usuário criado
func (c *EventConsumer) HandleUserCreated(ctx context.Context, message []byte) error {
	var event pkgevents.UserCreated
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("user_id", event.User.ID).
			Msg("processando evento user.created")
		
		// Atualiza projeção de usuário
		if err := c.userProjection.HandleUserCreated(ctx, event); err != nil {
			return fmt.Errorf("erro ao processar user.created: %w", err)
		}
		
		// Atualiza projeção de pedido (para incluir dados do usuário)
		if err := c.orderProjection.HandleUserCreated(ctx, event); err != nil {
			return fmt.Errorf("erro ao atualizar pedidos com dados do usuário: %w", err)
		}
		
		return nil
	})
}

// HandleProductCreated processa evento de produto criado
func (c *EventConsumer) HandleProductCreated(ctx context.Context, message []byte) error {
	var event pkgevents.ProductCreated
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("product_id", event.Product.ID).
			Msg("processando evento product.created")
		
		// Atualiza projeção de produto
		if err := c.productProjection.HandleProductCreated(ctx, event); err != nil {
			return fmt.Errorf("erro ao processar product.created: %w", err)
		}
		
		// Atualiza projeção de pedido (para incluir dados do produto)
		if err := c.orderProjection.HandleProductCreated(ctx, event); err != nil {
			return fmt.Errorf("erro ao atualizar pedidos com dados do produto: %w", err)
		}
		
		return nil
	})
}

// HandleProductUpdated processa evento de produto atualizado
func (c *EventConsumer) HandleProductUpdated(ctx context.Context, message []byte) error {
	var event pkgevents.ProductUpdated
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("product_id", event.Product.ID).
			Msg("processando evento product.updated")
		
		// Atualiza projeção de produto
		if err := c.productProjection.HandleProductUpdated(ctx, event); err != nil {
			return fmt.Errorf("erro ao processar product.updated: %w", err)
		}
		
		// Atualiza projeção de pedido (para incluir dados do produto)
		if err := c.orderProjection.HandleProductUpdated(ctx, event); err != nil {
			return fmt.Errorf("erro ao atualizar pedidos com dados do produto: %w", err)
		}
		
		return nil
	})
}

// HandleOrderCreated processa evento de pedido criado
func (c *EventConsumer) HandleOrderCreated(ctx context.Context, message []byte) error {
	var event pkgevents.OrderCreated
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("order_id", event.Order.ID).
			Msg("processando evento order.created")
		
		// Busca informações do usuário
		user, err := c.userProjection.GetByID(ctx, int(event.Order.UserID))
		if err != nil {
			log.Warn().Err(err).Uint("user_id", event.Order.UserID).Msg("usuário não encontrado, continuando sem dados do usuário")
		}
		
		// Busca informações dos produtos
		productInfos := make(map[int]*projections.ProductProjectionView)
		for _, item := range event.Order.Items {
			product, err := c.productProjection.GetByID(ctx, int(item.ProductID))
			if err != nil {
				log.Warn().Err(err).Uint("product_id", item.ProductID).Msg("produto não encontrado, continuando sem dados do produto")
			} else {
				productInfos[int(item.ProductID)] = product
			}
		}
		
		// Atualiza projeção de pedido com dados completos
		if err := c.orderProjection.HandleOrderCreatedWithData(ctx, event, user, productInfos); err != nil {
			return fmt.Errorf("erro ao processar order.created: %w", err)
		}
		
		return nil
	})
}

// HandleOrderPaid processa evento de pedido pago
func (c *EventConsumer) HandleOrderPaid(ctx context.Context, message []byte) error {
	var event pkgevents.OrderPaid
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("order_id", event.OrderID).
			Msg("processando evento order.paid")
		
		// Atualiza projeção de pedido
		if err := c.orderProjection.HandleOrderPaid(ctx, event); err != nil {
			return fmt.Errorf("erro ao processar order.paid: %w", err)
		}
		
		return nil
	})
}

// HandleOrderCanceled processa evento de pedido cancelado
func (c *EventConsumer) HandleOrderCanceled(ctx context.Context, message []byte) error {
	var event pkgevents.OrderCanceled
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("order_id", event.OrderID).
			Msg("processando evento order.canceled")
		
		// Atualiza projeção de pedido
		if err := c.orderProjection.HandleOrderCanceled(ctx, event); err != nil {
			return fmt.Errorf("erro ao processar order.canceled: %w", err)
		}
		
		return nil
	})
}

// HandleStockReserved processa evento de estoque reservado
func (c *EventConsumer) HandleStockReserved(ctx context.Context, message []byte) error {
	var event pkgevents.StockReserved
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("product_id", event.ProductID).
			Int("quantity", event.Quantity).
			Msg("processando evento stock.reserved")
		
		// Atualiza projeção de produto
		if err := c.productProjection.HandleStockReserved(ctx, event); err != nil {
			return fmt.Errorf("erro ao processar stock.reserved: %w", err)
		}
		
		return nil
	})
}

// HandleStockReleased processa evento de estoque liberado
func (c *EventConsumer) HandleStockReleased(ctx context.Context, message []byte) error {
	var event pkgevents.StockReleased
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("erro ao deserializar evento: %w", err)
	}
	
	return c.idempotencyHandler.ProcessWithIdempotency(ctx, event.EventID, func() error {
		log.Info().
			Str("event_id", event.EventID).
			Uint("product_id", event.ProductID).
			Int("quantity", event.Quantity).
			Msg("processando evento stock.released")
		
		// Atualiza projeção de produto
		if err := c.productProjection.HandleStockReleased(ctx, event); err != nil {
			return fmt.Errorf("erro ao processar stock.released: %w", err)
		}
		
		return nil
	})
}
