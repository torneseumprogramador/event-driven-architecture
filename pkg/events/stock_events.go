package events

// StockReserved evento de estoque reservado
type StockReserved struct {
	BaseEvent
	OrderID   uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// StockReleased evento de estoque liberado
type StockReleased struct {
	BaseEvent
	OrderID   uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
