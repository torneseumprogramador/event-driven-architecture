# 📋 Scripts de Teste da API

Este diretório contém scripts para testar a API do sistema Event-Driven Architecture.

## 🚀 Scripts Disponíveis

### 1. `test-api.sh` - Teste Interativo Completo

Script interativo com menu para testar todas as funcionalidades da API.

**Características:**
- ✅ Interface colorida e amigável
- ✅ Menu interativo com opções numeradas
- ✅ Validação de entrada
- ✅ Armazenamento de IDs entre operações
- ✅ Formatação JSON com `jq` (opcional)
- ✅ Tratamento de erros
- ✅ Health checks dos serviços

**Como usar:**
```bash
# Executar via make
make test-api

# Ou executar diretamente
./scripts/test-api.sh
```

**Funcionalidades:**
1. 📝 Cadastrar Usuário
2. 📦 Cadastrar Produto
3. 🛒 Criar Pedido
4. 💳 Pagar Pedido
5. ❌ Cancelar Pedido
6. 📋 Listar Usuários
7. 📋 Listar Produtos
8. 📋 Listar Pedidos
9. 🔍 Consultar Dados (Query Service)
10. 🏥 Health Check dos Serviços
11. 📊 Status dos Dados

### 2. `quick-test.sh` - Teste Rápido Automatizado

Script para teste rápido e automatizado do fluxo completo.

**Características:**
- ✅ Execução automática sem interação
- ✅ Testa o fluxo completo: usuário → produto → pedido → pagamento
- ✅ Verificação de health dos serviços
- ✅ Consulta ao Query Service
- ✅ Relatório final com IDs gerados

**Como usar:**
```bash
# Executar via make
make quick-test

# Ou executar diretamente
./scripts/quick-test.sh
```

**Fluxo executado:**
1. Health check dos serviços
2. Cadastrar usuário (João Silva)
3. Cadastrar produto (iPhone 15)
4. Criar pedido (2 unidades)
5. Pagar pedido
6. Listar dados
7. Consultar Query Service

## 📋 Pré-requisitos

### Dependências Obrigatórias
- **curl** - Para fazer requisições HTTP
- **bash** - Para executar os scripts

### Dependências Opcionais
- **jq** - Para formatação JSON (recomendado)

**Instalar jq:**
```bash
# macOS
brew install jq

# Ubuntu/Debian
sudo apt-get install jq

# CentOS/RHEL
sudo yum install jq
```

## 🔧 Configuração

### URLs dos Serviços
Os scripts estão configurados para as seguintes URLs padrão:

```bash
USER_SERVICE="http://localhost:8081"
PRODUCT_SERVICE="http://localhost:8082"
ORDER_SERVICE="http://localhost:8083"
QUERY_SERVICE="http://localhost:8084"
```

### Variáveis de Ambiente
Os scripts não dependem de variáveis de ambiente, mas você pode modificar as URLs diretamente nos arquivos se necessário.

## 🎯 Casos de Uso

### Cenário 1: Teste Inicial
```bash
# 1. Iniciar infraestrutura
make up

# 2. Executar teste rápido
make quick-test
```

### Cenário 2: Teste Interativo
```bash
# 1. Iniciar infraestrutura
make up

# 2. Executar teste interativo
make test-api
```

### Cenário 3: Teste Manual
```bash
# 1. Iniciar infraestrutura
make up

# 2. Executar serviços (em terminais separados)
make run-user
make run-product
make run-order
make run-query

# 3. Executar script
./scripts/test-api.sh
```

## 📊 Exemplos de Saída

### Teste Rápido
```
🚀 TESTE RÁPIDO DA API - Event-Driven Architecture
==================================================
1. Verificando health dos serviços...
✅ User Service - OK
✅ Product Service - OK
✅ Order Service - OK
✅ Query Service - OK

2. Cadastrando usuário...
✅ Usuário cadastrado com sucesso!
User ID: 1

3. Cadastrando produto...
✅ Produto cadastrado com sucesso!
Product ID: 1

4. Criando pedido...
✅ Pedido criado com sucesso!
Order ID: 1

5. Pagando pedido...
✅ Pedido pago com sucesso!

6. Listando dados...
Usuários: 1
Produtos: 1
Pedidos: 1

7. Consultando Query Service...
✅ Query Service funcionando!
1

🎉 Teste concluído com sucesso!
IDs gerados:
User ID: 1
Product ID: 1
Order ID: 1
```

### Teste Interativo
```
╔══════════════════════════════════════════════════════════════╗
║                🚀 TESTE INTERATIVO DA API                    ║
║              Event-Driven Architecture + CQRS                ║
╚══════════════════════════════════════════════════════════════╝

┌─────────────────────────────────────────────────────────────┐
│                        MENU PRINCIPAL                       │
├─────────────────────────────────────────────────────────────┤
│ 1. 📝 Cadastrar Usuário                                     │
│ 2. 📦 Cadastrar Produto                                     │
│ 3. 🛒 Criar Pedido                                          │
│ 4. 💳 Pagar Pedido                                          │
│ 5. ❌ Cancelar Pedido                                        │
│ 6. 📋 Listar Usuários                                       │
│ 7. 📋 Listar Produtos                                       │
│ 8. 📋 Listar Pedidos                                        │
│ 9. 🔍 Consultar Dados (Query Service)                       │
│ 10. 🏥 Health Check dos Serviços                            │
│ 11. 📊 Status dos Dados                                     │
│ 0. 🚪 Sair                                                   │
└─────────────────────────────────────────────────────────────┘

Escolha uma opção:
```

## 🛠️ Troubleshooting

### Erro: "curl não está instalado"
```bash
# macOS
brew install curl

# Ubuntu/Debian
sudo apt-get install curl

# CentOS/RHEL
sudo yum install curl
```

### Erro: "jq não está instalado"
O script funcionará sem jq, mas a formatação JSON será mais simples.

### Erro: "Connection refused"
Verifique se os serviços estão rodando:
```bash
make health
```

### Erro: "Service not found"
Verifique se a infraestrutura está ativa:
```bash
make status
```

## 🔄 Personalização

### Modificar URLs dos Serviços
Edite as variáveis no início dos scripts:
```bash
USER_SERVICE="http://localhost:8081"
PRODUCT_SERVICE="http://localhost:8082"
ORDER_SERVICE="http://localhost:8083"
QUERY_SERVICE="http://localhost:8084"
```

### Adicionar Novos Testes
Para adicionar novos testes, crie uma nova função no script e adicione a opção no menu.

### Modificar Dados de Teste
Edite as variáveis `user_data`, `product_data`, etc. nos scripts para usar dados diferentes.

## 📝 Logs e Debug

### Verbose Mode
Para ver mais detalhes das requisições, modifique a função `make_request`:
```bash
# Adicione -v ao curl para verbose
curl -v -s -w "\n%{http_code}" -X "$method" "$url" \
    -H "Content-Type: application/json" \
    -d "$data"
```

### Logs de Erro
Os scripts mostram erros detalhados quando as requisições falham.

## 🤝 Contribuição

Para adicionar novos scripts ou melhorar os existentes:

1. Crie o script no diretório `scripts/`
2. Adicione permissão de execução: `chmod +x scripts/nome-do-script.sh`
3. Atualize o Makefile se necessário
4. Documente no README.md

---

**💡 Dica:** Use `make quick-test` para verificar rapidamente se tudo está funcionando antes de usar o teste interativo.
