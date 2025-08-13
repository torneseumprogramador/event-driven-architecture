#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# URLs dos serviços
USER_SERVICE="http://localhost:8081"
PRODUCT_SERVICE="http://localhost:8082"
ORDER_SERVICE="http://localhost:8083"
QUERY_SERVICE="http://localhost:8084"

# Variáveis para armazenar dados
USER_ID=""
PRODUCT_ID=""
ORDER_ID=""

# Função para imprimir cabeçalho
print_header() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                🚀 TESTE INTERATIVO DA API                    ║"
    echo "║              Event-Driven Architecture + CQRS                ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Função para imprimir menu
print_menu() {
    echo -e "${CYAN}"
    echo "┌─────────────────────────────────────────────────────────────┐"
    echo "│                        MENU PRINCIPAL                       │"
    echo "├─────────────────────────────────────────────────────────────┤"
    echo "│ 1. 📝 Cadastrar Usuário                                     │"
    echo "│ 2. 📦 Cadastrar Produto                                     │"
    echo "│ 3. 🛒 Criar Pedido                                          │"
    echo "│ 4. 💳 Pagar Pedido                                          │"
    echo "│ 5. ❌ Cancelar Pedido                                        │"
    echo "│ 6. 📋 Listar Usuários                                       │"
    echo "│ 7. 📋 Listar Produtos                                       │"
    echo "│ 8. 📋 Listar Pedidos                                        │"
    echo "│ 9. 🔍 Consultar Dados (Query Service)                       │"
    echo "│ 10. 🏥 Health Check dos Serviços                            │"
    echo "│ 11. 📊 Status dos Dados                                     │"
    echo "│ 0. 🚪 Sair                                                   │"
    echo "└─────────────────────────────────────────────────────────────┘"
    echo -e "${NC}"
}

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
    
    # Separar body e status code
    body=$(echo "$response" | head -n -1)
    status_code=$(echo "$response" | tail -n 1)
    
    echo "$body"
    return $status_code
}

# Função para cadastrar usuário
cadastrar_usuario() {
    echo -e "${YELLOW}📝 CADASTRAR USUÁRIO${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    read -p "Nome do usuário: " nome
    read -p "Email do usuário: " email
    
    if [ -z "$nome" ] || [ -z "$email" ]; then
        echo -e "${RED}❌ Nome e email são obrigatórios!${NC}"
        return
    fi
    
    data="{\"name\":\"$nome\",\"email\":\"$email\"}"
    
    echo -e "${BLUE}Enviando requisição...${NC}"
    response=$(make_request "POST" "$USER_SERVICE/users" "$data")
    status_code=$?
    
    if [ $status_code -eq 201 ]; then
        echo -e "${GREEN}✅ Usuário cadastrado com sucesso!${NC}"
        USER_ID=$(echo "$response" | jq -r '.data.id' 2>/dev/null)
        echo -e "${GREEN}ID do usuário: $USER_ID${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Erro ao cadastrar usuário${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função para cadastrar produto
cadastrar_produto() {
    echo -e "${YELLOW}📦 CADASTRAR PRODUTO${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    read -p "Nome do produto: " nome
    read -p "Preço do produto: " preco
    read -p "Quantidade em estoque: " estoque
    
    if [ -z "$nome" ] || [ -z "$preco" ] || [ -z "$estoque" ]; then
        echo -e "${RED}❌ Nome, preço e estoque são obrigatórios!${NC}"
        return
    fi
    
    data="{\"name\":\"$nome\",\"price\":$preco,\"stock\":$estoque}"
    
    echo -e "${BLUE}Enviando requisição...${NC}"
    response=$(make_request "POST" "$PRODUCT_SERVICE/products" "$data")
    status_code=$?
    
    if [ $status_code -eq 201 ]; then
        echo -e "${GREEN}✅ Produto cadastrado com sucesso!${NC}"
        PRODUCT_ID=$(echo "$response" | jq -r '.data.id' 2>/dev/null)
        echo -e "${GREEN}ID do produto: $PRODUCT_ID${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Erro ao cadastrar produto${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função para criar pedido
criar_pedido() {
    echo -e "${YELLOW}🛒 CRIAR PEDIDO${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    if [ -z "$USER_ID" ]; then
        read -p "ID do usuário (ou pressione ENTER para usar o último cadastrado): " user_id_input
        if [ -n "$user_id_input" ]; then
            USER_ID=$user_id_input
        else
            echo -e "${RED}❌ Nenhum usuário cadastrado! Cadastre um usuário primeiro.${NC}"
            return
        fi
    fi
    
    if [ -z "$PRODUCT_ID" ]; then
        read -p "ID do produto (ou pressione ENTER para usar o último cadastrado): " product_id_input
        if [ -n "$product_id_input" ]; then
            PRODUCT_ID=$product_id_input
        else
            echo -e "${RED}❌ Nenhum produto cadastrado! Cadastre um produto primeiro.${NC}"
            return
        fi
    fi
    
    read -p "Quantidade do produto: " quantidade
    read -p "Preço unitário: " preco_unitario
    
    if [ -z "$quantidade" ] || [ -z "$preco_unitario" ]; then
        echo -e "${RED}❌ Quantidade e preço unitário são obrigatórios!${NC}"
        return
    fi
    
    data="{\"user_id\":$USER_ID,\"items\":[{\"product_id\":$PRODUCT_ID,\"quantity\":$quantidade,\"unit_price\":$preco_unitario}]}"
    
    echo -e "${BLUE}Enviando requisição...${NC}"
    response=$(make_request "POST" "$ORDER_SERVICE/orders" "$data")
    status_code=$?
    
    if [ $status_code -eq 201 ]; then
        echo -e "${GREEN}✅ Pedido criado com sucesso!${NC}"
        ORDER_ID=$(echo "$response" | jq -r '.data.id' 2>/dev/null)
        echo -e "${GREEN}ID do pedido: $ORDER_ID${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Erro ao criar pedido${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função para pagar pedido
pagar_pedido() {
    echo -e "${YELLOW}💳 PAGAR PEDIDO${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    if [ -z "$ORDER_ID" ]; then
        read -p "ID do pedido: " order_id_input
        if [ -n "$order_id_input" ]; then
            ORDER_ID=$order_id_input
        else
            echo -e "${RED}❌ Nenhum pedido disponível! Crie um pedido primeiro.${NC}"
            return
        fi
    fi
    
    read -p "Método de pagamento (credit_card, debit_card, pix): " metodo_pagamento
    
    if [ -z "$metodo_pagamento" ]; then
        metodo_pagamento="credit_card"
    fi
    
    data="{\"payment_method\":\"$metodo_pagamento\"}"
    
    echo -e "${BLUE}Enviando requisição...${NC}"
    response=$(make_request "POST" "$ORDER_SERVICE/orders/$ORDER_ID/pay" "$data")
    status_code=$?
    
    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ Pedido pago com sucesso!${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Erro ao pagar pedido${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função para cancelar pedido
cancelar_pedido() {
    echo -e "${YELLOW}❌ CANCELAR PEDIDO${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    if [ -z "$ORDER_ID" ]; then
        read -p "ID do pedido: " order_id_input
        if [ -n "$order_id_input" ]; then
            ORDER_ID=$order_id_input
        else
            echo -e "${RED}❌ Nenhum pedido disponível! Crie um pedido primeiro.${NC}"
            return
        fi
    fi
    
    read -p "Motivo do cancelamento: " motivo
    
    if [ -z "$motivo" ]; then
        motivo="Cancelamento solicitado pelo usuário"
    fi
    
    data="{\"reason\":\"$motivo\"}"
    
    echo -e "${BLUE}Enviando requisição...${NC}"
    response=$(make_request "POST" "$ORDER_SERVICE/orders/$ORDER_ID/cancel" "$data")
    status_code=$?
    
    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ Pedido cancelado com sucesso!${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Erro ao cancelar pedido${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função para listar usuários
listar_usuarios() {
    echo -e "${YELLOW}📋 LISTAR USUÁRIOS${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    echo -e "${BLUE}Buscando usuários...${NC}"
    response=$(make_request "GET" "$USER_SERVICE/users")
    status_code=$?
    
    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ Usuários encontrados:${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Erro ao buscar usuários${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função para listar produtos
listar_produtos() {
    echo -e "${YELLOW}📋 LISTAR PRODUTOS${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    echo -e "${BLUE}Buscando produtos...${NC}"
    response=$(make_request "GET" "$PRODUCT_SERVICE/products")
    status_code=$?
    
    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ Produtos encontrados:${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Erro ao buscar produtos${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função para listar pedidos
listar_pedidos() {
    echo -e "${YELLOW}📋 LISTAR PEDIDOS${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    echo -e "${BLUE}Buscando pedidos...${NC}"
    response=$(make_request "GET" "$ORDER_SERVICE/orders")
    status_code=$?
    
    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ Pedidos encontrados:${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Erro ao buscar pedidos${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função para consultar dados no Query Service
consultar_dados() {
    echo -e "${YELLOW}🔍 CONSULTAR DADOS (QUERY SERVICE)${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    echo "1. Consultar pedidos"
    echo "2. Consultar produtos"
    echo "3. Consultar usuários"
    echo "4. Voltar ao menu principal"
    
    read -p "Escolha uma opção: " opcao_consulta
    
    case $opcao_consulta in
        1)
            echo -e "${BLUE}Buscando pedidos no Query Service...${NC}"
            response=$(make_request "GET" "$QUERY_SERVICE/orders")
            status_code=$?
            ;;
        2)
            echo -e "${BLUE}Buscando produtos no Query Service...${NC}"
            response=$(make_request "GET" "$QUERY_SERVICE/products")
            status_code=$?
            ;;
        3)
            echo -e "${BLUE}Buscando usuários no Query Service...${NC}"
            response=$(make_request "GET" "$QUERY_SERVICE/users")
            status_code=$?
            ;;
        4)
            return
            ;;
        *)
            echo -e "${RED}❌ Opção inválida!${NC}"
            return
            ;;
    esac
    
    if [ $status_code -eq 200 ]; then
        echo -e "${GREEN}✅ Dados encontrados:${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        echo -e "${RED}❌ Erro ao buscar dados${NC}"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função para health check
health_check() {
    echo -e "${YELLOW}🏥 HEALTH CHECK DOS SERVIÇOS${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    services=("User Service" "Product Service" "Order Service" "Query Service")
    urls=("$USER_SERVICE/healthz" "$PRODUCT_SERVICE/healthz" "$ORDER_SERVICE/healthz" "$QUERY_SERVICE/healthz")
    
    for i in "${!services[@]}"; do
        echo -e "${BLUE}Verificando ${services[$i]}...${NC}"
        response=$(make_request "GET" "${urls[$i]}")
        status_code=$?
        
        if [ $status_code -eq 200 ]; then
            echo -e "${GREEN}✅ ${services[$i]} está funcionando${NC}"
        else
            echo -e "${RED}❌ ${services[$i]} não está respondendo${NC}"
        fi
        echo
    done
    
    read -p "Pressione ENTER para continuar..."
}

# Função para mostrar status dos dados
status_dados() {
    echo -e "${YELLOW}📊 STATUS DOS DADOS${NC}"
    echo "─────────────────────────────────────────────────────────────"
    
    echo -e "${CYAN}IDs armazenados nesta sessão:${NC}"
    echo "User ID: ${USER_ID:-Nenhum}"
    echo "Product ID: ${PRODUCT_ID:-Nenhum}"
    echo "Order ID: ${ORDER_ID:-Nenhum}"
    echo
    
    echo -e "${CYAN}Contadores de dados:${NC}"
    
    # Contar usuários
    response=$(make_request "GET" "$USER_SERVICE/users")
    if [ $? -eq 200 ]; then
        user_count=$(echo "$response" | jq '.total // 0' 2>/dev/null || echo "0")
        echo "Usuários cadastrados: $user_count"
    else
        echo "Usuários cadastrados: Erro ao consultar"
    fi
    
    # Contar produtos
    response=$(make_request "GET" "$PRODUCT_SERVICE/products")
    if [ $? -eq 200 ]; then
        product_count=$(echo "$response" | jq '.total // 0' 2>/dev/null || echo "0")
        echo "Produtos cadastrados: $product_count"
    else
        echo "Produtos cadastrados: Erro ao consultar"
    fi
    
    # Contar pedidos
    response=$(make_request "GET" "$ORDER_SERVICE/orders")
    if [ $? -eq 200 ]; then
        order_count=$(echo "$response" | jq '.total // 0' 2>/dev/null || echo "0")
        echo "Pedidos criados: $order_count"
    else
        echo "Pedidos criados: Erro ao consultar"
    fi
    
    echo
    read -p "Pressione ENTER para continuar..."
}

# Função principal
main() {
    # Verificar se jq está instalado
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}❌ jq não está instalado. Instale jq para melhor formatação JSON.${NC}"
        echo "macOS: brew install jq"
        echo "Ubuntu/Debian: sudo apt-get install jq"
        echo "CentOS/RHEL: sudo yum install jq"
        echo
    fi
    
    # Verificar se curl está instalado
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}❌ curl não está instalado. Instale curl para fazer requisições HTTP.${NC}"
        exit 1
    fi
    
    while true; do
        clear
        print_header
        print_menu
        
        read -p "Escolha uma opção: " opcao
        
        case $opcao in
            1) cadastrar_usuario ;;
            2) cadastrar_produto ;;
            3) criar_pedido ;;
            4) pagar_pedido ;;
            5) cancelar_pedido ;;
            6) listar_usuarios ;;
            7) listar_produtos ;;
            8) listar_pedidos ;;
            9) consultar_dados ;;
            10) health_check ;;
            11) status_dados ;;
            0) 
                echo -e "${GREEN}👋 Obrigado por usar o script de teste!${NC}"
                exit 0
                ;;
            *)
                echo -e "${RED}❌ Opção inválida!${NC}"
                read -p "Pressione ENTER para continuar..."
                ;;
        esac
    done
}

# Executar função principal
main
