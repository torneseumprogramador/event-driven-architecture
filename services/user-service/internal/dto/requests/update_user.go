package requests

// UpdateUserRequest request para atualizar usuário
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
