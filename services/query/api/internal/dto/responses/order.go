package responses

import (
	"time"
)



// OrderItemResponse representa um item do pedido na resposta
type OrderItemResponse struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
}

// OrderResponse representa a resposta de pedido
type OrderResponse struct {
	ID          int                 `json:"id"`
	UserID      int                 `json:"user_id"`
	Status      string              `json:"status"`
	Total       float64             `json:"total"`
	User        UserResponse        `json:"user"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	PaidAt      *time.Time          `json:"paid_at,omitempty"`
	CanceledAt  *time.Time          `json:"canceled_at,omitempty"`
}

// OrdersResponse representa a resposta de lista de pedidos
type OrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
}
