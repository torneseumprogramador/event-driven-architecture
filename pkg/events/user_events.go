package events

// UserCreated evento de usuário criado
type UserCreated struct {
	BaseEvent
	User UserData `json:"user"`
}

// UserData dados do usuário
type UserData struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserUpdated evento de usuário atualizado
type UserUpdated struct {
	BaseEvent
	User UserData `json:"user"`
}
