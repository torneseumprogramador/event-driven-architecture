package services

import (
	"context"
	"time"
	"query-consumer/internal/domain/entities"
	"query-consumer/internal/repository"
	pkgevents "pkg/events"

	"go.mongodb.org/mongo-driver/bson"
)

// OrderService interface para business logic de pedidos
type OrderService interface {
	HandleUserCreated(ctx context.Context, event pkgevents.UserCreated) error
	HandleProductCreated(ctx context.Context, event pkgevents.ProductCreated) error
	HandleProductUpdated(ctx context.Context, event pkgevents.ProductUpdated) error
	HandleOrderCreated(ctx context.Context, event pkgevents.OrderCreated) error
	HandleOrderCreatedWithData(ctx context.Context, event pkgevents.OrderCreated, user *entities.UserProjectionView, productInfos map[int]*entities.ProductProjectionView) error
	HandleOrderPaid(ctx context.Context, event pkgevents.OrderPaid) error
	HandleOrderCanceled(ctx context.Context, event pkgevents.OrderCanceled) error
}

// OrderServiceImpl implementação do service de pedidos
type OrderServiceImpl struct {
	orderRepository repository.OrderRepository
}

// NewOrderService cria uma nova instância do service
func NewOrderService(orderRepository repository.OrderRepository) OrderService {
	return &OrderServiceImpl{
		orderRepository: orderRepository,
	}
}

// HandleUserCreated processa evento de usuário criado
func (s *OrderServiceImpl) HandleUserCreated(ctx context.Context, event pkgevents.UserCreated) error {
	// Atualiza usuários em pedidos existentes
	filter := bson.M{"user_id": event.User.ID}
	update := bson.M{
		"$set": bson.M{
			"user": entities.UserView{
				ID:    int(event.User.ID),
				Name:  event.User.Name,
				Email: event.User.Email,
			},
			"updated_at": time.Now(),
		},
	}
	
	return s.orderRepository.UpdateMany(ctx, filter, update)
}

// HandleProductCreated processa evento de produto criado
func (s *OrderServiceImpl) HandleProductCreated(ctx context.Context, event pkgevents.ProductCreated) error {
	// Atualiza produtos em itens de pedidos existentes
	filter := bson.M{"items.product_id": event.Product.ID}
	update := bson.M{
		"$set": bson.M{
			"items.$.product": entities.ProductView{
				ID:    int(event.Product.ID),
				Name:  event.Product.Name,
				Price: event.Product.Price,
				Stock: event.Product.Stock,
			},
			"updated_at": time.Now(),
		},
	}
	
	return s.orderRepository.UpdateMany(ctx, filter, update)
}

// HandleProductUpdated processa evento de produto atualizado
func (s *OrderServiceImpl) HandleProductUpdated(ctx context.Context, event pkgevents.ProductUpdated) error {
	// Atualiza produtos em itens de pedidos existentes
	filter := bson.M{"items.product_id": event.Product.ID}
	update := bson.M{
		"$set": bson.M{
			"items.$.product": entities.ProductView{
				ID:    int(event.Product.ID),
				Name:  event.Product.Name,
				Price: event.Product.Price,
				Stock: event.Product.Stock,
			},
			"updated_at": time.Now(),
		},
	}
	
	return s.orderRepository.UpdateMany(ctx, filter, update)
}

// HandleOrderCreated processa evento de pedido criado
func (s *OrderServiceImpl) HandleOrderCreated(ctx context.Context, event pkgevents.OrderCreated) error {
	// Converte itens para o formato da projeção
	items := make([]entities.OrderItemView, len(event.Order.Items))
	for i, item := range event.Order.Items {
		items[i] = entities.OrderItemView{
			ProductID: int(item.ProductID),
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}
	
	// Cria a projeção do pedido
	orderView := &entities.OrderView{
		ID:          int(event.Order.ID),
		UserID:      int(event.Order.UserID),
		Status:      event.Order.Status,
		TotalAmount: event.Order.TotalAmount,
		CreatedAt:   event.OccurredAt,
		UpdatedAt:   time.Now(),
		Items:       items,
	}
	
	// Cria novo pedido
	return s.orderRepository.Create(ctx, orderView)
}

// HandleOrderCreatedWithData processa evento de pedido criado com dados completos
func (s *OrderServiceImpl) HandleOrderCreatedWithData(ctx context.Context, event pkgevents.OrderCreated, user *entities.UserProjectionView, productInfos map[int]*entities.ProductProjectionView) error {
	// Converte itens para o formato da projeção com dados completos
	items := make([]entities.OrderItemView, len(event.Order.Items))
	for i, item := range event.Order.Items {
		productInfo := productInfos[int(item.ProductID)]
		var productName string
		var productPrice float64
		var productStock int
		
		if productInfo != nil {
			productName = productInfo.Name
			productPrice = productInfo.Price
			productStock = productInfo.Stock
		}
		
		items[i] = entities.OrderItemView{
			ProductID:  int(item.ProductID),
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			Product: entities.ProductView{
				ID:    int(item.ProductID),
				Name:  productName,
				Price: productPrice,
				Stock: productStock,
			},
		}
	}
	
	// Prepara dados do usuário
	var userName string
	var userEmail string
	
	if user != nil {
		userName = user.Name
		userEmail = user.Email
	}
	
	// Cria a projeção do pedido com dados completos
	orderView := &entities.OrderView{
		ID:          int(event.Order.ID),
		UserID:      int(event.Order.UserID),
		Status:      event.Order.Status,
		TotalAmount: event.Order.TotalAmount,
		CreatedAt:   event.OccurredAt,
		UpdatedAt:   time.Now(),
		User: entities.UserView{
			ID:    int(event.Order.UserID),
			Name:  userName,
			Email: userEmail,
		},
		Items: items,
	}
	
	// Cria novo pedido com dados completos
	return s.orderRepository.Create(ctx, orderView)
}

// HandleOrderPaid processa evento de pedido pago
func (s *OrderServiceImpl) HandleOrderPaid(ctx context.Context, event pkgevents.OrderPaid) error {
	filter := bson.M{"_id": event.OrderID}
	update := bson.M{
		"$set": bson.M{
			"status":     "PAID",
			"updated_at": time.Now(),
		},
	}
	
	return s.orderRepository.Update(ctx, filter, update)
}

// HandleOrderCanceled processa evento de pedido cancelado
func (s *OrderServiceImpl) HandleOrderCanceled(ctx context.Context, event pkgevents.OrderCanceled) error {
	filter := bson.M{"_id": event.OrderID}
	update := bson.M{
		"$set": bson.M{
			"status":     "CANCELED",
			"updated_at": time.Now(),
		},
	}
	
	return s.orderRepository.Update(ctx, filter, update)
}
