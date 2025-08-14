package controllers

import (
	"net/http"
	"strconv"
	"query-api/internal/dto"
	"query-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QueryController define o controller de consultas
type QueryController struct {
	queryService services.QueryService
}

// NewQueryController cria uma nova instância de QueryController
func NewQueryController(queryService services.QueryService) *QueryController {
	return &QueryController{
		queryService: queryService,
	}
}

// GetUsers retorna todos os usuários
func (c *QueryController) GetUsers(ctx *gin.Context) {
	users, err := c.queryService.GetUsers(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("erro ao buscar usuários")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno do servidor"})
		return
	}

	response := dto.ToUsersResponse(users)
	ctx.JSON(http.StatusOK, response)
}

// GetUserByID retorna um usuário pelo ID
func (c *QueryController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	user, err := c.queryService.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Int("user_id", id).Msg("erro ao buscar usuário")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno do servidor"})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	response := dto.ToUserResponse(*user)
	ctx.JSON(http.StatusOK, response)
}

// GetProducts retorna todos os produtos
func (c *QueryController) GetProducts(ctx *gin.Context) {
	products, err := c.queryService.GetProducts(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("erro ao buscar produtos")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno do servidor"})
		return
	}

	response := dto.ToProductsResponse(products)
	ctx.JSON(http.StatusOK, response)
}

// GetProductByID retorna um produto pelo ID
func (c *QueryController) GetProductByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	product, err := c.queryService.GetProductByID(ctx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Int("product_id", id).Msg("erro ao buscar produto")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno do servidor"})
		return
	}

	if product == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "produto não encontrado"})
		return
	}

	response := dto.ToProductResponse(*product)
	ctx.JSON(http.StatusOK, response)
}

// GetOrders retorna todos os pedidos
func (c *QueryController) GetOrders(ctx *gin.Context) {
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
func (c *QueryController) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	order, err := c.queryService.GetOrderByID(ctx.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("order_id", idStr).Msg("erro ao buscar pedido")
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
