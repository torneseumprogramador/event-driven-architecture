package dto

import (
	"time"
	"order-service/internal/domain"
)

// CreateOrderRequest request para criar pedido
type CreateOrderRequest struct {
	UserID uint                `json:"user_id" binding:"required"`
	Items  []CreateOrderItem   `json:"items" binding:"required,min=1"`
}

// CreateOrderItem item do pedido para criação
type CreateOrderItem struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"required,gt=0"`
}

// UpdateOrderRequest request para atualizar pedido
type UpdateOrderRequest struct {
	Status *string `json:"status,omitempty"`
}

// OrderResponse response do pedido
type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   time.Time           `json:"created_at"`
}

// OrderItemResponse response do item do pedido
type OrderItemResponse struct {
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

// PayOrderRequest request para pagar pedido
type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

// CancelOrderRequest request para cancelar pedido
type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// ToOrderResponse converte Order para OrderResponse
func ToOrderResponse(order *domain.Order) *OrderResponse {
	items := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemResponse{
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return &OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Items:       items,
		CreatedAt:   order.CreatedAt,
	}
}

// ToOrderResponseList converte lista de Order para OrderResponse
func ToOrderResponseList(orders []domain.Order) []OrderResponse {
	var responses []OrderResponse
	for _, order := range orders {
		responses = append(responses, *ToOrderResponse(&order))
	}
	return responses
}
