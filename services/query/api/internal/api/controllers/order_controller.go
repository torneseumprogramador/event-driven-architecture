package controllers

import (
	"net/http"
	"strconv"
	"query-api/internal/dto"
	"query-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// OrderController define o controller de pedidos
type OrderController struct {
	queryService services.QueryService
}

// NewOrderController cria uma nova instância de OrderController
func NewOrderController(queryService services.QueryService) *OrderController {
	return &OrderController{
		queryService: queryService,
	}
}

// GetOrders retorna todos os pedidos
func (c *OrderController) GetOrders(ctx *gin.Context) {
	orders, err := c.queryService.GetOrders(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("erro ao buscar pedidos")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno do servidor"})
		return
	}

	response := dto.ToOrdersResponse(orders)
	ctx.JSON(http.StatusOK, response)
}

// GetOrder retorna um pedido pelo ID
func (c *OrderController) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	order, err := c.queryService.GetOrderByID(ctx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Int("order_id", id).Msg("erro ao buscar pedido")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno do servidor"})
		return
	}

	if order == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "pedido não encontrado"})
		return
	}

	response := dto.ToOrderResponse(*order)
	ctx.JSON(http.StatusOK, response)
}
