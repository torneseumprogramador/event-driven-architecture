package controllers

import (
	"net/http"
	"strconv"
	"product-service/internal/domain"
	"product-service/internal/dto"
	"product-service/internal/services"
	"product-service/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ProductController controller para operações de produto
type ProductController struct {
	productRepo    repo.ProductRepository
	productService *services.ProductService
}

// NewProductController cria um novo controller de produto
func NewProductController(productRepo repo.ProductRepository, productService *services.ProductService) *ProductController {
	return &ProductController{
		productRepo:    productRepo,
		productService: productService,
	}
}

// CreateProduct cria um novo produto
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar dados do produto")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}
	
	// Cria o produto com evento na outbox
	product := &domain.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	
	if err := c.productService.CreateProductWithEvent(ctx.Request.Context(), product); err != nil {
		log.Error().Err(err).Msg("erro ao criar produto")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().Uint("product_id", product.ID).Str("name", product.Name).Msg("produto criado com sucesso")
	
	ctx.JSON(http.StatusCreated, gin.H{
		"data":    dto.ToProductResponse(product),
		"message": "Produto criado com sucesso",
	})
}

// GetProduct obtém um produto por ID
func (c *ProductController) GetProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	product, err := c.productRepo.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("product_id", uint(id)).Msg("erro ao buscar produto")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": dto.ToProductResponse(product)})
}

// ListProducts lista todos os produtos
func (c *ProductController) ListProducts(ctx *gin.Context) {
	products, err := c.productRepo.GetAll(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("erro ao listar produtos")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	responses := dto.ToProductResponseList(products)
	ctx.JSON(http.StatusOK, gin.H{"data": responses, "total": len(responses)})
}

// UpdateProduct atualiza um produto
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	var req dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar dados do produto")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
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
	
	if err := c.productService.UpdateProductWithEvent(ctx.Request.Context(), uint(id), updates); err != nil {
		if err.Error() == "produto não encontrado" {
			log.Error().Uint("product_id", uint(id)).Msg("produto não encontrado")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
			return
		}
		log.Error().Err(err).Msg("erro ao atualizar produto")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().Uint("product_id", uint(id)).Msg("produto atualizado com sucesso")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Produto atualizado com sucesso",
	})
}

// DeleteProduct remove um produto
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	product, err := c.productRepo.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("product_id", uint(id)).Msg("produto não encontrado")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}

	if err := c.productRepo.Delete(ctx.Request.Context(), uint(id)); err != nil {
		log.Error().Err(err).Uint("product_id", uint(id)).Msg("erro ao remover produto")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	log.Info().Uint("product_id", uint(id)).Str("name", product.Name).Msg("produto removido com sucesso")
	ctx.JSON(http.StatusOK, gin.H{"message": "Produto removido com sucesso"})
}
