package services

import (
	"context"
	"fmt"
	"order-service/internal/domain/entities"
	pkgoutbox "pkg/outbox"
	pkgevents "pkg/events"

	"gorm.io/gorm"
)

// OrderService serviço para gerenciar operações de pedido
type OrderService struct {
	db *gorm.DB
}

// NewOrderService cria um novo serviço de pedido
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

// CreateOrderWithEvent cria um pedido e grava o evento na outbox na mesma transação
func (s *OrderService) CreateOrderWithEvent(ctx context.Context, order *entities.Order) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Cria o pedido
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		
		// Cria os itens do pedido
		for i := range order.Items {
			order.Items[i].OrderID = order.ID
		}
		
		if err := tx.Create(&order.Items).Error; err != nil {
			return err
		}
		
		// Converte itens para o formato do evento
		eventItems := make([]pkgevents.OrderItem, len(order.Items))
		for i, item := range order.Items {
			eventItems[i] = pkgevents.OrderItem{
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				UnitPrice:  item.UnitPrice,
			}
		}
		
		// Cria o evento
		event := pkgevents.OrderCreated{
			BaseEvent: pkgevents.NewBaseEvent(),
			Order: pkgevents.OrderData{
				ID:          order.ID,
				UserID:      order.UserID,
				Status:      order.Status,
				TotalAmount: order.TotalAmount,
				Items:       eventItems,
			},
		}
		
		// Cria a mensagem da outbox
		outboxMessage, err := pkgoutbox.CreateOutboxMessage("order", "order.created", event)
		if err != nil {
			return err
		}
		
		// Grava na outbox
		if err := tx.Create(outboxMessage).Error; err != nil {
			return err
		}
		
		return nil
	})
}

// PayOrderWithEvent paga um pedido e grava o evento na outbox na mesma transação
func (s *OrderService) PayOrderWithEvent(ctx context.Context, orderID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Verifica se o pedido existe e está no status correto
		var order entities.Order
		if err := tx.Where("id = ?", orderID).First(&order).Error; err != nil {
			return fmt.Errorf("pedido não encontrado")
		}
		
		if order.Status != "CREATED" {
			return fmt.Errorf("pedido não pode ser pago no status atual: %s", order.Status)
		}
		
		// Atualiza o status do pedido
		if err := tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("status", "PAID").Error; err != nil {
			return err
		}
		
		// Cria o evento
		event := pkgevents.OrderPaid{
			BaseEvent: pkgevents.NewBaseEvent(),
			OrderID:   orderID,
		}
		
		// Cria a mensagem da outbox
		outboxMessage, err := pkgoutbox.CreateOutboxMessage("order", "order.paid", event)
		if err != nil {
			return err
		}
		
		// Grava na outbox
		if err := tx.Create(outboxMessage).Error; err != nil {
			return err
		}
		
		return nil
	})
}

// CancelOrderWithEvent cancela um pedido e grava o evento na outbox na mesma transação
func (s *OrderService) CancelOrderWithEvent(ctx context.Context, orderID uint, reason string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Verifica se o pedido existe e pode ser cancelado
		var order entities.Order
		if err := tx.Where("id = ?", orderID).First(&order).Error; err != nil {
			return fmt.Errorf("pedido não encontrado")
		}
		
		if order.Status == "CANCELED" {
			return fmt.Errorf("pedido já está cancelado")
		}
		
		if order.Status == "PAID" {
			return fmt.Errorf("pedido pago não pode ser cancelado")
		}
		
		// Atualiza o status do pedido
		if err := tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("status", "CANCELED").Error; err != nil {
			return err
		}
		
		// Cria o evento
		event := pkgevents.OrderCanceled{
			BaseEvent: pkgevents.NewBaseEvent(),
			OrderID:   orderID,
			Reason:    reason,
		}
		
		// Cria a mensagem da outbox
		outboxMessage, err := pkgoutbox.CreateOutboxMessage("order", "order.canceled", event)
		if err != nil {
			return err
		}
		
		// Grava na outbox
		if err := tx.Create(outboxMessage).Error; err != nil {
			return err
		}
		
		return nil
	})
}
