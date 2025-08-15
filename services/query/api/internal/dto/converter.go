package dto

import (
	"query-api/internal/domain/entities"
	"query-api/internal/dto/responses"
)

// ToUserResponse converte UserView para UserResponse
func ToUserResponse(user entities.UserView) responses.UserResponse {
	return responses.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUsersResponse converte slice de UserView para UsersResponse
func ToUsersResponse(users []entities.UserView) responses.UsersResponse {
	userResponses := make([]responses.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}
	
	return responses.UsersResponse{
		Users: userResponses,
		Total: len(users),
	}
}

// ToProductResponse converte ProductView para ProductResponse
func ToProductResponse(product entities.ProductView) responses.ProductResponse {
	return responses.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// ToProductsResponse converte slice de ProductView para ProductsResponse
func ToProductsResponse(products []entities.ProductView) responses.ProductsResponse {
	productResponses := make([]responses.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = ToProductResponse(product)
	}
	
	return responses.ProductsResponse{
		Products: productResponses,
		Total:    len(products),
	}
}

// ToOrderItemResponse converte OrderItemView para OrderItemResponse
func ToOrderItemResponse(item entities.OrderItemView) responses.OrderItemResponse {
	return responses.OrderItemResponse{
		ProductID:   item.ProductID,
		ProductName: item.ProductName,
		Quantity:    item.Quantity,
		Price:       item.Price,
		Total:       item.Total,
	}
}

// ToOrderResponse converte OrderView para OrderResponse
func ToOrderResponse(order entities.OrderView) responses.OrderResponse {
	itemResponses := make([]responses.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		itemResponses[i] = ToOrderItemResponse(item)
	}
	
	return responses.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      order.Status,
		Total:       order.Total,
		User: responses.UserResponse{
			ID:    order.User.ID,
			Name:  order.User.Name,
			Email: order.User.Email,
		},
		Items:       itemResponses,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		PaidAt:      order.PaidAt,
		CanceledAt:  order.CanceledAt,
	}
}

// ToOrdersResponse converte slice de OrderView para OrdersResponse
func ToOrdersResponse(orders []entities.OrderView) responses.OrdersResponse {
	orderResponses := make([]responses.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = ToOrderResponse(order)
	}
	
	return responses.OrdersResponse{
		Orders: orderResponses,
		Total:  len(orders),
	}
}
