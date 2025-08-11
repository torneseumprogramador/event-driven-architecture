package api

import (
	"net/http"
	"strconv"
	"query-service/internal/projections"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// QueryHandler handler para consultas do read model
type QueryHandler struct {
	orderProjection   *projections.OrderProjection
	productProjection *projections.ProductProjection
	userProjection    *projections.UserProjection
}

// NewQueryHandler cria um novo handler de consulta
func NewQueryHandler(
	orderProjection *projections.OrderProjection,
	productProjection *projections.ProductProjection,
	userProjection *projections.UserProjection,
) *QueryHandler {
	return &QueryHandler{
		orderProjection:   orderProjection,
		productProjection: productProjection,
		userProjection:    userProjection,
	}
}

// GetOrder busca pedido por ID
func (h *QueryHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	order, err := h.orderProjection.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("order_id", uint(id)).Msg("erro ao buscar pedido")
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}

// GetOrders busca pedidos por usuário
func (h *QueryHandler) GetOrders(c *gin.Context) {
	userIDStr := c.Query("user_id")
	status := c.Query("status")
	
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id é obrigatório"})
		return
	}
	
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("user_id", userIDStr).Msg("user_id inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id inválido"})
		return
	}
	
	orders, err := h.orderProjection.GetByUser(c.Request.Context(), uint(userID), status)
	if err != nil {
		log.Error().Err(err).Uint("user_id", uint(userID)).Msg("erro ao buscar pedidos")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": orders,
	})
}

// GetProducts busca todos os produtos
func (h *QueryHandler) GetProducts(c *gin.Context) {
	products, err := h.productProjection.GetAll(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("erro ao buscar produtos")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}

// GetUsers busca todos os usuários
func (h *QueryHandler) GetUsers(c *gin.Context) {
	users, err := h.userProjection.GetAll(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("erro ao buscar usuários")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}
