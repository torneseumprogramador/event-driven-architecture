package dto

import (
	"time"
	"product-service/internal/domain"
)

// CreateProductRequest request para criar produto
type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
	Stock int     `json:"stock" binding:"required,gte=0"`
}

// UpdateProductRequest request para atualizar produto
type UpdateProductRequest struct {
	Name  *string  `json:"name,omitempty"`
	Price *float64 `json:"price,omitempty"`
	Stock *int     `json:"stock,omitempty"`
}

// ProductResponse response do produto
type ProductResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

// ToProductResponse converte Product para ProductResponse
func ToProductResponse(product *domain.Product) *ProductResponse {
	return &ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
	}
}

// ToProductResponseList converte lista de Product para ProductResponse
func ToProductResponseList(products []domain.Product) []ProductResponse {
	var responses []ProductResponse
	for _, product := range products {
		responses = append(responses, *ToProductResponse(&product))
	}
	return responses
}
