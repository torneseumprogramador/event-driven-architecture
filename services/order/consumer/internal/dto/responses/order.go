package responses

import (
	"time"
)

// OrderResponse response do pedido
type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   time.Time           `json:"created_at"`
}
