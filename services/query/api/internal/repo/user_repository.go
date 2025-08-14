package repo

import (
	"context"
	"query-api/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository define a interface do repositório de usuários
type UserRepository interface {
	FindAll(ctx context.Context) ([]entities.UserView, error)
	FindByID(ctx context.Context, id int) (*entities.UserView, error)
}

// MongoUserRepository implementa UserRepository usando MongoDB
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository cria uma nova instância de MongoUserRepository
func NewMongoUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		collection: db.Collection("views.users"),
	}
}

// FindAll retorna todos os usuários
func (r *MongoUserRepository) FindAll(ctx context.Context) ([]entities.UserView, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []entities.UserView
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// FindByID retorna um usuário pelo ID
func (r *MongoUserRepository) FindByID(ctx context.Context, id int) (*entities.UserView, error) {
	var user entities.UserView
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
