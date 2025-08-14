package dto

import (
	"user-api/internal/domain/entities"
	"user-api/internal/dto/responses"
)

// ToUserResponse converte User para UserResponse
func ToUserResponse(user *entities.User) *responses.UserResponse {
	return &responses.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

// ToUserResponseList converte lista de User para UserResponse
func ToUserResponseList(users []entities.User) []responses.UserResponse {
	var responses []responses.UserResponse
	for _, user := range users {
		responses = append(responses, *ToUserResponse(&user))
	}
	return responses
}
