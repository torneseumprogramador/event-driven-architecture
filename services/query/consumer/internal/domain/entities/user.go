package entities

import (
	"time"
)

// UserProjectionView representa a projeção de usuário no MongoDB
type UserProjectionView struct {
	ID        int       `bson:"_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
