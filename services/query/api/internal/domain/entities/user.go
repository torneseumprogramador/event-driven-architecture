package entities

import (
	"time"
)

// UserView representa a view de usuário no MongoDB
type UserView struct {
	ID        int       `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
