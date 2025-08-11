# E-commerce Event-Driven Architecture com CQRS

Este projeto implementa um sistema de e-commerce completo usando arquitetura Event-Driven (EDA) com CQRS (Command Query Responsibility Segregation) em Go. O sistema utiliza Kafka para comunicação assíncrona, MySQL para o write model e MongoDB para o read model.

## 🏗️ Arquitetura

### Componentes Principais

- **user-service**: Gerencia usuários (porta 8081)
- **product-service**: Gerencia produtos e estoque (porta 8082)
- **order-service**: Gerencia pedidos (porta 8083)
- **query-service**: Consultas do read model (porta 8084)

### Infraestrutura

- **Kafka**: Comunicação assíncrona entre serviços
- **MySQL**: Write model (dados transacionais)
- **MongoDB**: Read model (projeções para consultas)
- **Kafka UI**: Interface web para monitoramento (porta 8080)
- **Mongo Express**: Interface web para MongoDB (porta 8081)

### Padrões Implementados

- **Outbox Pattern**: Garante entrega de eventos mesmo em caso de falha
- **CQRS**: Separação entre comandos (write) e consultas (read)
- **Idempotência**: Evita processamento duplicado de eventos
- **Retry com Backoff**: Recuperação automática de falhas
- **DLQ (Dead Letter Queue)**: Tópicos para mensagens com falha

## 🚀 Quick Start

### Pré-requisitos

- Docker e Docker Compose
- Go 1.21+
- Make (opcional, mas recomendado)

### 1. Configuração Inicial

```bash
# Clone o repositório
git clone <repository-url>
cd event-driven-architecture

# Configure o ambiente
make dev-setup
```

### 2. Subir Infraestrutura

```bash
# Sobe toda a infraestrutura (Kafka, MySQL, MongoDB)
make up
```

### 3. Executar Serviços

Em terminais separados:

```bash
# Terminal 1 - User Service
make run-user

# Terminal 2 - Product Service
make run-product

# Terminal 3 - Order Service
make run-order

# Terminal 4 - Query Service
make run-query
```

### 4. Verificar Status

```bash
# Verifica se todos os serviços estão rodando
make health

# Verifica status dos containers
make status
```

## 📋 Endpoints da API

### User Service (Porta 8081)

```bash
# Criar usuário
curl -X POST http://localhost:8081/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Danilo Silva",
    "email": "danilo@exemplo.com"
  }'

# Buscar usuário
curl http://localhost:8081/api/v1/users/1
```

### Product Service (Porta 8082)

```bash
# Criar produto
curl -X POST http://localhost:8082/api/v1/products \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Camiseta Básica",
    "price": 59.90,
    "stock": 100
  }'

# Atualizar produto
curl -X PATCH http://localhost:8082/api/v1/products/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "price": 49.90,
    "stock": 80
  }'
```

### Order Service (Porta 8083)

```bash
# Criar pedido
curl -X POST http://localhost:8083/api/v1/orders \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'

# Pagar pedido
curl -X POST http://localhost:8083/api/v1/orders/1/pay

# Cancelar pedido
curl -X POST http://localhost:8083/api/v1/orders/1/cancel
```

### Query Service (Porta 8084)

```bash
# Buscar pedido com dados denormalizados
curl http://localhost:8084/q/orders/1

# Buscar pedidos por usuário
curl "http://localhost:8084/q/orders?user_id=1&status=PAID"

# Listar produtos
curl http://localhost:8084/q/products

# Listar usuários
curl http://localhost:8084/q/users
```

## 🔄 Fluxo de Eventos

### 1. Criação de Pedido

1. **order-service** recebe requisição para criar pedido
2. Grava `orders` e `order_products` no MySQL (write model)
3. Insere evento `order.created` na tabela `outbox`
4. **Dispatcher** publica evento no tópico `order.created`
5. **product-service** consome evento e tenta reservar estoque
6. Se sucesso: publica `stock.reserved`
7. Se falha: publica `order.canceled`
8. **query-service** consome todos os eventos e atualiza projeções no MongoDB

### 2. Pagamento de Pedido

1. **order-service** recebe requisição de pagamento
2. Atualiza status para `PAID` no MySQL
3. Insere evento `order.paid` na outbox
4. **Dispatcher** publica evento no tópico `order.paid`
5. **query-service** consome evento e atualiza projeção

## 🛠️ Comandos Úteis

```bash
# Ver todos os comandos disponíveis
make help

# Logs específicos
make logs-kafka    # Logs do Kafka
make logs-mysql    # Logs do MySQL
make logs-mongo    # Logs do MongoDB

# Desenvolvimento
make lint          # Executa linter
make test          # Executa testes
make build-all     # Compila todos os serviços

# Limpeza
make down          # Para infraestrutura
make clean         # Limpa tudo
```

## 📊 Monitoramento

### Kafka UI
- URL: http://localhost:8080
- Monitora tópicos, consumidores e mensagens

### Mongo Express
- URL: http://localhost:8081
- Usuário: admin
- Senha: admin
- Visualiza projeções no MongoDB

### Health Checks
```bash
# Verifica saúde de todos os serviços
make health
```

## 🔧 Configuração

### Variáveis de Ambiente

Copie `env.example` para `.env` e ajuste conforme necessário:

```bash
# Configurações gerais
ENV=development

# MySQL (Write Model)
MYSQL_DSN=ecommerce:ecommerce@tcp(mysql:3306)/ecommerce?parseTime=true

# MongoDB (Read Model)
MONGO_URI=mongodb://admin:admin@mongo:27017/ecommerce?authSource=admin

# Kafka
KAFKA_BROKERS=kafka:9092

# Outbox
OUTBOX_POLL_INTERVAL=1s
```

### Portas dos Serviços

- **Kafka UI**: 8080
- **Mongo Express**: 8081
- **User Service**: 8081
- **Product Service**: 8082
- **Order Service**: 8083
- **Query Service**: 8084

## 🏛️ Estrutura do Projeto

```
event-driven-architecture/
├── docker-compose.yml          # Infraestrutura Docker
├── docker/
│   ├── mysql/init.sql         # Script de inicialização MySQL
│   └── kafka/create-topics.sh # Script para criar tópicos
├── pkg/                       # Pacotes compartilhados
│   ├── config/               # Configuração
│   ├── kafka/                # Wrappers Kafka
│   ├── outbox/               # Padrão Outbox
│   ├── idempotency/          # Controle de idempotência
│   ├── http/                 # Middlewares HTTP
│   ├── log/                  # Logging
│   └── events/               # Contratos de eventos
├── services/
│   ├── user-service/         # Serviço de usuários
│   ├── product-service/      # Serviço de produtos
│   ├── order-service/        # Serviço de pedidos
│   └── query-service/        # Serviço de consultas
├── Makefile                  # Comandos de automação
├── env.example              # Exemplo de variáveis
└── README.md                # Esta documentação
```

## 🧪 Testando o Sistema

### Cenário Completo

1. **Criar usuário**:
```bash
curl -X POST http://localhost:8081/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"João Silva","email":"joao@exemplo.com"}'
```

2. **Criar produto**:
```bash
curl -X POST http://localhost:8082/api/v1/products \
  -H 'Content-Type: application/json' \
  -d '{"name":"Notebook Dell","price":3500.00,"stock":10}'
```

3. **Criar pedido**:
```bash
curl -X POST http://localhost:8083/api/v1/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_id":1,"items":[{"product_id":1,"quantity":1}]}'
```

4. **Pagar pedido**:
```bash
curl -X POST http://localhost:8083/api/v1/orders/1/pay
```

5. **Consultar read model**:
```bash
curl http://localhost:8084/q/orders/1
```

### Verificando Eventos

1. Acesse **Kafka UI** (http://localhost:8080)
2. Verifique os tópicos criados
3. Monitore mensagens em tempo real
4. Verifique DLQs se houver falhas

## 🔍 Troubleshooting

### Problemas Comuns

1. **Serviços não iniciam**:
   - Verifique se a infraestrutura está rodando: `make status`
   - Verifique logs: `make logs`

2. **Eventos não fluem**:
   - Verifique se os tópicos foram criados: `make create-topics`
   - Verifique logs do Kafka: `make logs-kafka`

3. **Projeções não atualizam**:
   - Verifique conexão com MongoDB: `make logs-mongo`
   - Verifique logs do query-service: `make logs-query`

### Logs Estruturados

Todos os serviços usam logs estruturados com zerolog. Os logs incluem:
- Service name
- Timestamp
- Log level
- Contexto específico (event_id, user_id, etc.)

## 📚 Conceitos Implementados

### Outbox Pattern
- Eventos são gravados na mesma transação dos dados
- Dispatcher processa outbox periodicamente
- Garante entrega mesmo em caso de falha

### CQRS
- **Commands**: Modificam estado (MySQL)
- **Queries**: Leem dados (MongoDB)
- **Projeções**: Dados denormalizados para consultas rápidas

### Idempotência
- Controle de eventos processados
- Evita processamento duplicado
- Tabela `processed_events` por serviço

### Retry e DLQ
- Retry exponencial (1s, 2s, 4s, 8s, 16s)
- Máximo 5 tentativas
- Falhas vão para DLQ (`<topic>.dlq`)

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature
3. Implemente as mudanças
4. Execute testes: `make test`
5. Execute linter: `make lint`
6. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.
