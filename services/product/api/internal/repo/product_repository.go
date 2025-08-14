package repo

import (
	"context"
	"product-api/internal/domain/entities"

	"gorm.io/gorm"
)

// ProductRepository interface para operações de produto
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id uint) (*entities.Product, error)
	GetAll(ctx context.Context) ([]entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id uint) error
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
func (r *GormProductRepository) Create(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID busca produto por ID
func (r *GormProductRepository) GetByID(ctx context.Context, id uint) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update atualiza um produto
func (r *GormProductRepository) Update(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// ReserveStock reserva estoque de forma atômica
func (r *GormProductRepository) ReserveStock(ctx context.Context, productID uint, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product entities.Product
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
		var product entities.Product
		if err := tx.Where("id = ?", productID).First(&product).Error; err != nil {
			return err
		}
		
		product.ReleaseStock(quantity)
		return tx.Save(&product).Error
	})
}

// GetAll busca todos os produtos
func (r *GormProductRepository) GetAll(ctx context.Context) ([]entities.Product, error) {
	var products []entities.Product
	err := r.db.WithContext(ctx).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Delete remove um produto
func (r *GormProductRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Product{}, id).Error
}
