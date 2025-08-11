# 🚀 Event-Driven Architecture - Sistema E-commerce com CQRS

## 📚 Sobre o Projeto

Este projeto foi desenvolvido como parte do **Desafio de Arquiteturas de Software** do curso [Arquiteturas de Software Modernas](https://www.torneseumprogramador.com.br/cursos/arquiteturas_software) ministrado pelo **Prof. Danilo Aparecido** na plataforma [Torne-se um Programador](https://www.torneseumprogramador.com.br/).

### 🎯 Objetivo

Implementar um sistema de e-commerce completo utilizando **Event-Driven Architecture (EDA)** com **CQRS (Command Query Responsibility Segregation)**, demonstrando arquiteturas modernas, padrões de resiliência e comunicação assíncrona entre microserviços.

## 🏗️ Arquitetura

O projeto segue os princípios da **Event-Driven Architecture** com **CQRS** e **microserviços**:

```
┌─────────────────────────────────────────────────────────────────┐
│                        API Gateway                              │
├─────────────────────────────────────────────────────────────────┤
│  User Service  │  Product Service  │  Order Service  │  Query  │
│  (Write Model) │   (Write Model)   │  (Write Model)  │ Service │
│                │                   │                 │(Read)   │
├─────────────────────────────────────────────────────────────────┤
│                    Kafka (Event Bus)                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │user.created │ │product.creat│ │order.created│ │stock.reserve│ │
│  │user.updated │ │product.updat│ │order.paid   │ │stock.release│ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│  MySQL (Write Model)  │  MongoDB (Read Model)  │  DLQ Topics  │
│  - Users             │  - views.orders         │  - *.dlq     │
│  - Products          │  - views.products       │              │
│  - Orders            │  - views.users          │              │
│  - Outbox            │                         │              │
└─────────────────────────────────────────────────────────────────┘
```

### 📁 Estrutura do Projeto

```
event-driven-architecture/
├── docker/                          # Configurações Docker
│   ├── mysql/init.sql              # Schema MySQL
│   └── kafka/create-topics.sh      # Criação de tópicos
├── pkg/                            # Pacotes compartilhados
│   ├── config/                     # Configuração da aplicação
│   ├── kafka/                      # Producer/Consumer Kafka
│   ├── outbox/                     # Padrão Outbox
│   ├── idempotency/                # Controle de idempotência
│   ├── http/                       # Middlewares HTTP
│   ├── log/                        # Logging estruturado
│   └── events/                     # Contratos de eventos
├── services/                       # Microserviços
│   ├── user-service/               # Gestão de usuários
│   │   ├── cmd/main.go
│   │   ├── internal/domain/
│   │   ├── internal/repo/
│   │   ├── internal/api/
│   │   └── internal/outbox/
│   ├── product-service/            # Gestão de produtos
│   │   ├── cmd/main.go
│   │   ├── internal/domain/
│   │   ├── internal/repo/
│   │   ├── internal/api/
│   │   ├── internal/outbox/
│   │   └── internal/consumer/
│   ├── order-service/              # Gestão de pedidos
│   │   ├── cmd/main.go
│   │   ├── internal/domain/
│   │   ├── internal/repo/
│   │   ├── internal/api/
│   │   └── internal/outbox/
│   └── query-service/              # Modelo de leitura
│       ├── cmd/main.go
│       ├── internal/projections/
│       ├── internal/api/
│       └── internal/consumer/
├── docker-compose.yml              # Orquestração Docker
├── go.work                         # Workspace Go
├── Makefile                        # Comandos de automação
├── env.example                     # Variáveis de ambiente
└── README.md                       # Esta documentação
```

## 🚀 Tecnologias Utilizadas

- **Go 1.24** - Linguagem de programação
- **Gin** - Framework HTTP
- **GORM** - ORM para MySQL
- **MongoDB Driver** - Driver para MongoDB
- **Kafka** - Plataforma de streaming de eventos
- **MySQL** - Banco de dados transacional (Write Model)
- **MongoDB** - Banco de dados de consulta (Read Model)
- **Docker Compose** - Orquestração de containers
- **Zerolog** - Logging estruturado
- **Viper** - Gerenciamento de configuração

## 📋 Pré-requisitos

- [Go 1.24+](https://golang.org/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [Git](https://git-scm.com/)
- [Make](https://www.gnu.org/software/make/) (opcional)

## ⚡ Como Executar

### Método Rápido (Recomendado)

```bash
# Clone o repositório
git clone https://github.com/torneseumprogramador/event-driven-architecture.git
cd event-driven-architecture

# Configure as variáveis de ambiente
cp env.example .env

# Execute o setup completo
make dev-setup

# Inicie a infraestrutura
make up

# Em terminais separados, execute os serviços:
make run-user      # Porta 8081
make run-product   # Porta 8082
make run-order     # Porta 8083
make run-query     # Porta 8084
```

### Método Manual

```bash
# 1. Iniciar infraestrutura
docker-compose up -d

# 2. Criar tópicos Kafka
./docker/kafka/create-topics.sh

# 3. Sincronizar dependências Go
go work sync

# 4. Executar serviços (em terminais separados)
cd services/user-service && go run cmd/main.go
cd services/product-service && go run cmd/main.go
cd services/order-service && go run cmd/main.go
cd services/query-service && go run cmd/main.go
```

### Comandos Disponíveis no Makefile

```bash
make up              # Inicia infraestrutura Docker
make down            # Para infraestrutura Docker
make logs            # Mostra logs de todos os serviços
make clean           # Limpa containers e volumes
make dev-setup       # Setup completo do ambiente
make create-topics   # Cria tópicos Kafka
make status          # Status dos containers
make health          # Health checks dos serviços
make run-user        # Executa user-service
make run-product     # Executa product-service
make run-order       # Executa order-service
make run-query       # Executa query-service
make lint            # Executa linter
make test            # Executa testes
```

## 🌐 Acessando os Serviços

Após executar o projeto, os serviços estarão disponíveis em:

| Serviço | URL | Descrição |
|---------|-----|-----------|
| **User Service** | http://localhost:8081 | Gestão de usuários |
| **Product Service** | http://localhost:8082 | Gestão de produtos |
| **Order Service** | http://localhost:8083 | Gestão de pedidos |
| **Query Service** | http://localhost:8084 | Consultas (Read Model) |
| **Kafka UI** | http://localhost:8080 | Interface Kafka |
| **Mongo Express** | http://localhost:8081 | Interface MongoDB |

### Health Checks

Todos os serviços possuem endpoint de health check:

```bash
curl http://localhost:8081/healthz  # User Service
curl http://localhost:8082/healthz  # Product Service
curl http://localhost:8083/healthz  # Order Service
curl http://localhost:8084/healthz  # Query Service
```

## 📖 Endpoints da API

### 👥 Usuários (User Service - Porta 8081)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/users` | Criar usuário |
| GET | `/users/{id}` | Buscar usuário por ID |

### 📦 Produtos (Product Service - Porta 8082)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/products` | Criar produto |
| PATCH | `/products/{id}` | Atualizar produto |

### 🛒 Pedidos (Order Service - Porta 8083)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/orders` | Criar pedido |
| POST | `/orders/{id}/pay` | Pagar pedido |
| POST | `/orders/{id}/cancel` | Cancelar pedido |

### 🔍 Consultas (Query Service - Porta 8084)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/q/orders/{id}` | Buscar pedido por ID |
| GET | `/q/orders` | Listar pedidos (com filtros) |
| GET | `/q/products` | Listar produtos |
| GET | `/q/users` | Listar usuários |

## 🏛️ Padrões Arquiteturais Implementados

### 📦 Event-Driven Architecture (EDA)

- **Comunicação assíncrona** via Kafka
- **Desacoplamento** entre serviços
- **Escalabilidade** horizontal
- **Resiliência** a falhas

### 🔄 CQRS (Command Query Responsibility Segregation)

- **Write Model**: MySQL com GORM
- **Read Model**: MongoDB com projeções
- **Separação** de responsabilidades
- **Otimização** para consultas

### 📮 Outbox Pattern

- **Atomicidade** entre transação e evento
- **Reliability** na entrega de eventos
- **Idempotência** garantida
- **Retry** automático

### 🛡️ Padrões de Resiliência

- **Retry com Exponential Backoff**
- **Dead Letter Queue (DLQ)**
- **Circuit Breaker** (implícito)
- **Health Checks**

## 🧪 Exemplos de Uso

### Criar Usuário

```bash
curl -X POST "http://localhost:8081/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@email.com"
  }'
```

### Criar Produto

```bash
curl -X POST "http://localhost:8082/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smartphone XYZ",
    "price": 1299.99,
    "stock": 50
  }'
```

### Criar Pedido

```bash
curl -X POST "http://localhost:8083/orders" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'
```

### Consultar Pedido

```bash
curl "http://localhost:8084/q/orders/1"
```

## 📊 Fluxos de Eventos

### 🛒 Fluxo: Criar Pedido

```
1. POST /orders → Order Service
2. Order Service → Salva em MySQL + Outbox
3. Outbox Dispatcher → Publish order.created
4. Product Service → Consome order.created
5. Product Service → Reserva estoque
6. Product Service → Publish stock.reserved
7. Query Service → Consome order.created
8. Query Service → Atualiza MongoDB
```

### 💳 Fluxo: Pagamento

```
1. POST /orders/{id}/pay → Order Service
2. Order Service → Atualiza status + Outbox
3. Outbox Dispatcher → Publish order.paid
4. Query Service → Consome order.paid
5. Query Service → Atualiza MongoDB
```

## 🔧 Configuração da Infraestrutura

### 🐳 Docker Compose

O projeto inclui toda a infraestrutura necessária:

- **Zookeeper**: Coordenação do Kafka
- **Kafka**: Plataforma de eventos
- **Kafka UI**: Interface web para Kafka
- **MySQL**: Banco transacional
- **MongoDB**: Banco de consultas
- **Mongo Express**: Interface web para MongoDB

### 📡 Kafka Topics

Tópicos criados automaticamente:

| Tópico | Partições | Replicação | Descrição |
|--------|-----------|------------|-----------|
| `user.created` | 3 | 1 | Usuário criado |
| `user.updated` | 3 | 1 | Usuário atualizado |
| `product.created` | 3 | 1 | Produto criado |
| `product.updated` | 3 | 1 | Produto atualizado |
| `order.created` | 3 | 1 | Pedido criado |
| `order.paid` | 3 | 1 | Pedido pago |
| `order.canceled` | 3 | 1 | Pedido cancelado |
| `stock.reserved` | 3 | 1 | Estoque reservado |
| `stock.released` | 3 | 1 | Estoque liberado |
| `*.dlq` | 3 | 1 | Dead Letter Queues |

## 🛡️ Tratamento de Erros

### 📋 Tipos de Erro

| Código | Tipo | Descrição |
|--------|------|-----------|
| 400 | `ValidationError` | Erro de validação |
| 404 | `NotFoundError` | Recurso não encontrado |
| 409 | `ConflictError` | Conflito (email duplicado) |
| 500 | `InternalServerError` | Erro interno |

### 📝 Formato da Resposta de Erro

```json
{
  "error": "Email já cadastrado",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 📝 Eventos e Contratos

### 📨 Estrutura Base do Evento

```json
{
  "event_id": "uuid-v4",
  "occurred_at": "2024-01-15T10:30:00Z",
  "data": {
    // Dados específicos do evento
  }
}
```

### 🔄 Eventos Implementados

- **UserCreated**: Usuário criado
- **UserUpdated**: Usuário atualizado
- **ProductCreated**: Produto criado
- **ProductUpdated**: Produto atualizado
- **OrderCreated**: Pedido criado
- **OrderPaid**: Pedido pago
- **OrderCanceled**: Pedido cancelado
- **StockReserved**: Estoque reservado
- **StockReleased**: Estoque liberado

## 🎓 Aprendizados do Curso

Este projeto demonstra os seguintes conceitos aprendidos no curso:

1. **Event-Driven Architecture (EDA)**
   - Comunicação assíncrona
   - Desacoplamento de serviços
   - Padrões de eventos

2. **CQRS (Command Query Responsibility Segregation)**
   - Separação de modelos
   - Otimização de consultas
   - Projeções de dados

3. **Microserviços**
   - Decomposição de domínio
   - Comunicação entre serviços
   - Independência de deploy

4. **Padrões de Resiliência**
   - Outbox Pattern
   - Idempotência
   - Retry com backoff
   - Dead Letter Queue

5. **Boas Práticas**
   - Clean Architecture
   - Domain-Driven Design
   - Structured Logging
   - Health Checks

## 👨‍🏫 Sobre o Professor

**Prof. Danilo Aparecido** é instrutor na plataforma [Torne-se um Programador](https://www.torneseumprogramador.com.br/), especializado em arquiteturas de software, microserviços e desenvolvimento de sistemas escaláveis.

## 📚 Curso Completo

Para aprender mais sobre arquiteturas de software e aprofundar seus conhecimentos, acesse o curso completo:

**[Arquiteturas de Software Modernas](https://www.torneseumprogramador.com.br/cursos/arquiteturas_software)**

## 🧪 Testando o Sistema

### Scripts de Teste

```bash
# Verificar estrutura do projeto
./test-structure.sh

# Executar exemplos de teste
./test-examples.sh
```

### Monitoramento

- **Kafka UI**: http://localhost:8080
- **Mongo Express**: http://localhost:8081
- **Logs estruturados** em todos os serviços

## 🤝 Contribuição

Este projeto foi desenvolvido como parte de um desafio educacional. Contribuições são bem-vindas através de issues e pull requests.

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

**Desenvolvido com ❤️ para o curso de Arquiteturas de Software do [Torne-se um Programador](https://www.torneseumprogramador.com.br/)**
