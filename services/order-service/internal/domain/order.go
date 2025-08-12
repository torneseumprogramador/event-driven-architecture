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


