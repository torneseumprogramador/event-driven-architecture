package responses

import (
	"time"
)

// UserResponse representa a resposta de usuário
type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UsersResponse representa a resposta de lista de usuários
type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}
