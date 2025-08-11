package api

import (
	"net/http"
	"strconv"
	"order-service/internal/domain"
	"order-service/internal/outbox"
	"order-service/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// OrderHandler handler para operações de pedido
type OrderHandler struct {
	orderRepo     repo.OrderRepository
	outboxService *outbox.OutboxService
}

// NewOrderHandler cria um novo handler de pedido
func NewOrderHandler(orderRepo repo.OrderRepository, outboxService *outbox.OutboxService) *OrderHandler {
	return &OrderHandler{
		orderRepo:     orderRepo,
		outboxService: outboxService,
	}
}

// CreateOrder cria um novo pedido
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	
	// Cria os itens do pedido
	items := make([]domain.OrderProduct, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.OrderProduct{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			// UnitPrice será calculado pelo serviço de produto
		}
	}
	
	// Cria o pedido com evento na outbox
	order := &domain.Order{
		UserID:      req.UserID,
		Status:      "CREATED",
		TotalAmount: 0, // Será calculado
		Items:       items,
	}
	
	if err := h.outboxService.CreateOrderWithEvent(c.Request.Context(), order); err != nil {
		log.Error().Err(err).Msg("erro ao criar pedido")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().
		Uint("order_id", order.ID).
		Uint("user_id", order.UserID).
		Msg("pedido criado com sucesso")
	
	c.JSON(http.StatusCreated, gin.H{
		"data": order.ToResponse(),
		"message": "Pedido criado com sucesso",
	})
}

// PayOrder paga um pedido
func (h *OrderHandler) PayOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	// Atualiza o status para PAID com evento na outbox
	if err := h.outboxService.PayOrderWithEvent(c.Request.Context(), uint(id)); err != nil {
		if err.Error() == "pedido não encontrado" {
			log.Error().Uint("order_id", uint(id)).Msg("pedido não encontrado")
			c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
			return
		}
		if err.Error()[:len("pedido não pode ser pago")] == "pedido não pode ser pago" {
			log.Error().Uint("order_id", uint(id)).Msg("pedido não pode ser pago")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("erro ao pagar pedido")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().
		Uint("order_id", uint(id)).
		Msg("pedido pago com sucesso")
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Pedido pago com sucesso",
	})
}

// CancelOrder cancela um pedido
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	// Atualiza o status para CANCELED com evento na outbox
	if err := h.outboxService.CancelOrderWithEvent(c.Request.Context(), uint(id), "Cancelado pelo usuário"); err != nil {
		if err.Error() == "pedido não encontrado" {
			log.Error().Uint("order_id", uint(id)).Msg("pedido não encontrado")
			c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
			return
		}
		if err.Error() == "pedido já está cancelado" || err.Error() == "pedido pago não pode ser cancelado" {
			log.Error().Uint("order_id", uint(id)).Msg("pedido não pode ser cancelado")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("erro ao cancelar pedido")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().
		Uint("order_id", uint(id)).
		Msg("pedido cancelado com sucesso")
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Pedido cancelado com sucesso",
	})
}
