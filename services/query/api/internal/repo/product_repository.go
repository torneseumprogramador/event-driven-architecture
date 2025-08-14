package repo

import (
	"context"
	"query-api/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProductRepository define a interface do repositório de produtos
type ProductRepository interface {
	FindAll(ctx context.Context) ([]entities.ProductView, error)
	FindByID(ctx context.Context, id int) (*entities.ProductView, error)
}

// MongoProductRepository implementa ProductRepository usando MongoDB
type MongoProductRepository struct {
	collection *mongo.Collection
}

// NewMongoProductRepository cria uma nova instância de MongoProductRepository
func NewMongoProductRepository(db *mongo.Database) ProductRepository {
	return &MongoProductRepository{
		collection: db.Collection("views.products"),
	}
}

// FindAll retorna todos os produtos
func (r *MongoProductRepository) FindAll(ctx context.Context) ([]entities.ProductView, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []entities.ProductView
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// FindByID retorna um produto pelo ID
func (r *MongoProductRepository) FindByID(ctx context.Context, id int) (*entities.ProductView, error) {
	var product entities.ProductView
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}
