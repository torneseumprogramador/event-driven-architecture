package entities

import (
	"time"
)

// ProductProjectionView representa a projeção de produto no MongoDB
type ProductProjectionView struct {
	ID        int       `bson:"_id"`
	Name      string    `bson:"name"`
	Price     float64   `bson:"price"`
	Stock     int       `bson:"stock"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
