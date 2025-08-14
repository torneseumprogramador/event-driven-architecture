package repo

import (
	"context"
	"query-api/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// OrderRepository define a interface do repositório de pedidos
type OrderRepository interface {
	FindAll(ctx context.Context) ([]entities.OrderView, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.OrderView, error)
}

// MongoOrderRepository implementa OrderRepository usando MongoDB
type MongoOrderRepository struct {
	collection *mongo.Collection
}

// NewMongoOrderRepository cria uma nova instância de MongoOrderRepository
func NewMongoOrderRepository(db *mongo.Database) OrderRepository {
	return &MongoOrderRepository{
		collection: db.Collection("views.orders"),
	}
}

// FindAll retorna todos os pedidos
func (r *MongoOrderRepository) FindAll(ctx context.Context) ([]entities.OrderView, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
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

// FindByID retorna um pedido pelo ID
func (r *MongoOrderRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.OrderView, error) {
	var order entities.OrderView
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}
