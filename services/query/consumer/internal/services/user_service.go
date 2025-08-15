package services

import (
	"context"
	"time"
	"query-consumer/internal/domain/entities"
	"query-consumer/internal/repository"
	pkgevents "pkg/events"

	"go.mongodb.org/mongo-driver/bson"
)

// UserService interface para business logic de usuários
type UserService interface {
	HandleUserCreated(ctx context.Context, event pkgevents.UserCreated) error
	HandleUserUpdated(ctx context.Context, event pkgevents.UserUpdated) error
}

// UserServiceImpl implementação do service de usuários
type UserServiceImpl struct {
	userRepository repository.UserRepository
}

// NewUserService cria uma nova instância do service
func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

// HandleUserCreated processa evento de usuário criado
func (s *UserServiceImpl) HandleUserCreated(ctx context.Context, event pkgevents.UserCreated) error {
	userView := &entities.UserProjectionView{
		ID:        int(event.User.ID),
		Name:      event.User.Name,
		Email:     event.User.Email,
		CreatedAt: event.OccurredAt,
		UpdatedAt: time.Now(),
	}
	
	// Cria novo usuário
	return s.userRepository.Create(ctx, userView)
}

// HandleUserUpdated processa evento de usuário atualizado
func (s *UserServiceImpl) HandleUserUpdated(ctx context.Context, event pkgevents.UserUpdated) error {
	filter := bson.M{"_id": event.User.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       event.User.Name,
			"email":      event.User.Email,
			"updated_at": time.Now(),
		},
	}
	
	return s.userRepository.Update(ctx, filter, update)
}
