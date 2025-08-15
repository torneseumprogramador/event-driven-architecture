package controllers

import (
	"net/http"
	"strconv"
	"query-api/internal/dto"
	"query-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ProductController define o controller de produtos
type ProductController struct {
	queryService services.QueryService
}

// NewProductController cria uma nova instância de ProductController
func NewProductController(queryService services.QueryService) *ProductController {
	return &ProductController{
		queryService: queryService,
	}
}

// GetProducts retorna todos os produtos
func (c *ProductController) GetProducts(ctx *gin.Context) {
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
func (c *ProductController) GetProductByID(ctx *gin.Context) {
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
