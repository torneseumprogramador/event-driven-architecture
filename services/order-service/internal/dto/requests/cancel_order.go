package requests

// CancelOrderRequest request para cancelar pedido
type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}
