package dto

import (
	"time"
	"user-service/internal/domain"
)

// CreateUserRequest request para criar usuário
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateUserRequest request para atualizar usuário
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserResponse response do usuário
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converte User para UserResponse
func ToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

// ToUserResponseList converte lista de User para UserResponse
func ToUserResponseList(users []domain.User) []UserResponse {
	var responses []UserResponse
	for _, user := range users {
		responses = append(responses, *ToUserResponse(&user))
	}
	return responses
}
