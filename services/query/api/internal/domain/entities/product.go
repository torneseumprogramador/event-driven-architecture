package entities

import (
	"time"
)

// ProductView representa a view de produto no MongoDB
type ProductView struct {
	ID          int       `bson:"_id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Price       float64   `bson:"price" json:"price"`
	Stock       int       `bson:"stock" json:"stock"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
