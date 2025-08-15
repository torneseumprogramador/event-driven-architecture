package services

import (
	"context"
	"query-api/internal/domain/entities"
	"query-api/internal/repo"
)

// QueryService define a interface do serviço de consultas
type QueryService interface {
	GetUsers(ctx context.Context) ([]entities.UserView, error)
	GetUserByID(ctx context.Context, id int) (*entities.UserView, error)
	GetProducts(ctx context.Context) ([]entities.ProductView, error)
	GetProductByID(ctx context.Context, id int) (*entities.ProductView, error)
	GetOrders(ctx context.Context) ([]entities.OrderView, error)
	GetOrderByID(ctx context.Context, id int) (*entities.OrderView, error)
}

// QueryServiceImpl implementa QueryService
type QueryServiceImpl struct {
	userRepo    repo.UserRepository
	productRepo repo.ProductRepository
	orderRepo   repo.OrderRepository
}

// NewQueryService cria uma nova instância de QueryService
func NewQueryService(
	userRepo repo.UserRepository,
	productRepo repo.ProductRepository,
	orderRepo repo.OrderRepository,
) QueryService {
	return &QueryServiceImpl{
		userRepo:    userRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

// GetUsers retorna todos os usuários
func (s *QueryServiceImpl) GetUsers(ctx context.Context) ([]entities.UserView, error) {
	return s.userRepo.FindAll(ctx)
}

// GetUserByID retorna um usuário pelo ID
func (s *QueryServiceImpl) GetUserByID(ctx context.Context, id int) (*entities.UserView, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GetProducts retorna todos os produtos
func (s *QueryServiceImpl) GetProducts(ctx context.Context) ([]entities.ProductView, error) {
	return s.productRepo.FindAll(ctx)
}

// GetProductByID retorna um produto pelo ID
func (s *QueryServiceImpl) GetProductByID(ctx context.Context, id int) (*entities.ProductView, error) {
	return s.productRepo.FindByID(ctx, id)
}

// GetOrders retorna todos os pedidos
func (s *QueryServiceImpl) GetOrders(ctx context.Context) ([]entities.OrderView, error) {
	return s.orderRepo.FindAll(ctx)
}

// GetOrderByID retorna um pedido pelo ID
func (s *QueryServiceImpl) GetOrderByID(ctx context.Context, id int) (*entities.OrderView, error) {
	return s.orderRepo.FindByID(ctx, id)
}
