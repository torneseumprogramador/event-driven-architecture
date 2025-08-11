package api

import (
	"net/http"
	"strconv"
	"product-service/internal/domain"
	"product-service/internal/outbox"
	"product-service/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ProductHandler handler para operações de produto
type ProductHandler struct {
	productRepo   repo.ProductRepository
	outboxService *outbox.OutboxService
}

// NewProductHandler cria um novo handler de produto
func NewProductHandler(productRepo repo.ProductRepository, outboxService *outbox.OutboxService) *ProductHandler {
	return &ProductHandler{
		productRepo:   productRepo,
		outboxService: outboxService,
	}
}

// CreateProduct cria um novo produto
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	
	// Cria o produto com evento na outbox
	product := &domain.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	
	if err := h.outboxService.CreateProductWithEvent(c.Request.Context(), product); err != nil {
		log.Error().Err(err).Msg("erro ao criar produto")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().
		Uint("product_id", product.ID).
		Str("name", product.Name).
		Msg("produto criado com sucesso")
	
	c.JSON(http.StatusCreated, gin.H{
		"data": product.ToResponse(),
		"message": "Produto criado com sucesso",
	})
}

// UpdateProduct atualiza um produto
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	
	// Prepara as atualizações
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	
	if err := h.outboxService.UpdateProductWithEvent(c.Request.Context(), uint(id), updates); err != nil {
		if err.Error() == "produto não encontrado" {
			log.Error().Uint("product_id", uint(id)).Msg("produto não encontrado")
			c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
			return
		}
		log.Error().Err(err).Msg("erro ao atualizar produto")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().
		Uint("product_id", uint(id)).
		Msg("produto atualizado com sucesso")
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Produto atualizado com sucesso",
	})
}
