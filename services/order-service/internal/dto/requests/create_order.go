package requests

// CreateOrderRequest request para criar pedido
type CreateOrderRequest struct {
	UserID uint                `json:"user_id" binding:"required"`
	Items  []CreateOrderItem   `json:"items" binding:"required,min=1"`
}

// CreateOrderItem item do pedido para criação
type CreateOrderItem struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"required,gt=0"`
}
