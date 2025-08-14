package repo

import (
	"context"
	"user-consumer/internal/domain/entities"

	"gorm.io/gorm"
)

// UserRepository interface para operações de usuário
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uint) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetAll(ctx context.Context) ([]entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uint) error
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
func (r *GormUserRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID busca usuário por ID
func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail busca usuário por email
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll busca todos os usuários
func (r *GormUserRepository) GetAll(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Update atualiza um usuário
func (r *GormUserRepository) Update(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete remove um usuário
func (r *GormUserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.User{}, id).Error
}
