package repository

import (
	"context"
	"query-consumer/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OrderRepository interface para acesso a dados de pedidos
type OrderRepository interface {
	Create(ctx context.Context, order *entities.OrderView) error
	Update(ctx context.Context, filter bson.M, update bson.M) error
	UpdateMany(ctx context.Context, filter bson.M, update bson.M) error
	GetByID(ctx context.Context, id int) (*entities.OrderView, error)
	GetByUser(ctx context.Context, userID int, status string) ([]entities.OrderView, error)
}

// MongoOrderRepository implementação do repository usando MongoDB
type MongoOrderRepository struct {
	collection *mongo.Collection
}

// NewMongoOrderRepository cria uma nova instância do repository
func NewMongoOrderRepository(db *mongo.Database) OrderRepository {
	return &MongoOrderRepository{
		collection: db.Collection("views.orders"),
	}
}

// Create cria um novo pedido
func (r *MongoOrderRepository) Create(ctx context.Context, order *entities.OrderView) error {
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

// Update atualiza um pedido
func (r *MongoOrderRepository) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateMany atualiza múltiplos pedidos
func (r *MongoOrderRepository) UpdateMany(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// GetByID busca pedido por ID
func (r *MongoOrderRepository) GetByID(ctx context.Context, id int) (*entities.OrderView, error) {
	var order entities.OrderView
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUser busca pedidos por usuário
func (r *MongoOrderRepository) GetByUser(ctx context.Context, userID int, status string) ([]entities.OrderView, error) {
	filter := bson.M{"user_id": userID}
	if status != "" {
		filter["status"] = status
	}
	
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var orders []entities.OrderView
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	
	return orders, nil
}
