package events

// OrderCreated evento de pedido criado
type OrderCreated struct {
	BaseEvent
	Order OrderData `json:"order"`
}

// OrderData dados do pedido
type OrderData struct {
	ID           uint           `json:"id"`
	UserID       uint           `json:"user_id"`
	Status       string         `json:"status"`
	TotalAmount  float64        `json:"total_amount"`
	Items        []OrderItem    `json:"items"`
}

// OrderItem item do pedido
type OrderItem struct {
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}

// OrderPaid evento de pedido pago
type OrderPaid struct {
	BaseEvent
	OrderID uint `json:"order_id"`
}

// OrderCanceled evento de pedido cancelado
type OrderCanceled struct {
	BaseEvent
	OrderID uint   `json:"order_id"`
	Reason  string `json:"reason,omitempty"`
}
