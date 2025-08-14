package dto

import (
	"product-consumer/internal/domain/entities"
	"product-consumer/internal/dto/responses"
)

// ToProductResponse converte Product para ProductResponse
func ToProductResponse(product *entities.Product) *responses.ProductResponse {
	return &responses.ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
	}
}

// ToProductResponseList converte lista de Product para ProductResponse
func ToProductResponseList(products []entities.Product) []responses.ProductResponse {
	var responses []responses.ProductResponse
	for _, product := range products {
		responses = append(responses, *ToProductResponse(&product))
	}
	return responses
}
