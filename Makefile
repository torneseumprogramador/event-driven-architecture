.PHONY: help up down logs clean lint test run-user run-product run-order run-query create-topics

# Variáveis
DOCKER_COMPOSE = docker-compose
GO = go

# Comandos principais
help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

up: ## Sobe toda a infraestrutura (Kafka, MySQL, MongoDB)
	$(DOCKER_COMPOSE) up -d
	@echo "Aguardando serviços ficarem prontos..."
	@sleep 30
	@echo "Criando tópicos do Kafka..."
	@./docker/kafka/create-topics.sh

down: ## Para e remove toda a infraestrutura
	$(DOCKER_COMPOSE) down -v

logs: ## Mostra logs de todos os serviços
	$(DOCKER_COMPOSE) logs -f

clean: ## Limpa volumes e containers
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker system prune -f

lint: ## Executa linter em todos os serviços
	@echo "Executando linter..."
	@cd pkg && $(GO) mod tidy && golangci-lint run
	@cd services/user-service && $(GO) mod tidy && golangci-lint run
	@cd services/product-service && $(GO) mod tidy && golangci-lint run
	@cd services/order-service && $(GO) mod tidy && golangci-lint run
	@cd services/query-service && $(GO) mod tidy && golangci-lint run

test: ## Executa testes em todos os serviços
	@echo "Executando testes..."
	@cd pkg && $(GO) test ./...
	@cd services/user-service && $(GO) test ./...
	@cd services/product-service && $(GO) test ./...
	@cd services/order-service && $(GO) test ./...
	@cd services/query-service && $(GO) test ./...

# Comandos para executar serviços individualmente
run-user: ## Executa o user-service
	@echo "Executando user-service..."
	@cd services/user-service && $(GO) run cmd/main.go

run-product: ## Executa o product-service
	@echo "Executando product-service..."
	@cd services/product-service && $(GO) run cmd/main.go

run-order: ## Executa o order-service
	@echo "Executando order-service..."
	@cd services/order-service && $(GO) run cmd/main.go

run-query: ## Executa o query-service
	@echo "Executando query-service..."
	@cd services/query-service && $(GO) run cmd/main.go

# Comandos para build
build-user: ## Compila o user-service
	@echo "Compilando user-service..."
	@cd services/user-service && $(GO) build -o bin/user-service cmd/main.go

build-product: ## Compila o product-service
	@echo "Compilando product-service..."
	@cd services/product-service && $(GO) build -o bin/product-service cmd/main.go

build-order: ## Compila o order-service
	@echo "Compilando order-service..."
	@cd services/order-service && $(GO) build -o bin/order-service cmd/main.go

build-query: ## Compila o query-service
	@echo "Compilando query-service..."
	@cd services/query-service && $(GO) build -o bin/query-service cmd/main.go

build-all: build-user build-product build-order build-query ## Compila todos os serviços

# Comandos para desenvolvimento
dev-setup: ## Configura ambiente de desenvolvimento
	@echo "Configurando ambiente de desenvolvimento..."
	@cp env.example .env
	@echo "Arquivo .env criado. Edite conforme necessário."

create-topics: ## Cria tópicos do Kafka manualmente
	@echo "Criando tópicos do Kafka..."
	@./docker/kafka/create-topics.sh

# Comandos para verificar status
status: ## Mostra status dos serviços Docker
	$(DOCKER_COMPOSE) ps

health: ## Verifica healthcheck de todos os serviços
	@echo "Verificando healthcheck dos serviços..."
	@curl -s http://localhost:8081/healthz || echo "user-service: ❌"
	@curl -s http://localhost:8082/healthz || echo "product-service: ❌"
	@curl -s http://localhost:8083/healthz || echo "order-service: ❌"
	@curl -s http://localhost:8084/healthz || echo "query-service: ❌"

# Comandos para logs específicos
logs-kafka: ## Mostra logs do Kafka
	$(DOCKER_COMPOSE) logs -f kafka

logs-mysql: ## Mostra logs do MySQL
	$(DOCKER_COMPOSE) logs -f mysql

logs-mongo: ## Mostra logs do MongoDB
	$(DOCKER_COMPOSE) logs -f mongo

logs-user: ## Mostra logs do user-service
	$(DOCKER_COMPOSE) logs -f user-service

logs-product: ## Mostra logs do product-service
	$(DOCKER_COMPOSE) logs -f product-service

logs-order: ## Mostra logs do order-service
	$(DOCKER_COMPOSE) logs -f order-service

logs-query: ## Mostra logs do query-service
	$(DOCKER_COMPOSE) logs -f query-service
