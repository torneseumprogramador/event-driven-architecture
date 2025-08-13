package controllers

import (
	"net/http"
	"strconv"
	"order-service/internal/domain/entities"
	"order-service/internal/dto"
	"order-service/internal/dto/requests"
	"order-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// OrderController controller para operações de pedido
type OrderController struct {
	orderService *services.OrderService
}

// NewOrderController cria um novo controller de pedido
func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// CreateOrder cria um novo pedido
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req requests.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar dados do pedido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}
	
	// Cria os itens do pedido
	items := make([]entities.OrderProduct, len(req.Items))
	for i, item := range req.Items {
		items[i] = entities.OrderProduct{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}
	
	// Cria o pedido com evento na outbox
	order := &entities.Order{
		UserID:      req.UserID,
		Status:      "CREATED",
		TotalAmount: 0, // Será calculado
		Items:       items,
	}
	
	if err := c.orderService.CreateOrderWithEvent(ctx.Request.Context(), order); err != nil {
		log.Error().Err(err).Msg("erro ao criar pedido")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().Uint("order_id", order.ID).Uint("user_id", order.UserID).Msg("pedido criado com sucesso")
	
	ctx.JSON(http.StatusCreated, gin.H{
		"data":    dto.ToOrderResponse(order),
		"message": "Pedido criado com sucesso",
	})
}

// PayOrder paga um pedido
func (c *OrderController) PayOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	// Atualiza o status para PAID com evento na outbox
	if err := c.orderService.PayOrderWithEvent(ctx.Request.Context(), uint(id)); err != nil {
		if err.Error() == "pedido não encontrado" {
			log.Error().Uint("order_id", uint(id)).Msg("pedido não encontrado")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
			return
		}
		if err.Error()[:len("pedido não pode ser pago")] == "pedido não pode ser pago" {
			log.Error().Uint("order_id", uint(id)).Msg("pedido não pode ser pago")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("erro ao pagar pedido")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().Uint("order_id", uint(id)).Msg("pedido pago com sucesso")
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Pedido pago com sucesso",
	})
}

// GetOrder obtém um pedido por ID
func (c *OrderController) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	order, err := c.orderService.GetOrderByID(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("erro ao buscar pedido")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": dto.ToOrderResponse(order)})
}

// ListOrders lista todos os pedidos
func (c *OrderController) ListOrders(ctx *gin.Context) {
	orders, err := c.orderService.GetAllOrders(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("erro ao listar pedidos")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	responses := dto.ToOrderResponseList(orders)
	ctx.JSON(http.StatusOK, gin.H{"data": responses, "total": len(responses)})
}

// UpdateOrder atualiza um pedido
func (c *OrderController) UpdateOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req requests.UpdateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar dados do pedido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	order, err := c.orderService.GetOrderByID(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("pedido não encontrado")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
		return
	}

	// Atualiza os campos
	if req.Status != nil {
		order.Status = *req.Status
	}

	if err := c.orderService.UpdateOrder(ctx.Request.Context(), order); err != nil {
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("erro ao atualizar pedido")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	log.Info().Uint("order_id", uint(id)).Msg("pedido atualizado com sucesso")
	ctx.JSON(http.StatusOK, gin.H{
		"data":    dto.ToOrderResponse(order),
		"message": "Pedido atualizado com sucesso",
	})
}

// DeleteOrder remove um pedido
func (c *OrderController) DeleteOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	order, err := c.orderService.GetOrderByID(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("pedido não encontrado")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
		return
	}

	if err := c.orderService.DeleteOrder(ctx.Request.Context(), uint(id)); err != nil {
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("erro ao remover pedido")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	log.Info().Uint("order_id", uint(id)).Uint("user_id", order.UserID).Msg("pedido removido com sucesso")
	ctx.JSON(http.StatusOK, gin.H{"message": "Pedido removido com sucesso"})
}

// CancelOrder cancela um pedido
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	var req requests.CancelOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar dados do pedido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}
	
	// Atualiza o status para CANCELED com evento na outbox
	if err := c.orderService.CancelOrderWithEvent(ctx.Request.Context(), uint(id), req.Reason); err != nil {
		if err.Error() == "pedido não encontrado" {
			log.Error().Uint("order_id", uint(id)).Msg("pedido não encontrado")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
			return
		}
		if err.Error() == "pedido já está cancelado" || err.Error() == "pedido pago não pode ser cancelado" {
			log.Error().Uint("order_id", uint(id)).Msg("pedido não pode ser cancelado")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("erro ao cancelar pedido")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().Uint("order_id", uint(id)).Msg("pedido cancelado com sucesso")
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Pedido cancelado com sucesso",
	})
}
