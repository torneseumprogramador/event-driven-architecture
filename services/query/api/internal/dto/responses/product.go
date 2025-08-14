package responses

import (
	"time"
)

// ProductResponse representa a resposta de produto
type ProductResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductsResponse representa a resposta de lista de produtos
type ProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
}
