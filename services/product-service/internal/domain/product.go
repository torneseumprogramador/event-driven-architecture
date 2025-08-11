package domain

import (
	"fmt"
	"time"
)

// Product representa a entidade produto
type Product struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	Stock     int       `json:"stock" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

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

// ToResponse converte Product para ProductResponse
func (p *Product) ToResponse() ProductResponse {
	return ProductResponse{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
	}
}

// ReserveStock reserva estoque do produto
func (p *Product) ReserveStock(quantity int) error {
	if p.Stock < quantity {
		return fmt.Errorf("estoque insuficiente: disponível %d, solicitado %d", p.Stock, quantity)
	}
	p.Stock -= quantity
	return nil
}

// ReleaseStock libera estoque do produto
func (p *Product) ReleaseStock(quantity int) {
	p.Stock += quantity
}
