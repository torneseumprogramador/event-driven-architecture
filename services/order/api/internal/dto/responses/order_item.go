package responses

// OrderItemResponse response do item do pedido
type OrderItemResponse struct {
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
