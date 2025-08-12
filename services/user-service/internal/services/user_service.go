package services

import (
	"context"
	"fmt"
	"user-service/internal/domain"
	pkgoutbox "pkg/outbox"
	pkgevents "pkg/events"

	"gorm.io/gorm"
)

// UserService serviço para gerenciar operações de usuário
type UserService struct {
	db *gorm.DB
}

// NewUserService cria um novo serviço de usuário
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUserWithEvent cria um usuário e grava o evento na outbox na mesma transação
func (s *UserService) CreateUserWithEvent(ctx context.Context, user *domain.User) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Verifica se email já existe
		var existingUser domain.User
		if err := tx.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			return fmt.Errorf("email já cadastrado")
		}
		
		// Cria o usuário
		if err := tx.Create(user).Error; err != nil {
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
		
		// Cria a mensagem da outbox
		outboxMessage, err := pkgoutbox.CreateOutboxMessage("user", "user.created", event)
		if err != nil {
			return err
		}
		
		// Grava na outbox
		if err := tx.Create(outboxMessage).Error; err != nil {
			return err
		}
		
		return nil
	})
}
