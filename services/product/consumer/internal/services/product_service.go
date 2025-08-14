package services

import (
	"context"
	"fmt"
	"product-consumer/internal/domain/entities"
	"product-consumer/internal/repo"
	pkgoutboxservices "pkg/outbox/services"
	pkgevents "pkg/events"

	"gorm.io/gorm"
)

// ProductService serviço para gerenciar operações de produto
type ProductService struct {
	productRepo   repo.ProductRepository
	outboxService pkgoutboxservices.OutboxService
	db            *gorm.DB // Mantido para transações
}

// NewProductService cria um novo serviço de produto
func NewProductService(productRepo repo.ProductRepository, outboxService pkgoutboxservices.OutboxService, db *gorm.DB) *ProductService {
	return &ProductService{
		productRepo:   productRepo,
		outboxService: outboxService,
		db:            db,
	}
}

// CreateProductWithEvent cria um produto e grava o evento na outbox na mesma transação
func (s *ProductService) CreateProductWithEvent(ctx context.Context, product *entities.Product) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Cria o produto usando o repositório
		if err := s.productRepo.Create(ctx, product); err != nil {
			return err
		}
		
		// Cria o evento
		event := pkgevents.ProductCreated{
			BaseEvent: pkgevents.NewBaseEvent(),
			Product: pkgevents.ProductData{
				ID:    product.ID,
				Name:  product.Name,
				Price: product.Price,
				Stock: product.Stock,
			},
		}
		
		// Cria a mensagem da outbox usando o serviço dentro da transação
		_, err := s.outboxService.CreateMessageInTransaction(ctx, tx, "product", "product.created", event)
		if err != nil {
			return err
		}
		
		return nil
	})
}

// UpdateProductWithEvent atualiza um produto e grava o evento na outbox na mesma transação
func (s *ProductService) UpdateProductWithEvent(ctx context.Context, productID uint, updates map[string]interface{}) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Busca o produto atual usando o repositório
		product, err := s.productRepo.GetByID(ctx, productID)
		if err != nil {
			return fmt.Errorf("produto não encontrado")
		}
		
		// Aplica as atualizações
		if err := tx.Model(product).Updates(updates).Error; err != nil {
			return err
		}
		
		// Recarrega o produto para ter os dados atualizados
		updatedProduct, err := s.productRepo.GetByID(ctx, productID)
		if err != nil {
			return err
		}
		
		// Cria o evento
		event := pkgevents.ProductUpdated{
			BaseEvent: pkgevents.NewBaseEvent(),
			Product: pkgevents.ProductData{
				ID:    updatedProduct.ID,
				Name:  updatedProduct.Name,
				Price: updatedProduct.Price,
				Stock: updatedProduct.Stock,
			},
		}
		
		// Cria a mensagem da outbox usando o serviço dentro da transação
		_, err = s.outboxService.CreateMessageInTransaction(ctx, tx, "product", "product.updated", event)
		if err != nil {
			return err
		}
		
		return nil
	})
}

// CreateProduct cria um produto sem evento
func (s *ProductService) CreateProduct(ctx context.Context, product *entities.Product) error {
	return s.productRepo.Create(ctx, product)
}

// GetProductByID busca produto por ID
func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*entities.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

// GetAllProducts busca todos os produtos
func (s *ProductService) GetAllProducts(ctx context.Context) ([]entities.Product, error) {
	return s.productRepo.GetAll(ctx)
}

// UpdateProduct atualiza um produto
func (s *ProductService) UpdateProduct(ctx context.Context, product *entities.Product) error {
	return s.productRepo.Update(ctx, product)
}

// DeleteProduct remove um produto
func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	return s.productRepo.Delete(ctx, id)
}

// ReserveStock reserva estoque
func (s *ProductService) ReserveStock(ctx context.Context, productID uint, quantity int) error {
	return s.productRepo.ReserveStock(ctx, productID, quantity)
}

// ReleaseStock libera estoque
func (s *ProductService) ReleaseStock(ctx context.Context, productID uint, quantity int) error {
	return s.productRepo.ReleaseStock(ctx, productID, quantity)
}
