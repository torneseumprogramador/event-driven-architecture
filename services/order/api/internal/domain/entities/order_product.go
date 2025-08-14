package entities

// OrderProduct representa um item do pedido
type OrderProduct struct {
	OrderID   uint    `json:"order_id" gorm:"primaryKey"`
	ProductID uint    `json:"product_id" gorm:"primaryKey"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	UnitPrice float64 `json:"unit_price" gorm:"not null"`
}

// TableName especifica o nome da tabela
func (OrderProduct) TableName() string {
	return "order_products"
}
