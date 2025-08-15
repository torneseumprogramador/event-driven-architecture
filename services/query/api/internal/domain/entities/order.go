package entities

import (
	"time"
)

// OrderView representa a view de pedido no MongoDB
type OrderView struct {
	ID          int                `bson:"_id" json:"id"`
	UserID      int                `bson:"user_id" json:"user_id"`
	Status      string             `bson:"status" json:"status"`
	Total       float64            `bson:"total_amount" json:"total"`
	User        UserView           `bson:"user" json:"user"`
	Items       []OrderItemView    `bson:"items" json:"items"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	PaidAt      *time.Time         `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
	CanceledAt  *time.Time         `bson:"canceled_at,omitempty" json:"canceled_at,omitempty"`
}





// OrderItemView representa um item do pedido na view
type OrderItemView struct {
	ProductID  int     `bson:"product_id" json:"product_id"`
	Quantity   int     `bson:"quantity" json:"quantity"`
	UnitPrice  float64 `bson:"unit_price" json:"unit_price"`
	Product    ProductView `bson:"product" json:"product"`
}
