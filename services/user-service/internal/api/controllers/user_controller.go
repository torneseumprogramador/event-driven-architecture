package controllers

import (
	"net/http"
	"strconv"
	"user-service/internal/domain/entities"
	"user-service/internal/dto"
	"user-service/internal/dto/requests"
	"user-service/internal/services"
	"user-service/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UserController controller para operações de usuário
type UserController struct {
	userRepo   repo.UserRepository
	userService *services.UserService
}

// NewUserController cria um novo controller de usuário
func NewUserController(userRepo repo.UserRepository, userService *services.UserService) *UserController {
	return &UserController{
		userRepo:    userRepo,
		userService: userService,
	}
}

// CreateUser cria um novo usuário
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req requests.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar dados do usuário")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	user := &entities.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := c.userService.CreateUserWithEvent(ctx.Request.Context(), user); err != nil {
		if err.Error() == "email já cadastrado" {
			log.Error().Str("email", req.Email).Msg("email já existe")
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
			return
		}
		log.Error().Err(err).Msg("erro ao criar usuário")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	log.Info().Uint("user_id", user.ID).Str("email", user.Email).Msg("usuário criado com sucesso")
	ctx.JSON(http.StatusCreated, gin.H{
		"data":    dto.ToUserResponse(user),
		"message": "Usuário criado com sucesso",
	})
}

// GetUser busca usuário por ID
func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	user, err := c.userRepo.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("user_id", uint(id)).Msg("erro ao buscar usuário")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"data": dto.ToUserResponse(user),
	})
}

// ListUsers lista todos os usuários
func (c *UserController) ListUsers(ctx *gin.Context) {
	users, err := c.userRepo.GetAll(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("erro ao listar usuários")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	responses := dto.ToUserResponseList(users)
	ctx.JSON(http.StatusOK, gin.H{"data": responses, "total": len(responses)})
}

// UpdateUser atualiza um usuário
func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req requests.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar dados do usuário")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	user, err := c.userRepo.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("user_id", uint(id)).Msg("usuário não encontrado")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// Atualiza os campos
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := c.userRepo.Update(ctx.Request.Context(), user); err != nil {
		log.Error().Err(err).Uint("user_id", uint(id)).Msg("erro ao atualizar usuário")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	log.Info().Uint("user_id", uint(id)).Msg("usuário atualizado com sucesso")
	ctx.JSON(http.StatusOK, gin.H{
		"data":    dto.ToUserResponse(user),
		"message": "Usuário atualizado com sucesso",
	})
}

// DeleteUser remove um usuário
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	user, err := c.userRepo.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("user_id", uint(id)).Msg("usuário não encontrado")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	if err := c.userRepo.Delete(ctx.Request.Context(), uint(id)); err != nil {
		log.Error().Err(err).Uint("user_id", uint(id)).Msg("erro ao remover usuário")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	log.Info().Uint("user_id", uint(id)).Str("email", user.Email).Msg("usuário removido com sucesso")
	ctx.JSON(http.StatusOK, gin.H{"message": "Usuário removido com sucesso"})
}
