package requests

// PayOrderRequest request para pagar pedido
type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}
