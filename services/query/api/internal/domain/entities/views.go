package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserView representa a view de usuário no MongoDB
type UserView struct {
	ID        int       `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

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

// OrderView representa a view de pedido no MongoDB
type OrderView struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	UserID      int                `bson:"user_id" json:"user_id"`
	Status      string             `bson:"status" json:"status"`
	Total       float64            `bson:"total" json:"total"`
	Items       []OrderItemView    `bson:"items" json:"items"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	PaidAt      *time.Time         `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
	CanceledAt  *time.Time         `bson:"canceled_at,omitempty" json:"canceled_at,omitempty"`
}

// OrderItemView representa um item do pedido na view
type OrderItemView struct {
	ProductID   int     `bson:"product_id" json:"product_id"`
	ProductName string  `bson:"product_name" json:"product_name"`
	Quantity    int     `bson:"quantity" json:"quantity"`
	Price       float64 `bson:"price" json:"price"`
	Total       float64 `bson:"total" json:"total"`
}
