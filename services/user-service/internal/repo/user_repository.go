package repo

import (
	"context"
	"user-service/internal/domain"

	"gorm.io/gorm"
)

// UserRepository interface para operações de usuário
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// GormUserRepository implementação usando GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository cria um novo repositório GORM
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Create cria um novo usuário
func (r *GormUserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID busca usuário por ID
func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail busca usuário por email
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
