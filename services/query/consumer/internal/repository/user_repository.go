package repository

import (
	"context"
	"query-consumer/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository interface para acesso a dados de usuários
type UserRepository interface {
	Create(ctx context.Context, user *entities.UserProjectionView) error
	Update(ctx context.Context, filter bson.M, update bson.M) error
	GetAll(ctx context.Context) ([]entities.UserProjectionView, error)
	GetByID(ctx context.Context, id int) (*entities.UserProjectionView, error)
}

// MongoUserRepository implementação do repository usando MongoDB
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository cria uma nova instância do repository
func NewMongoUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		collection: db.Collection("views.users"),
	}
}

// Create cria um novo usuário
func (r *MongoUserRepository) Create(ctx context.Context, user *entities.UserProjectionView) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Update atualiza um usuário
func (r *MongoUserRepository) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetAll busca todos os usuários
func (r *MongoUserRepository) GetAll(ctx context.Context) ([]entities.UserProjectionView, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var users []entities.UserProjectionView
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	
	return users, nil
}

// GetByID busca usuário por ID
func (r *MongoUserRepository) GetByID(ctx context.Context, id int) (*entities.UserProjectionView, error) {
	var user entities.UserProjectionView
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
