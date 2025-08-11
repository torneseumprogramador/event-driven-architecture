     package domain

import (
	"time"
)

// Order representa a entidade pedido
type Order struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	Status       string    `json:"status" gorm:"not null;default:'CREATED'"`
	TotalAmount  float64   `json:"total_amount" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null"`
	Items        []OrderProduct `json:"items" gorm:"foreignKey:OrderID"`
}

// OrderProduct representa um item do pedido
type OrderProduct struct {
	OrderID   uint    `json:"order_id" gorm:"primaryKey"`
	ProductID uint    `json:"product_id" gorm:"primaryKey"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	UnitPrice float64 `json:"unit_price" gorm:"not null"`
}

// TableName especifica o nome da tabela
func (OrderProduct) TableName() string {
	return "order_products"
}

// CreateOrderRequest request para criar pedido
type CreateOrderRequest struct {
	UserID uint           `json:"user_id" binding:"required"`
	Items  []OrderItemReq `json:"items" binding:"required,min=1"`
}

// OrderItemReq item do pedido no request
type OrderItemReq struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// OrderResponse response do pedido
type OrderResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Status       string    `json:"status"`
	TotalAmount  float64   `json:"total_amount"`
	CreatedAt    time.Time `json:"created_at"`
	Items        []OrderItemResponse `json:"items"`
}

// OrderItemResponse item do pedido no response
type OrderItemResponse struct {
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}

// ToResponse converte Order para OrderResponse
func (o *Order) ToResponse() OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}
	
	return OrderResponse{
		ID:          o.ID,
		UserID:      o.UserID,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		CreatedAt:   o.CreatedAt,
		Items:       items,
	}
}
