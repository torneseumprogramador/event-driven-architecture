package outbox

import (
	"context"
	"fmt"
	"product-service/internal/domain"
	pkgoutbox "pkg/outbox"
	pkgevents "pkg/events"

	"gorm.io/gorm"
)

// OutboxService serviço para gerenciar outbox
type OutboxService struct {
	db *gorm.DB
}

// NewOutboxService cria um novo serviço de outbox
func NewOutboxService(db *gorm.DB) *OutboxService {
	return &OutboxService{db: db}
}

// CreateProductWithEvent cria um produto e grava o evento na outbox na mesma transação
func (s *OutboxService) CreateProductWithEvent(ctx context.Context, product *domain.Product) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Cria o produto
		if err := tx.Create(product).Error; err != nil {
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
		
		// Cria a mensagem da outbox
		outboxMessage, err := pkgoutbox.CreateOutboxMessage("product", "product.created", event)
		if err != nil {
			return err
		}
		
		// Grava na outbox
		if err := tx.Create(outboxMessage).Error; err != nil {
			return err
		}
		
		return nil
	})
}

// UpdateProductWithEvent atualiza um produto e grava o evento na outbox na mesma transação
func (s *OutboxService) UpdateProductWithEvent(ctx context.Context, productID uint, updates map[string]interface{}) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Busca o produto atual
		var product domain.Product
		if err := tx.Where("id = ?", productID).First(&product).Error; err != nil {
			return fmt.Errorf("produto não encontrado")
		}
		
		// Aplica as atualizações
		if err := tx.Model(&product).Updates(updates).Error; err != nil {
			return err
		}
		
		// Recarrega o produto para ter os dados atualizados
		if err := tx.First(&product, productID).Error; err != nil {
			return err
		}
		
		// Cria o evento
		event := pkgevents.ProductUpdated{
			BaseEvent: pkgevents.NewBaseEvent(),
			Product: pkgevents.ProductData{
				ID:    product.ID,
				Name:  product.Name,
				Price: product.Price,
				Stock: product.Stock,
			},
		}
		
		// Cria a mensagem da outbox
		outboxMessage, err := pkgoutbox.CreateOutboxMessage("product", "product.updated", event)
		if err != nil {
			return err
		}
		
		// Grava na outbox
		if err := tx.Create(outboxMessage).Error; err != nil {
			return err
		}
		
		return nil
	})
}
