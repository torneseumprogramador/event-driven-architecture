package services

import (
	"context"
	"fmt"
	"user-service/internal/domain/entities"
	"user-service/internal/repo"
	pkgoutboxservices "pkg/outbox/services"
	pkgevents "pkg/events"

	"gorm.io/gorm"
)

// UserService serviço para gerenciar operações de usuário
type UserService struct {
	userRepo      repo.UserRepository
	outboxService pkgoutboxservices.OutboxService
	db            *gorm.DB // Mantido para transações
}

// NewUserService cria um novo serviço de usuário
func NewUserService(userRepo repo.UserRepository, outboxService pkgoutboxservices.OutboxService, db *gorm.DB) *UserService {
	return &UserService{
		userRepo:      userRepo,
		outboxService: outboxService,
		db:            db,
	}
}

// CreateUserWithEvent cria um usuário e grava o evento na outbox na mesma transação
func (s *UserService) CreateUserWithEvent(ctx context.Context, user *entities.User) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Verifica se email já existe usando o repositório
		existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
		if err == nil && existingUser != nil {
			return fmt.Errorf("email já cadastrado")
		}
		
		// Cria o usuário usando o repositório
		if err := s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		
		// Cria o evento
		event := pkgevents.UserCreated{
			BaseEvent: pkgevents.NewBaseEvent(),
			User: pkgevents.UserData{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
		}
		
		// Cria a mensagem da outbox usando o serviço dentro da transação
		_, err = s.outboxService.CreateMessageInTransaction(ctx, tx, "user", "user.created", event)
		if err != nil {
			return err
		}
		
		return nil
	})
}

// CreateUser cria um usuário sem evento
func (s *UserService) CreateUser(ctx context.Context, user *entities.User) error {
	// Verifica se email já existe
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("email já cadastrado")
	}
	
	return s.userRepo.Create(ctx, user)
}

// GetUserByID busca usuário por ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail busca usuário por email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetAllUsers busca todos os usuários
func (s *UserService) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	return s.userRepo.GetAll(ctx)
}

// UpdateUser atualiza um usuário
func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser remove um usuário
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}
