package services

import (
	"context"
	"time"
	"query-consumer/internal/domain/entities"
	"query-consumer/internal/repository"
	pkgevents "pkg/events"

	"go.mongodb.org/mongo-driver/bson"
)

// ProductService interface para business logic de produtos
type ProductService interface {
	HandleProductCreated(ctx context.Context, event pkgevents.ProductCreated) error
	HandleProductUpdated(ctx context.Context, event pkgevents.ProductUpdated) error
	HandleStockReserved(ctx context.Context, event pkgevents.StockReserved) error
	HandleStockReleased(ctx context.Context, event pkgevents.StockReleased) error
}

// ProductServiceImpl implementação do service de produtos
type ProductServiceImpl struct {
	productRepository repository.ProductRepository
}

// NewProductService cria uma nova instância do service
func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		productRepository: productRepository,
	}
}

// HandleProductCreated processa evento de produto criado
func (s *ProductServiceImpl) HandleProductCreated(ctx context.Context, event pkgevents.ProductCreated) error {
	productView := &entities.ProductProjectionView{
		ID:        int(event.Product.ID),
		Name:      event.Product.Name,
		Price:     event.Product.Price,
		Stock:     event.Product.Stock,
		CreatedAt: event.OccurredAt,
		UpdatedAt: time.Now(),
	}
	
	// Cria novo produto
	return s.productRepository.Create(ctx, productView)
}

// HandleProductUpdated processa evento de produto atualizado
func (s *ProductServiceImpl) HandleProductUpdated(ctx context.Context, event pkgevents.ProductUpdated) error {
	filter := bson.M{"_id": event.Product.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       event.Product.Name,
			"price":      event.Product.Price,
			"stock":      event.Product.Stock,
			"updated_at": time.Now(),
		},
	}
	
	return s.productRepository.Update(ctx, filter, update)
}

// HandleStockReserved processa evento de estoque reservado
func (s *ProductServiceImpl) HandleStockReserved(ctx context.Context, event pkgevents.StockReserved) error {
	filter := bson.M{"_id": event.ProductID}
	update := bson.M{
		"$inc": bson.M{
			"stock": -event.Quantity,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
	
	return s.productRepository.Update(ctx, filter, update)
}

// HandleStockReleased processa evento de estoque liberado
func (s *ProductServiceImpl) HandleStockReleased(ctx context.Context, event pkgevents.StockReleased) error {
	filter := bson.M{"_id": event.ProductID}
	update := bson.M{
		"$inc": bson.M{
			"stock": event.Quantity,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
	
	return s.productRepository.Update(ctx, filter, update)
}
