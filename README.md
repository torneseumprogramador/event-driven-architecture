# 🚀 Event-Driven Architecture com CQRS em Go

Um sistema de e-commerce completo construído com **Event-Driven Architecture (EDA)** e **CQRS (Command Query Responsibility Segregation)** em Go, utilizando Kafka para comunicação assíncrona, MySQL para o modelo de escrita e MongoDB para o modelo de leitura.

## 📋 Índice

- [🏗️ Arquitetura](#️-arquitetura)
- [🛠️ Tecnologias](#️-tecnologias)
- [📁 Estrutura do Projeto](#-estrutura-do-projeto)
- [🚀 Como Executar](#-como-executar)
- [📚 Padrões Implementados](#-padrões-implementados)
- [🔧 Configuração](#-configuração)
- [📖 Exemplos de Uso](#-exemplos-de-uso)
- [🧪 Testes](#-testes)
- [📊 Monitoramento](#-monitoramento)

## 🏗️ Arquitetura

### Diagrama da Arquitetura

![Diagrama da Arquitetura](diagrama-arquitetura2.png)

### Visão Geral

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User API      │    │ Product API     │    │  Order API      │    │  Query Service  │
│   (HTTP Only)   │    │   (HTTP Only)   │    │   (HTTP Only)   │    │   (Read Model)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │                       │
         └───────────────────────┼───────────────────────┼───────────────────────┘
                                 │                       │
                    ┌─────────────────────────────────────────────────────────────┐
                    │                        Kafka                                │
                    │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
                    │  │user.created │ │product.*    │ │order.*      │          │
                    │  │user.updated │ │             │ │             │          │
                    │  └─────────────┘ └─────────────┘ └─────────────┘          │
                    └─────────────────────────────────────────────────────────────┘
                                 │                       │
                    ┌─────────────────────────────────────────────────────────────┐
                    │                        MySQL                                 │
                    │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
                    │  │   users     │ │  products   │ │   orders    │          │
                    │  │   outbox    │ │   outbox    │ │order_products│         │
                    │  │processed_evt│ │processed_evt│ │   outbox    │          │
                    │  └─────────────┘ └─────────────┘ │processed_evt│          │
                    └─────────────────────────────────────────────────────────────┘
                                                       │
                    ┌─────────────────────────────────────────────────────────────┐
                    │                       MongoDB                                │
                    │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
                    │  │ views.users │ │views.products│ │ views.orders│          │
                    │  └─────────────┘ └─────────────┘ └─────────────┘          │
                    └─────────────────────────────────────────────────────────────┘
```

### Arquitetura Separada: APIs e Consumers

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              APIS (HTTP Only)                                   │
├─────────────────┬─────────────────┬─────────────────┬─────────────────────────┤
│   User API      │  Product API    │   Order API     │    Query Service        │
│   Port: 8081    │   Port: 8082    │   Port: 8083    │    Port: 8084          │
│   - HTTP        │   - HTTP        │   - HTTP        │    - HTTP               │
│   - Outbox      │   - Outbox      │   - Outbox      │    - MongoDB            │
│   - MySQL       │   - MySQL       │   - MySQL       │    - Views              │
└─────────────────┴─────────────────┴─────────────────┴─────────────────────────┘
                                 │
                    ┌─────────────────────────────────────────────────────────────┐
                    │                        Kafka                                │
                    └─────────────────────────────────────────────────────────────┘
                                 │
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            CONSUMERS (Kafka Only)                               │
├─────────────────┬─────────────────┬─────────────────┬─────────────────────────┤
│User Consumer    │Product Consumer │Order Consumer   │    Query Service        │
│- Kafka          │- Kafka          │- Kafka          │    - Kafka              │
│- Outbox         │- Outbox         │- Outbox         │    - MongoDB            │
│- Dispatcher     │- Dispatcher     │- Dispatcher     │    - Views              │
│- MySQL          │- MySQL          │- MySQL          │    - Projections        │
└─────────────────┴─────────────────┴─────────────────┴─────────────────────────┘
```

### Fluxo de Dados

#### 1. Criação de Pedido
```
1. Cliente → Order Service (POST /orders)
2. Order Service → MySQL (transação)
   - Insere pedido
   - Insere itens do pedido
   - Insere evento na outbox
3. Outbox Dispatcher → Kafka (order.created)
4. Product Service → Consome evento → Reserva estoque
5. Query Service → Consome evento → Atualiza view
```

#### 2. Pagamento
```
1. Cliente → Order Service (POST /orders/{id}/pay)
2. Order Service → MySQL (transação)
   - Atualiza status para PAID
   - Insere evento na outbox
3. Outbox Dispatcher → Kafka (order.paid)
4. Query Service → Consome evento → Atualiza view
```

## 🛠️ Tecnologias

### Backend
- **Go 1.21+** - Linguagem principal
- **Gin** - Framework HTTP
- **GORM** - ORM para MySQL
- **MongoDB Go Driver** - Driver para MongoDB
- **Kafka Go** - Cliente Kafka
- **Zerolog** - Logging estruturado
- **Viper** - Configuração
- **UUID** - Geração de IDs únicos

### Infraestrutura
- **Docker & Docker Compose** - Containerização
- **Apache Kafka** - Message Broker
- **Zookeeper** - Coordenação do Kafka
- **MySQL 8.0** - Banco de dados transacional
- **MongoDB** - Banco de dados de leitura
- **Kafka UI** - Interface web para Kafka
- **Mongo Express** - Interface web para MongoDB

### Padrões
- **Event-Driven Architecture (EDA)**
- **CQRS (Command Query Responsibility Segregation)**
- **Outbox Pattern**
- **Idempotency Pattern**
- **Retry with Exponential Backoff**
- **Dead Letter Queue (DLQ)**

## 📁 Estrutura do Projeto

```
event-driven-architecture/
├── docker/                          # Configurações Docker
│   ├── mysql/
│   │   └── init.sql                 # Script de inicialização MySQL
│   └── kafka/
│       └── create-topics.sh         # Script de criação de tópicos
├── pkg/                             # Pacotes compartilhados
│   ├── config/                      # Configuração
│   ├── kafka/                       # Cliente Kafka
│   ├── outbox/                      # Padrão Outbox
│   ├── idempotency/                 # Controle de idempotência
│   ├── log/                         # Logging
│   └── http/                        # Utilitários HTTP
├── services/                        # Microserviços
│   ├── user/                        # Domínio de usuários
│   │   ├── api/                     # API HTTP (Port: 8081)
│   │   │   ├── internal/
│   │   │   │   ├── domain/
│   │   │   │   │   └── entities/    # Entidades de domínio
│   │   │   │   │       └── user.go
│   │   │   │   ├── dto/
│   │   │   │   │   ├── requests/    # DTOs de entrada
│   │   │   │   │   │   ├── create_user.go
│   │   │   │   │   │   └── update_user.go
│   │   │   │   │   ├── responses/   # DTOs de saída
│   │   │   │   │   │   └── user.go
│   │   │   │   │   └── converter.go # Conversores DTO
│   │   │   │   ├── api/
│   │   │   │   │   ├── controllers/ # Controllers HTTP
│   │   │   │   │   │   └── user_controller.go
│   │   │   │   │   └── routes/      # Configuração de rotas
│   │   │   │   │       └── routes.go
│   │   │   │   ├── repo/            # Repositórios
│   │   │   │   │   └── user_repository.go
│   │   │   │   └── services/        # Serviços de domínio
│   │   │   │       └── user_service.go
│   │   │   └── cmd/
│   │   │       └── main.go
│   │   └── consumer/                # Consumer Kafka
│   │       ├── internal/
│   │       │   ├── domain/
│   │       │   │   └── entities/    # Entidades de domínio
│   │       │   │       └── user.go
│   │       │   ├── repo/            # Repositórios
│   │       │   │   └── user_repository.go
│   │       │   └── services/        # Serviços de domínio
│   │       │       └── user_service.go
│   │       └── cmd/
│   │           └── main.go
│   ├── product/                     # Domínio de produtos
│   │   ├── api/                     # API HTTP (Port: 8082)
│   │   │   ├── internal/
│   │   │   │   ├── domain/
│   │   │   │   │   └── entities/    # Entidades de domínio
│   │   │   │   │       └── product.go
│   │   │   │   ├── dto/
│   │   │   │   │   ├── requests/    # DTOs de entrada
│   │   │   │   │   │   ├── create_product.go
│   │   │   │   │   │   └── update_product.go
│   │   │   │   │   ├── responses/   # DTOs de saída
│   │   │   │   │   │   └── product.go
│   │   │   │   │   └── converter.go # Conversores DTO
│   │   │   │   ├── api/
│   │   │   │   │   ├── controllers/ # Controllers HTTP
│   │   │   │   │   │   └── product_controller.go
│   │   │   │   │   └── routes/      # Configuração de rotas
│   │   │   │   │       └── routes.go
│   │   │   │   ├── repo/            # Repositórios
│   │   │   │   │   └── product_repository.go
│   │   │   │   └── services/        # Serviços de domínio
│   │   │   │       └── product_service.go
│   │   │   └── cmd/
│   │   │       └── main.go
│   │   └── consumer/                # Consumer Kafka
│   │       ├── internal/
│   │       │   ├── domain/
│   │       │   │   └── entities/    # Entidades de domínio
│   │       │   │       └── product.go
│   │       │   ├── consumer/        # Consumidores de eventos
│   │       │   │   └── order_consumer.go
│   │       │   ├── repo/            # Repositórios
│   │       │   │   └── product_repository.go
│   │       │   └── services/        # Serviços de domínio
│   │       │       └── product_service.go
│   │       └── cmd/
│   │           └── main.go
│   ├── order/                       # Domínio de pedidos
│   │   ├── api/                     # API HTTP (Port: 8083)
│   │   │   ├── internal/
│   │   │   │   ├── domain/
│   │   │   │   │   └── entities/    # Entidades de domínio
│   │   │   │   │   ├── order.go
│   │   │   │   │   └── order_product.go
│   │   │   │   ├── dto/
│   │   │   │   │   ├── requests/    # DTOs de entrada
│   │   │   │   │   │   ├── create_order.go
│   │   │   │   │   │   ├── update_order.go
│   │   │   │   │   │   ├── pay_order.go
│   │   │   │   │   │   └── cancel_order.go
│   │   │   │   │   ├── responses/   # DTOs de saída
│   │   │   │   │   │   ├── order.go
│   │   │   │   │   │   └── order_item.go
│   │   │   │   │   └── converter.go # Conversores DTO
│   │   │   │   ├── api/
│   │   │   │   │   ├── controllers/ # Controllers HTTP
│   │   │   │   │   │   └── order_controller.go
│   │   │   │   │   └── routes/      # Configuração de rotas
│   │   │   │   │       └── routes.go
│   │   │   │   ├── repo/            # Repositórios
│   │   │   │   │   └── order_repository.go
│   │   │   │   └── services/        # Serviços de domínio
│   │   │   │       └── order_service.go
│   │   │   └── cmd/
│   │   │       └── main.go
│   │   └── consumer/                # Consumer Kafka
│   │       ├── internal/
│   │       │   ├── domain/
│   │       │   │   └── entities/    # Entidades de domínio
│   │       │   │   ├── order.go
│   │       │   │   └── order_product.go
│   │       │   ├── repo/            # Repositórios
│   │       │   │   └── order_repository.go
│   │       │   └── services/        # Serviços de domínio
│   │       │       └── order_service.go
│   │       └── cmd/
│   │           └── main.go
│   └── query-service/               # Serviço de consultas (Port: 8084)
│       ├── internal/
│       │   ├── domain/
│       │   │   └── entities/        # Entidades de domínio
│       │   │       └── views.go
│       │   ├── dto/
│       │   │   ├── requests/        # DTOs de entrada
│       │   │   └── responses/       # DTOs de saída
│       │   ├── api/
│       │   │   ├── controllers/     # Controllers HTTP
│       │   │   └── routes/          # Configuração de rotas
│       │   ├── repo/                # Repositórios MongoDB
│       │   ├── services/            # Serviços de domínio
│       │   └── consumer/            # Consumidores de eventos
│       └── cmd/
│           └── main.go
├── docker-compose.yml               # Orquestração Docker
├── go.work                          # Workspace Go
├── Makefile                         # Comandos de automação
├── .env.example                     # Exemplo de variáveis de ambiente
└── README.md                        # Esta documentação
```

## 🚀 Como Executar

### Pré-requisitos

- **Go 1.21+**
- **Docker & Docker Compose**
- **Git**

### 1. Clone o Repositório

```bash
git clone git@github.com:torneseumprogramador/event-driven-architecture.git
cd event-driven-architecture
```

### 2. Configure as Variáveis de Ambiente

```bash
cp .env.example .env
# Edite o arquivo .env conforme necessário
```

### 3. Inicie a Infraestrutura

```bash
make up
```

Este comando irá:
- Iniciar Zookeeper, Kafka, MySQL, MongoDB
- Criar tópicos Kafka automaticamente
- Configurar bancos de dados

### 4. Execute os Serviços

#### Opção 1: Executar Tudo (APIs + Consumers)
```bash
# Executa todas as APIs e Consumers em terminais separados
make run-all
```

#### Opção 2: Executar Apenas APIs
```bash
# Executa todas as APIs (HTTP)
make run-apis
```

#### Opção 3: Executar Apenas Consumers
```bash
# Executa todos os consumers (Kafka)
make run-consumers
```

#### Opção 4: Executar Serviços Individualmente
```bash
# Terminal 1 - User API
make run-user-api

# Terminal 2 - Product API  
make run-product-api

# Terminal 3 - Order API
make run-order-api

# Terminal 4 - Query Service
make run-query-service

# Terminal 5 - User Consumer
make run-user-consumer

# Terminal 6 - Product Consumer
make run-product-consumer

# Terminal 7 - Order Consumer
make run-order-consumer
```

### 5. Verifique os Serviços

- **User API**: http://localhost:8081
- **Product API**: http://localhost:8082
- **Order API**: http://localhost:8083
- **Query Service**: http://localhost:8084
- **Kafka UI**: http://localhost:8080
- **Mongo Express**: http://localhost:8085

### 6. Benefícios da Separação APIs/Consumers

#### 🚀 **Escalabilidade Independente**
- **APIs**: Podem escalar horizontalmente para alta demanda HTTP
- **Consumers**: Podem escalar independentemente para processamento de eventos

#### ⚡ **Performance Otimizada**
- **APIs**: Foco em HTTP, sem overhead do consumer
- **Consumers**: Foco em processamento de eventos Kafka

#### 🔧 **Deployments Independentes**
- **APIs**: Deploy rápido para correções HTTP
- **Consumers**: Deploy independente para lógica de eventos

#### 📊 **Monitoramento Específico**
- **APIs**: Métricas HTTP (requests, latência, erros)
- **Consumers**: Métricas Kafka (lag, throughput, processamento)

#### 🛡️ **Isolamento de Falhas**
- Falha no consumer não afeta a API
- Falha na API não afeta o processamento de eventos

## 📚 Padrões Implementados

### 1. Event-Driven Architecture (EDA)

**Princípio**: Componentes se comunicam através de eventos assíncronos.

**Implementação**:
- Eventos publicados no Kafka
- Consumidores processam eventos independentemente
- Desacoplamento entre serviços

### 2. CQRS (Command Query Responsibility Segregation)

**Princípio**: Separação entre operações de escrita (Commands) e leitura (Queries).

**Implementação**:
- **Write Model**: MySQL com transações ACID
- **Read Model**: MongoDB com views denormalizadas
- **Commands**: Criar, atualizar, deletar entidades
- **Queries**: Consultas otimizadas para leitura

### 3. Outbox Pattern

**Princípio**: Garantir atomicidade entre escrita no banco e publicação de eventos.

**Implementação**:
```sql
-- Tabela outbox em cada serviço
CREATE TABLE outbox (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    aggregate VARCHAR(50) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL,
    headers JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP NULL
);
```

### 4. Idempotency Pattern

**Princípio**: Garantir que operações podem ser executadas múltiplas vezes sem efeitos colaterais.

**Implementação**:
```sql
-- Tabela de eventos processados
CREATE TABLE processed_events (
    event_id VARCHAR(36) PRIMARY KEY,
    service_name VARCHAR(50) NOT NULL,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 5. Retry with Exponential Backoff

**Princípio**: Reexecutar operações falhadas com delay exponencial.

**Implementação**:
- Máximo 5 tentativas
- Delay: 1s, 2s, 4s, 8s, 16s
- Dead Letter Queue para falhas persistentes

## 🔧 Configuração

### Variáveis de Ambiente

```bash
# Configurações gerais
ENV=development

# MySQL (Write Model)
MYSQL_DSN=ecommerce:ecommerce@tcp(localhost:3306)/ecommerce?parseTime=true

# MongoDB (Read Model)
MONGO_URI=mongodb://admin:admin@localhost:27017/ecommerce?authSource=admin

# Kafka
KAFKA_BROKERS=localhost:9092

# Outbox
OUTBOX_POLL_INTERVAL=1s
```

### Tópicos Kafka

| Tópico | Descrição | Partições | Replicação |
|--------|-----------|-----------|------------|
| `user.created` | Usuário criado | 3 | 1 |
| `user.updated` | Usuário atualizado | 3 | 1 |
| `product.created` | Produto criado | 3 | 1 |
| `product.updated` | Produto atualizado | 3 | 1 |
| `stock.reserved` | Estoque reservado | 3 | 1 |
| `stock.released` | Estoque liberado | 3 | 1 |
| `order.created` | Pedido criado | 3 | 1 |
| `order.paid` | Pedido pago | 3 | 1 |
| `order.canceled` | Pedido cancelado | 3 | 1 |

Cada tópico possui sua respectiva **Dead Letter Queue** (`*.dlq`).

## 📖 Exemplos de Uso

### User Service

#### Criar Usuário
```bash
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com"
  }'
```

#### Listar Usuários
```bash
curl http://localhost:8081/users
```

#### Obter Usuário por ID
```bash
curl http://localhost:8081/users/1
```

#### Atualizar Usuário
```bash
curl -X PUT http://localhost:8081/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva Atualizado",
    "email": "joao.novo@example.com"
  }'
```

#### Remover Usuário
```bash
curl -X DELETE http://localhost:8081/users/1
```

### Product Service

#### Criar Produto
```bash
curl -X POST http://localhost:8082/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "price": 5999.99,
    "stock": 10
  }'
```

#### Listar Produtos
```bash
curl http://localhost:8082/products
```

#### Obter Produto por ID
```bash
curl http://localhost:8082/products/1
```

#### Atualizar Produto
```bash
curl -X PUT http://localhost:8082/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 5499.99,
    "stock": 15
  }'
```

#### Remover Produto
```bash
curl -X DELETE http://localhost:8082/products/1
```

### Order Service

#### Criar Pedido
```bash
curl -X POST http://localhost:8083/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 2,
        "unit_price": 5999.99
      }
    ]
  }'
```

#### Listar Pedidos
```bash
curl http://localhost:8083/orders
```

#### Obter Pedido por ID
```bash
curl http://localhost:8083/orders/1
```

#### Atualizar Pedido
```bash
curl -X PUT http://localhost:8083/orders/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "PROCESSING"
  }'
```

#### Pagar Pedido
```bash
curl -X POST http://localhost:8083/orders/1/pay \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "credit_card"
  }'
```

#### Cancelar Pedido
```bash
curl -X POST http://localhost:8083/orders/1/cancel \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Cliente solicitou cancelamento"
  }'
```

#### Remover Pedido
```bash
curl -X DELETE http://localhost:8083/orders/1
```

### Query Service

#### Consultar Pedidos
```bash
curl http://localhost:8084/orders
```

#### Consultar Produtos
```bash
curl http://localhost:8084/products
```

#### Consultar Usuários
```bash
curl http://localhost:8084/users
```

## 🧪 Testes

### Executar Testes
```bash
make test
```

### Executar Linter
```bash
make lint
```

### Verificar Cobertura
```bash
make coverage
```

## 📊 Monitoramento

### Health Checks

Todos os serviços expõem endpoint de health check:

```bash
curl http://localhost:8081/healthz  # User Service
curl http://localhost:8082/healthz  # Product Service
curl http://localhost:8083/healthz  # Order Service
curl http://localhost:8084/healthz  # Query Service
```

### Logs Estruturados

Os serviços utilizam **Zerolog** para logs estruturados:

```json
{
  "level": "info",
  "service": "user-service",
  "message": "usuário criado com sucesso",
  "user_id": 1,
  "email": "joao@example.com",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Métricas

- **Kafka**: Tópicos, partições, offsets
- **MySQL**: Conexões, queries, transações
- **MongoDB**: Operações, conexões, índices
- **Aplicação**: Requests, latência, erros

## 🚀 Comandos Úteis

### Docker Compose
```bash
make up          # Iniciar infraestrutura
make down        # Parar infraestrutura
make logs        # Ver logs
make restart     # Reiniciar serviços
```

### Desenvolvimento
```bash
# Executar tudo
make run-all           # APIs + Consumers

# Executar apenas APIs
make run-apis          # Todas as APIs

# Executar apenas consumers
make run-consumers     # Todos os consumers

# Executar individualmente
make run-user-api      # User API
make run-product-api   # Product API
make run-order-api     # Order API
make run-query-service # Query Service
make run-user-consumer # User Consumer
make run-product-consumer # Product Consumer
make run-order-consumer   # Order Consumer
```

### Manutenção
```bash
make clean       # Limpar arquivos temporários
make deps        # Instalar dependências
make build       # Compilar todos os serviços
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👥 Autores

- **Danilo** - *Desenvolvimento inicial* - [torneseumprogramador](https://github.com/torneseumprogramador)

## 🙏 Agradecimentos

- Comunidade Go
- Apache Kafka
- GORM
- Gin Framework
- Docker

---

**⭐ Se este projeto te ajudou, considere dar uma estrela!**
