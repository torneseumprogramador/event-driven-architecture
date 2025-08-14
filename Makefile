.PHONY: help up down logs clean lint test run-apis run-consumers run-all create-topics test-api quick-test

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
	@cd services/user/api && $(GO) mod tidy && golangci-lint run
	@cd services/user/consumer && $(GO) mod tidy && golangci-lint run
	@cd services/product/api && $(GO) mod tidy && golangci-lint run
	@cd services/product/consumer && $(GO) mod tidy && golangci-lint run
	@cd services/order/api && $(GO) mod tidy && golangci-lint run
	@cd services/order/consumer && $(GO) mod tidy && golangci-lint run
	@cd services/query/api && $(GO) mod tidy && golangci-lint run

test: ## Executa testes em todos os serviços
	@echo "Executando testes..."
	@cd pkg && $(GO) test ./...
	@cd services/user/api && $(GO) test ./...
	@cd services/user/consumer && $(GO) test ./...
	@cd services/product/api && $(GO) test ./...
	@cd services/product/consumer && $(GO) test ./...
	@cd services/order/api && $(GO) test ./...
	@cd services/order/consumer && $(GO) test ./...
	@cd services/query/api && $(GO) test ./...

# =============================================================================
# EXECUÇÃO DE APIS
# =============================================================================

run-user-api: ## Executa a user-api
	@echo "Executando user-api..."
	@cd services/user/api && SERVICE_NAME=user-api PORT=8081 $(GO) run cmd/main.go

run-product-api: ## Executa a product-api
	@echo "Executando product-api..."
	@cd services/product/api && SERVICE_NAME=product-api PORT=8082 $(GO) run cmd/main.go

run-order-api: ## Executa a order-api
	@echo "Executando order-api..."
	@cd services/order/api && SERVICE_NAME=order-api PORT=8083 $(GO) run cmd/main.go

run-query-api: ## Executa a query-api
	@echo "Executando query-api..."
	@cd services/query/api && SERVICE_NAME=query-api PORT=8084 $(GO) run cmd/main.go

run-apis: run-user-api run-product-api run-order-api run-query-api ## Executa todas as APIs

# =============================================================================
# EXECUÇÃO DE CONSUMERS
# =============================================================================

run-user-consumer: ## Executa o user-consumer
	@echo "Executando user-consumer..."
	@cd services/user/consumer && SERVICE_NAME=user-consumer $(GO) run cmd/main.go

run-product-consumer: ## Executa o product-consumer
	@echo "Executando product-consumer..."
	@cd services/product/consumer && SERVICE_NAME=product-consumer $(GO) run cmd/main.go

run-order-consumer: ## Executa o order-consumer
	@echo "Executando order-consumer..."
	@cd services/order/consumer && SERVICE_NAME=order-consumer $(GO) run cmd/main.go

run-query-consumer: ## Executa o query-consumer
	@echo "Executando query-consumer..."
	@cd services/query/consumer && SERVICE_NAME=query-consumer $(GO) run cmd/main.go

run-consumers: run-user-consumer run-product-consumer run-order-consumer run-query-consumer ## Executa todos os consumers

# =============================================================================
# EXECUÇÃO COMPLETA
# =============================================================================

run-all: run-apis run-consumers ## Executa APIs e Consumers (em terminais separados)

# =============================================================================
# BUILD DOS SERVIÇOS
# =============================================================================

# Build das APIs
build-user-api: ## Compila a user-api
	@echo "Compilando user-api..."
	@cd services/user/api && $(GO) build -o bin/user-api cmd/main.go

build-product-api: ## Compila a product-api
	@echo "Compilando product-api..."
	@cd services/product/api && $(GO) build -o bin/product-api cmd/main.go

build-order-api: ## Compila a order-api
	@echo "Compilando order-api..."
	@cd services/order/api && $(GO) build -o bin/order-api cmd/main.go

build-query-api: ## Compila a query-api
	@echo "Compilando query-api..."
	@cd services/query/api && $(GO) build -o bin/query-api cmd/main.go

build-apis: build-user-api build-product-api build-order-api build-query-api ## Compila todas as APIs

# Build dos Consumers
build-user-consumer: ## Compila o user-consumer
	@echo "Compilando user-consumer..."
	@cd services/user/consumer && $(GO) build -o bin/user-consumer cmd/main.go

build-product-consumer: ## Compila o product-consumer
	@echo "Compilando product-consumer..."
	@cd services/product/consumer && $(GO) build -o bin/product-consumer cmd/main.go

build-order-consumer: ## Compila o order-consumer
	@echo "Compilando order-consumer..."
	@cd services/order/consumer && $(GO) build -o bin/order-consumer cmd/main.go

build-query-consumer: ## Compila o query-consumer
	@echo "Compilando query-consumer..."
	@cd services/query/consumer && $(GO) build -o bin/query-consumer cmd/main.go

build-consumers: build-user-consumer build-product-consumer build-order-consumer build-query-consumer ## Compila todos os consumers

build-all: build-apis build-consumers ## Compila todos os serviços

# =============================================================================
# DESENVOLVIMENTO
# =============================================================================

dev-setup: ## Configura ambiente de desenvolvimento
	@echo "Configurando ambiente de desenvolvimento..."
	@cp env.example .env
	@echo "Arquivo .env criado. Edite conforme necessário."

create-topics: ## Cria tópicos do Kafka manualmente
	@echo "Criando tópicos do Kafka..."
	@./docker/kafka/create-topics.sh

# =============================================================================
# STATUS E MONITORAMENTO
# =============================================================================

status: ## Mostra status dos serviços Docker
	$(DOCKER_COMPOSE) ps

health: ## Verifica healthcheck de todas as APIs
	@echo "Verificando healthcheck das APIs..."
	@curl -s http://localhost:8081/healthz || echo "user-api: ❌"
	@curl -s http://localhost:8082/healthz || echo "product-api: ❌"
	@curl -s http://localhost:8083/healthz || echo "order-api: ❌"
	@curl -s http://localhost:8084/healthz || echo "query-api: ❌"

# =============================================================================
# LOGS ESPECÍFICOS
# =============================================================================

# Logs da infraestrutura
logs-kafka: ## Mostra logs do Kafka
	$(DOCKER_COMPOSE) logs -f kafka

logs-mysql: ## Mostra logs do MySQL
	$(DOCKER_COMPOSE) logs -f mysql

logs-mongo: ## Mostra logs do MongoDB
	$(DOCKER_COMPOSE) logs -f mongo

# Logs das APIs (quando em Docker)
logs-user-api: ## Mostra logs da user-api
	$(DOCKER_COMPOSE) logs -f user-api

logs-product-api: ## Mostra logs da product-api
	$(DOCKER_COMPOSE) logs -f product-api

logs-order-api: ## Mostra logs da order-api
	$(DOCKER_COMPOSE) logs -f order-api

logs-query-api: ## Mostra logs do query-api
	@echo "Query API não roda no Docker, use: ps aux | grep query-api"

# Logs dos Consumers (quando em Docker)
logs-user-consumer: ## Mostra logs do user-consumer
	$(DOCKER_COMPOSE) logs -f user-consumer

logs-product-consumer: ## Mostra logs do product-consumer
	$(DOCKER_COMPOSE) logs -f product-consumer

logs-order-consumer: ## Mostra logs do order-consumer
	$(DOCKER_COMPOSE) logs -f order-consumer

logs-query-consumer: ## Mostra logs do query-consumer
	@echo "Query Consumer não roda no Docker, use: ps aux | grep query-consumer"

# =============================================================================
# TESTES
# =============================================================================

test-api: ## Executa teste interativo da API
	@echo "Executando teste interativo da API..."
	@./scripts/test-api.sh

quick-test: ## Executa teste rápido da API
	@echo "Executando teste rápido da API..."
	@./scripts/quick-test.sh
