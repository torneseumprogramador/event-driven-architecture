#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# URLs dos serviços
USER_SERVICE="http://localhost:8081"
PRODUCT_SERVICE="http://localhost:8082"
ORDER_SERVICE="http://localhost:8083"
QUERY_SERVICE="http://localhost:8084"

echo -e "${BLUE}🚀 TESTE RÁPIDO DA API - Event-Driven Architecture${NC}"
echo "=================================================="

# Função para fazer requisição HTTP
make_request() {
    local method=$1
    local url=$2
    local data=$3
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url")
    fi
    
    body=$(echo "$response" | head -n -1)
    status_code=$(echo "$response" | tail -n 1)
    
    echo "$body"
    return $status_code
}

# 1. Health Check
echo -e "${YELLOW}1. Verificando health dos serviços...${NC}"
services=("User Service" "Product Service" "Order Service" "Query Service")
urls=("$USER_SERVICE/healthz" "$PRODUCT_SERVICE/healthz" "$ORDER_SERVICE/healthz" "$QUERY_SERVICE/healthz")

for i in "${!services[@]}"; do
    response=$(make_request "GET" "${urls[$i]}")
    status_code=$?
    
    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ ${services[$i]} - OK${NC}"
    else
        echo -e "${RED}❌ ${services[$i]} - ERRO${NC}"
    fi
done

echo

# 2. Cadastrar Usuário
echo -e "${YELLOW}2. Cadastrando usuário...${NC}"
user_data='{"name":"João Silva","email":"joao@example.com"}'
response=$(make_request "POST" "$USER_SERVICE/users" "$user_data")
status_code=$?

if [ $status_code -eq 201 ]; then
    echo -e "${GREEN}✅ Usuário cadastrado com sucesso!${NC}"
    USER_ID=$(echo "$response" | jq -r '.data.id' 2>/dev/null)
    echo "User ID: $USER_ID"
else
    echo -e "${RED}❌ Erro ao cadastrar usuário${NC}"
    echo "$response"
    exit 1
fi

echo

# 3. Cadastrar Produto
echo -e "${YELLOW}3. Cadastrando produto...${NC}"
product_data='{"name":"iPhone 15","price":5999.99,"stock":10}'
response=$(make_request "POST" "$PRODUCT_SERVICE/products" "$product_data")
status_code=$?

if [ $status_code -eq 201 ]; then
    echo -e "${GREEN}✅ Produto cadastrado com sucesso!${NC}"
    PRODUCT_ID=$(echo "$response" | jq -r '.data.id' 2>/dev/null)
    echo "Product ID: $PRODUCT_ID"
else
    echo -e "${RED}❌ Erro ao cadastrar produto${NC}"
    echo "$response"
    exit 1
fi

echo

# 4. Criar Pedido
echo -e "${YELLOW}4. Criando pedido...${NC}"
order_data="{\"user_id\":$USER_ID,\"items\":[{\"product_id\":$PRODUCT_ID,\"quantity\":2,\"unit_price\":5999.99}]}"
response=$(make_request "POST" "$ORDER_SERVICE/orders" "$order_data")
status_code=$?

if [ $status_code -eq 201 ]; then
    echo -e "${GREEN}✅ Pedido criado com sucesso!${NC}"
    ORDER_ID=$(echo "$response" | jq -r '.data.id' 2>/dev/null)
    echo "Order ID: $ORDER_ID"
else
    echo -e "${RED}❌ Erro ao criar pedido${NC}"
    echo "$response"
    exit 1
fi

echo

# 5. Pagar Pedido
echo -e "${YELLOW}5. Pagando pedido...${NC}"
payment_data='{"payment_method":"credit_card"}'
response=$(make_request "POST" "$ORDER_SERVICE/orders/$ORDER_ID/pay" "$payment_data")
status_code=$?

if [ $status_code -eq 200 ]; then
    echo -e "${GREEN}✅ Pedido pago com sucesso!${NC}"
else
    echo -e "${RED}❌ Erro ao pagar pedido${NC}"
    echo "$response"
fi

echo

# 6. Listar Dados
echo -e "${YELLOW}6. Listando dados...${NC}"

echo -e "${BLUE}Usuários:${NC}"
response=$(make_request "GET" "$USER_SERVICE/users")
if [ $? -eq 200 ]; then
    echo "$response" | jq '.data | length' 2>/dev/null || echo "Dados encontrados"
else
    echo "Erro ao buscar usuários"
fi

echo -e "${BLUE}Produtos:${NC}"
response=$(make_request "GET" "$PRODUCT_SERVICE/products")
if [ $? -eq 200 ]; then
    echo "$response" | jq '.data | length' 2>/dev/null || echo "Dados encontrados"
else
    echo "Erro ao buscar produtos"
fi

echo -e "${BLUE}Pedidos:${NC}"
response=$(make_request "GET" "$ORDER_SERVICE/orders")
if [ $? -eq 200 ]; then
    echo "$response" | jq '.data | length' 2>/dev/null || echo "Dados encontrados"
else
    echo "Erro ao buscar pedidos"
fi

echo

# 7. Consultar Query Service
echo -e "${YELLOW}7. Consultando Query Service...${NC}"
response=$(make_request "GET" "$QUERY_SERVICE/orders")
if [ $? -eq 200 ]; then
    echo -e "${GREEN}✅ Query Service funcionando!${NC}"
    echo "$response" | jq '.data | length' 2>/dev/null || echo "Dados encontrados"
else
    echo -e "${RED}❌ Erro no Query Service${NC}"
    echo "$response"
fi

echo
echo -e "${GREEN}🎉 Teste concluído com sucesso!${NC}"
echo -e "${BLUE}IDs gerados:${NC}"
echo "User ID: $USER_ID"
echo "Product ID: $PRODUCT_ID"
echo "Order ID: $ORDER_ID"
