package entities

import (
	"time"
)

// OrderView representa a projeção de pedido no MongoDB
type OrderView struct {
	ID           int       `bson:"_id"`
	UserID       int       `bson:"user_id"`
	Status       string    `bson:"status"`
	TotalAmount  float64   `bson:"total_amount"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
	User         UserView  `bson:"user"`
	Items        []OrderItemView `bson:"items"`
}

// UserView representa a projeção de usuário
type UserView struct {
	ID    int    `bson:"id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

// OrderItemView representa a projeção de item do pedido
type OrderItemView struct {
	ProductID  int     `bson:"product_id"`
	Quantity   int     `bson:"quantity"`
	UnitPrice  float64 `bson:"unit_price"`
	Product    ProductView `bson:"product"`
}

// ProductView representa a projeção de produto
type ProductView struct {
	ID    int     `bson:"id"`
	Name  string  `bson:"name"`
	Price float64 `bson:"price"`
	Stock int     `bson:"stock"`
}
