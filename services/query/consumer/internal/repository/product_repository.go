package repository

import (
	"context"
	"query-consumer/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductRepository interface para acesso a dados de produtos
type ProductRepository interface {
	Create(ctx context.Context, product *entities.ProductProjectionView) error
	Update(ctx context.Context, filter bson.M, update bson.M) error
	GetAll(ctx context.Context) ([]entities.ProductProjectionView, error)
	GetByID(ctx context.Context, id int) (*entities.ProductProjectionView, error)
}

// MongoProductRepository implementação do repository usando MongoDB
type MongoProductRepository struct {
	collection *mongo.Collection
}

// NewMongoProductRepository cria uma nova instância do repository
func NewMongoProductRepository(db *mongo.Database) ProductRepository {
	return &MongoProductRepository{
		collection: db.Collection("views.products"),
	}
}

// Create cria um novo produto
func (r *MongoProductRepository) Create(ctx context.Context, product *entities.ProductProjectionView) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

// Update atualiza um produto
func (r *MongoProductRepository) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetAll busca todos os produtos
func (r *MongoProductRepository) GetAll(ctx context.Context) ([]entities.ProductProjectionView, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var products []entities.ProductProjectionView
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	
	return products, nil
}

// GetByID busca produto por ID
func (r *MongoProductRepository) GetByID(ctx context.Context, id int) (*entities.ProductProjectionView, error) {
	var product entities.ProductProjectionView
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
