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

// UserCreated evento de usuário criado
type UserCreated struct {
	BaseEvent
	User UserData `json:"user"`
}

// UserData dados do usuário
type UserData struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserUpdated evento de usuário atualizado
type UserUpdated struct {
	BaseEvent
	User UserData `json:"user"`
}

// ProductCreated evento de produto criado
type ProductCreated struct {
	BaseEvent
	Product ProductData `json:"product"`
}

// ProductData dados do produto
type ProductData struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// ProductUpdated evento de produto atualizado
type ProductUpdated struct {
	BaseEvent
	Product ProductData `json:"product"`
}

// OrderCreated evento de pedido criado
type OrderCreated struct {
	BaseEvent
	Order OrderData `json:"order"`
}

// OrderData dados do pedido
type OrderData struct {
	ID           uint           `json:"id"`
	UserID       uint           `json:"user_id"`
	Status       string         `json:"status"`
	TotalAmount  float64        `json:"total_amount"`
	Items        []OrderItem    `json:"items"`
}

// OrderItem item do pedido
type OrderItem struct {
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}

// OrderPaid evento de pedido pago
type OrderPaid struct {
	BaseEvent
	OrderID uint `json:"order_id"`
}

// OrderCanceled evento de pedido cancelado
type OrderCanceled struct {
	BaseEvent
	OrderID uint   `json:"order_id"`
	Reason  string `json:"reason,omitempty"`
}

// StockReserved evento de estoque reservado
type StockReserved struct {
	BaseEvent
	OrderID   uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// StockReleased evento de estoque liberado
type StockReleased struct {
	BaseEvent
	OrderID   uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
