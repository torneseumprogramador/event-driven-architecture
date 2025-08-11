package repo

import (
	"context"
	"product-service/internal/domain"

	"gorm.io/gorm"
)

// ProductRepository interface para operações de produto
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uint) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	ReserveStock(ctx context.Context, productID uint, quantity int) error
	ReleaseStock(ctx context.Context, productID uint, quantity int) error
}

// GormProductRepository implementação usando GORM
type GormProductRepository struct {
	db *gorm.DB
}

// NewGormProductRepository cria um novo repositório GORM
func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{db: db}
}

// Create cria um novo produto
func (r *GormProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID busca produto por ID
func (r *GormProductRepository) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update atualiza um produto
func (r *GormProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// ReserveStock reserva estoque de forma atômica
func (r *GormProductRepository) ReserveStock(ctx context.Context, productID uint, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product domain.Product
		if err := tx.Where("id = ?", productID).First(&product).Error; err != nil {
			return err
		}
		
		if err := product.ReserveStock(quantity); err != nil {
			return err
		}
		
		return tx.Save(&product).Error
	})
}

// ReleaseStock libera estoque de forma atômica
func (r *GormProductRepository) ReleaseStock(ctx context.Context, productID uint, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product domain.Product
		if err := tx.Where("id = ?", productID).First(&product).Error; err != nil {
			return err
		}
		
		product.ReleaseStock(quantity)
		return tx.Save(&product).Error
	})
}
