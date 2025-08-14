package entities

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
