package api

import (
	"net/http"
	"strconv"
	"user-service/internal/domain"
	"user-service/internal/outbox"
	"user-service/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UserHandler handler para operações de usuário
type UserHandler struct {
	userRepo      repo.UserRepository
	outboxService *outbox.OutboxService
}

// NewUserHandler cria um novo handler de usuário
func NewUserHandler(userRepo repo.UserRepository, outboxService *outbox.OutboxService) *UserHandler {
	return &UserHandler{
		userRepo:      userRepo,
		outboxService: outboxService,
	}
}

// CreateUser cria um novo usuário
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("erro ao validar request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	
	// Cria o usuário com evento na outbox
	user := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}
	
	if err := h.outboxService.CreateUserWithEvent(c.Request.Context(), user); err != nil {
		if err.Error() == "email já cadastrado" {
			log.Error().Str("email", req.Email).Msg("email já existe")
			c.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
			return
		}
		log.Error().Err(err).Msg("erro ao criar usuário")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}
	
	log.Info().
		Uint("user_id", user.ID).
		Str("email", user.Email).
		Msg("usuário criado com sucesso")
	
	c.JSON(http.StatusCreated, gin.H{
		"data": user.ToResponse(),
		"message": "Usuário criado com sucesso",
	})
}

// GetUser busca usuário por ID
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	user, err := h.userRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		log.Error().Err(err).Uint("user_id", uint(id)).Msg("erro ao buscar usuário")
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": user.ToResponse(),
	})
}
