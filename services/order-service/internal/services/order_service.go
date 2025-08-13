package services

import (
	"context"
	"fmt"
	"order-service/internal/domain/entities"
	"order-service/internal/repo"
	pkgoutbox "pkg/outbox"
	pkgevents "pkg/events"

	"gorm.io/gorm"
)

// OrderService serviço para gerenciar operações de pedido
type OrderService struct {
	orderRepo repo.OrderRepository
	db        *gorm.DB // Mantido para transações
}

// NewOrderService cria um novo serviço de pedido
func NewOrderService(orderRepo repo.OrderRepository, db *gorm.DB) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		db:        db,
	}
}

// CreateOrderWithEvent cria um pedido e grava o evento na outbox na mesma transação
func (s *OrderService) CreateOrderWithEvent(ctx context.Context, order *entities.Order) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Cria o pedido usando o repositório
		if err := s.orderRepo.Create(ctx, order); err != nil {
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
		// Verifica se o pedido existe e está no status correto usando o repositório
		order, err := s.orderRepo.GetByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("pedido não encontrado")
		}
		
		if order.Status != "CREATED" {
			return fmt.Errorf("pedido não pode ser pago no status atual: %s", order.Status)
		}
		
		// Atualiza o status do pedido usando o repositório
		if err := s.orderRepo.UpdateStatus(ctx, orderID, "PAID"); err != nil {
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
		// Verifica se o pedido existe e pode ser cancelado usando o repositório
		order, err := s.orderRepo.GetByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("pedido não encontrado")
		}
		
		if order.Status == "CANCELED" {
			return fmt.Errorf("pedido já está cancelado")
		}
		
		if order.Status == "PAID" {
			return fmt.Errorf("pedido pago não pode ser cancelado")
		}
		
		// Atualiza o status do pedido usando o repositório
		if err := s.orderRepo.UpdateStatus(ctx, orderID, "CANCELED"); err != nil {
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

// CreateOrder cria um pedido sem evento
func (s *OrderService) CreateOrder(ctx context.Context, order *entities.Order) error {
	return s.orderRepo.Create(ctx, order)
}

// GetOrderByID busca pedido por ID
func (s *OrderService) GetOrderByID(ctx context.Context, id uint) (*entities.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

// GetAllOrders busca todos os pedidos
func (s *OrderService) GetAllOrders(ctx context.Context) ([]entities.Order, error) {
	return s.orderRepo.GetAll(ctx)
}

// UpdateOrder atualiza um pedido
func (s *OrderService) UpdateOrder(ctx context.Context, order *entities.Order) error {
	return s.orderRepo.Update(ctx, order)
}

// DeleteOrder remove um pedido
func (s *OrderService) DeleteOrder(ctx context.Context, id uint) error {
	return s.orderRepo.Delete(ctx, id)
}

// UpdateOrderStatus atualiza o status do pedido
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id uint, status string) error {
	return s.orderRepo.UpdateStatus(ctx, id, status)
}
