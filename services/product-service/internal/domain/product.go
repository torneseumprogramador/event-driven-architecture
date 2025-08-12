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
