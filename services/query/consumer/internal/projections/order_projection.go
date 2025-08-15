package projections

import (
	"context"
	"time"
	pkgevents "pkg/events"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OrderView representa a projeção de pedido no MongoDB
type OrderView struct {
	ID           uint      `bson:"_id"`
	UserID       uint      `bson:"user_id"`
	Status       string    `bson:"status"`
	TotalAmount  float64   `bson:"total_amount"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
	User         UserView  `bson:"user"`
	Items        []OrderItemView `bson:"items"`
}

// UserView representa a projeção de usuário
type UserView struct {
	ID    uint   `bson:"id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

// OrderItemView representa a projeção de item do pedido
type OrderItemView struct {
	ProductID  uint    `bson:"product_id"`
	Quantity   int     `bson:"quantity"`
	UnitPrice  float64 `bson:"unit_price"`
	Product    ProductView `bson:"product"`
}

// ProductView representa a projeção de produto
type ProductView struct {
	ID    uint    `bson:"id"`
	Name  string  `bson:"name"`
	Price float64 `bson:"price"`
	Stock int     `bson:"stock"`
}

// OrderProjection gerencia as projeções de pedido
type OrderProjection struct {
	collection *mongo.Collection
}

// NewOrderProjection cria uma nova projeção de pedido
func NewOrderProjection(db *mongo.Database) *OrderProjection {
	return &OrderProjection{
		collection: db.Collection("views.orders"),
	}
}

// HandleUserCreated processa evento de usuário criado
func (p *OrderProjection) HandleUserCreated(ctx context.Context, event pkgevents.UserCreated) error {
	// Atualiza usuários em pedidos existentes
	filter := bson.M{"user_id": event.User.ID}
	update := bson.M{
		"$set": bson.M{
			"user": UserView{
				ID:    event.User.ID,
				Name:  event.User.Name,
				Email: event.User.Email,
			},
			"updated_at": time.Now(),
		},
	}
	
	_, err := p.collection.UpdateMany(ctx, filter, update)
	return err
}

// HandleProductCreated processa evento de produto criado
func (p *OrderProjection) HandleProductCreated(ctx context.Context, event pkgevents.ProductCreated) error {
	// Atualiza produtos em itens de pedidos existentes
	filter := bson.M{"items.product_id": event.Product.ID}
	update := bson.M{
		"$set": bson.M{
			"items.$.product": ProductView{
				ID:    event.Product.ID,
				Name:  event.Product.Name,
				Price: event.Product.Price,
				Stock: event.Product.Stock,
			},
			"updated_at": time.Now(),
		},
	}
	
	_, err := p.collection.UpdateMany(ctx, filter, update)
	return err
}

// HandleProductUpdated processa evento de produto atualizado
func (p *OrderProjection) HandleProductUpdated(ctx context.Context, event pkgevents.ProductUpdated) error {
	// Atualiza produtos em itens de pedidos existentes
	filter := bson.M{"items.product_id": event.Product.ID}
	update := bson.M{
		"$set": bson.M{
			"items.$.product": ProductView{
				ID:    event.Product.ID,
				Name:  event.Product.Name,
				Price: event.Product.Price,
				Stock: event.Product.Stock,
			},
			"updated_at": time.Now(),
		},
	}
	
	_, err := p.collection.UpdateMany(ctx, filter, update)
	return err
}

// HandleOrderCreated processa evento de pedido criado
func (p *OrderProjection) HandleOrderCreated(ctx context.Context, event pkgevents.OrderCreated) error {
	// Converte itens para o formato da projeção
	items := make([]OrderItemView, len(event.Order.Items))
	for i, item := range event.Order.Items {
		items[i] = OrderItemView{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}
	
	// Cria a projeção do pedido
	orderView := OrderView{
		ID:          event.Order.ID,
		UserID:      event.Order.UserID,
		Status:      event.Order.Status,
		TotalAmount: event.Order.TotalAmount,
		CreatedAt:   event.OccurredAt,
		UpdatedAt:   time.Now(),
		Items:       items,
	}
	
	// Usa ReplaceOne com upsert para evitar erro de chave duplicada
	filter := bson.M{"_id": event.Order.ID}
	_, err := p.collection.ReplaceOne(ctx, filter, orderView, options.Replace().SetUpsert(true))
	return err
}

// HandleOrderPaid processa evento de pedido pago
func (p *OrderProjection) HandleOrderPaid(ctx context.Context, event pkgevents.OrderPaid) error {
	filter := bson.M{"_id": event.OrderID}
	update := bson.M{
		"$set": bson.M{
			"status":     "PAID",
			"updated_at": time.Now(),
		},
	}
	
	_, err := p.collection.UpdateOne(ctx, filter, update)
	return err
}

// HandleOrderCanceled processa evento de pedido cancelado
func (p *OrderProjection) HandleOrderCanceled(ctx context.Context, event pkgevents.OrderCanceled) error {
	filter := bson.M{"_id": event.OrderID}
	update := bson.M{
		"$set": bson.M{
			"status":     "CANCELED",
			"updated_at": time.Now(),
		},
	}
	
	_, err := p.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetByID busca pedido por ID
func (p *OrderProjection) GetByID(ctx context.Context, id uint) (*OrderView, error) {
	var order OrderView
	err := p.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUser busca pedidos por usuário
func (p *OrderProjection) GetByUser(ctx context.Context, userID uint, status string) ([]OrderView, error) {
	filter := bson.M{"user_id": userID}
	if status != "" {
		filter["status"] = status
	}
	
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := p.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var orders []OrderView
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	
	return orders, nil
}
