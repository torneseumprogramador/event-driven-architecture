package repo

import (
	"context"
	"order-service/internal/domain"

	"gorm.io/gorm"
)

// OrderRepository interface para operações de pedido
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id uint) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

// GormOrderRepository implementação usando GORM
type GormOrderRepository struct {
	db *gorm.DB
}

// NewGormOrderRepository cria um novo repositório GORM
func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

// Create cria um novo pedido com seus itens
func (r *GormOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Cria o pedido
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		
		// Cria os itens do pedido
		for i := range order.Items {
			order.Items[i].OrderID = order.ID
		}
		
		if err := tx.Create(&order.Items).Error; err != nil {
			return err
		}
		
		return nil
	})
}

// GetByID busca pedido por ID com seus itens
func (r *GormOrderRepository) GetByID(ctx context.Context, id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateStatus atualiza o status do pedido
func (r *GormOrderRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}
