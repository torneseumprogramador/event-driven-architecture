package requests

// UpdateOrderRequest request para atualizar pedido
type UpdateOrderRequest struct {
	Status *string `json:"status,omitempty"`
}
