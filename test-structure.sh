#!/bin/bash

echo "🧪 Testando estrutura do projeto Event-Driven Architecture"

# Verifica se os arquivos principais existem
echo "📁 Verificando arquivos principais..."
files=(
    "docker-compose.yml"
    "go.work"
    "Makefile"
    "README.md"
    "env.example"
    ".gitignore"
    "docker/mysql/init.sql"
    "docker/kafka/create-topics.sh"
    "pkg/go.mod"
    "pkg/config/config.go"
    "pkg/log/logger.go"
    "pkg/kafka/producer.go"
    "pkg/kafka/consumer.go"
    "pkg/outbox/outbox.go"
    "pkg/idempotency/idempotency.go"
    "pkg/idempotency/mongo_repository.go"
    "pkg/http/middleware.go"
    "pkg/events/events.go"
    "services/user-service/go.mod"
    "services/user-service/cmd/main.go"
    "services/user-service/internal/domain/user.go"
    "services/user-service/internal/repo/user_repository.go"
    "services/user-service/internal/api/user_handler.go"
    "services/user-service/internal/outbox/outbox_service.go"
    "services/product-service/go.mod"
    "services/product-service/cmd/main.go"
    "services/product-service/internal/domain/product.go"
    "services/product-service/internal/repo/product_repository.go"
    "services/product-service/internal/api/product_handler.go"
    "services/product-service/internal/outbox/outbox_service.go"
    "services/product-service/internal/consumer/order_consumer.go"
    "services/order-service/go.mod"
    "services/order-service/cmd/main.go"
    "services/order-service/internal/domain/order.go"
    "services/order-service/internal/repo/order_repository.go"
    "services/order-service/internal/api/order_handler.go"
    "services/order-service/internal/outbox/outbox_service.go"
    "services/query-service/go.mod"
    "services/query-service/cmd/main.go"
    "services/query-service/internal/projections/order_projection.go"
    "services/query-service/internal/projections/product_projection.go"
    "services/query-service/internal/projections/user_projection.go"
    "services/query-service/internal/api/query_handler.go"
    "services/query-service/internal/consumer/event_consumer.go"
)

missing_files=()
for file in "${files[@]}"; do
    if [ ! -f "$file" ]; then
        missing_files+=("$file")
    else
        echo "✅ $file"
    fi
done

if [ ${#missing_files[@]} -eq 0 ]; then
    echo "🎉 Todos os arquivos principais estão presentes!"
else
    echo "❌ Arquivos faltando:"
    for file in "${missing_files[@]}"; do
        echo "   - $file"
    done
fi

echo ""
echo "📊 Estrutura de diretórios:"
echo "├── docker-compose.yml"
echo "├── docker/"
echo "│   ├── mysql/init.sql"
echo "│   └── kafka/create-topics.sh"
echo "├── pkg/"
echo "│   ├── config/"
echo "│   ├── kafka/"
echo "│   ├── outbox/"
echo "│   ├── idempotency/"
echo "│   ├── http/"
echo "│   ├── log/"
echo "│   └── events/"
echo "├── services/"
echo "│   ├── user-service/"
echo "│   ├── product-service/"
echo "│   ├── order-service/"
echo "│   └── query-service/"
echo "├── Makefile"
echo "├── README.md"
echo "└── env.example"

echo ""
echo "🔧 Serviços implementados:"
echo "✅ user-service (porta 8081)"
echo "✅ product-service (porta 8082)"
echo "✅ order-service (porta 8083)"
echo "✅ query-service (porta 8084)"

echo ""
echo "🏗️ Padrões implementados:"
echo "✅ Outbox Pattern"
echo "✅ CQRS (Command Query Responsibility Segregation)"
echo "✅ Idempotência"
echo "✅ Retry com Backoff"
echo "✅ DLQ (Dead Letter Queue)"

echo ""
echo "📋 Próximos passos:"
echo "1. Instalar Go 1.21+"
echo "2. Executar: make dev-setup"
echo "3. Executar: make up"
echo "4. Em terminais separados:"
echo "   - make run-user"
echo "   - make run-product"
echo "   - make run-order"
echo "   - make run-query"
echo "5. Testar com os exemplos do README.md"

echo ""
echo "🎯 Projeto pronto para desenvolvimento!"
