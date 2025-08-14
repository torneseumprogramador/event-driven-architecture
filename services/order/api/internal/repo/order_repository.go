package repo

import (
	"context"
	"order-api/internal/domain/entities"

	"gorm.io/gorm"
)

// OrderRepository interface para operações de pedido
type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id uint) (*entities.Order, error)
	GetAll(ctx context.Context) ([]entities.Order, error)
	Update(ctx context.Context, order *entities.Order) error
	Delete(ctx context.Context, id uint) error
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
func (r *GormOrderRepository) Create(ctx context.Context, order *entities.Order) error {
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
func (r *GormOrderRepository) GetByID(ctx context.Context, id uint) (*entities.Order, error) {
	var order entities.Order
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
		Model(&entities.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetAll busca todos os pedidos
func (r *GormOrderRepository) GetAll(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Update atualiza um pedido
func (r *GormOrderRepository) Update(ctx context.Context, order *entities.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete remove um pedido
func (r *GormOrderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Order{}, id).Error
}
