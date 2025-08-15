package projections

import (
	"context"
	"time"
	pkgevents "pkg/events"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserProjectionView representa a projeção de usuário no MongoDB
type UserProjectionView struct {
	ID        int       `bson:"_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// UserProjection gerencia as projeções de usuário
type UserProjection struct {
	collection *mongo.Collection
}

// NewUserProjection cria uma nova projeção de usuário
func NewUserProjection(db *mongo.Database) *UserProjection {
	return &UserProjection{
		collection: db.Collection("views.users"),
	}
}

// HandleUserCreated processa evento de usuário criado
func (p *UserProjection) HandleUserCreated(ctx context.Context, event pkgevents.UserCreated) error {
	userView := UserProjectionView{
		ID:        int(event.User.ID),
		Name:      event.User.Name,
		Email:     event.User.Email,
		CreatedAt: event.OccurredAt,
		UpdatedAt: time.Now(),
	}
	
	// Usa ReplaceOne com upsert para evitar erro de chave duplicada
	filter := bson.M{"_id": event.User.ID}
	_, err := p.collection.ReplaceOne(ctx, filter, userView, options.Replace().SetUpsert(true))
	return err
}

// HandleUserUpdated processa evento de usuário atualizado
func (p *UserProjection) HandleUserUpdated(ctx context.Context, event pkgevents.UserUpdated) error {
	filter := bson.M{"_id": event.User.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       event.User.Name,
			"email":      event.User.Email,
			"updated_at": time.Now(),
		},
	}
	
	_, err := p.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetAll busca todos os usuários
func (p *UserProjection) GetAll(ctx context.Context) ([]UserProjectionView, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	
	cursor, err := p.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var users []UserProjectionView
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	
	return users, nil
}

// GetByID busca usuário por ID
func (p *UserProjection) GetByID(ctx context.Context, id int) (*UserProjectionView, error) {
	var user UserProjectionView
	err := p.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
