package controllers

import (
	"net/http"
	"strconv"
	"query-api/internal/dto"
	"query-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UserController define o controller de usuários
type UserController struct {
	queryService services.QueryService
}

// NewUserController cria uma nova instância de UserController
func NewUserController(queryService services.QueryService) *UserController {
	return &UserController{
		queryService: queryService,
	}
}

// GetUsers retorna todos os usuários
func (c *UserController) GetUsers(ctx *gin.Context) {
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
func (c *UserController) GetUserByID(ctx *gin.Context) {
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
