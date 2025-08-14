package dto

import (
	"order-api/internal/domain/entities"
	"order-api/internal/dto/responses"
)

// ToOrderResponse converte Order para OrderResponse
func ToOrderResponse(order *entities.Order) *responses.OrderResponse {
	items := make([]responses.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = responses.OrderItemResponse{
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return &responses.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Items:       items,
		CreatedAt:   order.CreatedAt,
	}
}

// ToOrderResponseList converte lista de Order para OrderResponse
func ToOrderResponseList(orders []entities.Order) []responses.OrderResponse {
	var responses []responses.OrderResponse
	for _, order := range orders {
		responses = append(responses, *ToOrderResponse(&order))
	}
	return responses
}
