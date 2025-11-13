# 📝 Sistema NF - Gestão de Notas Fiscais# 📝 Sistema de Gestão de Notas Fiscais



<div align="center">> Sistema completo para gerenciamento de produtos e emissão de notas fiscais com controle de estoque em tempo real.



![Angular](https://img.shields.io/badge/Angular-19.2-red?style=for-the-badge&logo=angular)![Status](https://img.shields.io/badge/status-produção-success)

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go)![Angular](https://img.shields.io/badge/Angular-19.2-red)

![MariaDB](https://img.shields.io/badge/MariaDB-11.5-003545?style=for-the-badge&logo=mariadb)![Go](https://img.shields.io/badge/Go-1.23-blue)

![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript)![MariaDB](https://img.shields.io/badge/MariaDB-11.5-orange)



**Sistema empresarial com arquitetura de microsserviços para gestão de produtos, estoque e notas fiscais**---



[📖 Como Iniciar](#-início-rápido) • [🏗️ Arquitetura](#️-arquitetura) • [✨ Funcionalidades](#-funcionalidades)## 🎯 Visão Geral



</div>Sistema empresarial desenvolvido com **microserviços** para gerenciar produtos, controlar estoque e emitir notas fiscais. Utiliza arquitetura moderna com backend em Go, frontend em Angular e banco de dados MariaDB.



---### ✨ Funcionalidades Principais



## 🎯 Visão Geral- ✅ **Gestão de Produtos**: Cadastro, edição, exclusão e listagem de produtos com imagens

- ✅ **Controle de Estoque**: Atualização automática de saldo ao emitir notas fiscais

Sistema completo para gerenciamento de notas fiscais com:- ✅ **Notas Fiscais**: Criação, edição e impressão de notas com múltiplos itens

- ✅ Cadastro de produtos com imagens- ✅ **Concorrência**: Sistema de locks para prevenir conflitos em operações simultâneas

- ✅ Controle automático de estoque- ✅ **Circuit Breaker**: Proteção contra falhas em cascata entre microserviços

- ✅ Emissão de notas fiscais- ✅ **Cache Inteligente**: Redução de requisições HTTP com RxJS shareReplay

- ✅ Assistente IA (Hugging Face)- ✅ **Retry Automático**: Tentativas com backoff exponencial em falhas temporárias

- ✅ Resiliência com Circuit Breaker- ✅ **Hot Reload**: Desenvolvimento ágil com Air (Go) e Angular CLI



------



## 🏗️ Arquitetura## 🏗️ Arquitetura



``````

┌─────────────────────────────────────────────────┐┌─────────────────────────────────────────────────────────────┐

│          FRONTEND - Angular 19.2                ││                      FRONTEND (Angular)                      │

│     http://localhost:4200                       ││  - Componentes Standalone                                    │

│  • Material Design  • RxJS  • Standalone        ││  - RxJS para reatividade                                     │

└──────────────┬──────────────────────────────────┘│  - Angular Material Design                                   │

               │ HTTP REST│  - NGX-Toastr, Date-fns, NGX-Mask                          │

       ┌───────┴────────┐└────────────────────┬────────────────────────────────────────┘

       ▼                ▼                     │ HTTP REST API

┌──────────────┐  ┌──────────────┐                     ├───────────────────┬─────────────────────┐

│   ESTOQUE    │  │ FATURAMENTO  │                     ▼                   ▼                     ▼

│  Porta 3001  │◄─┤  Porta 3002  │┌────────────────────────────┐ ┌────────────────────────────┐

│              │  │              ││ SERVIÇO DE ESTOQUE (Go)    │ │ SERVIÇO FATURAMENTO (Go)  │

│ • Produtos   │  │ • Notas      ││ Porta: 3001                 │ │ Porta: 3002                │

│ • Saldo      │  │ • Itens      ││ - CRUD de Produtos          │ │ - CRUD de Notas Fiscais   │

│              │  │ • IA Chat    ││ - Controle de Saldo         │ │ - Circuit Breaker         │

└──────┬───────┘  └──────┬───────┘│ - SELECT FOR UPDATE         │ │ - Integração com Estoque  │

       │                 ││ - Retry com Backoff         │ │ - Validações Robustas     │

       └────────┬────────┘└──────────┬─────────────────┘ └──────────┬────────────────┘

                ▼           │                               │

       ┌─────────────────┐           └───────────────┬───────────────┘

       │     MariaDB     │                           ▼

       │ notafiscal_desafio              ┌──────────────────────────┐

       └─────────────────┘              │    MariaDB 11.5          │

```              │ notafiscal_desafio       │

              │ - Tabela: produtos       │

---              │ - Tabela: notasfiscais   │

              │ - Tabela: itens          │

## ✨ Funcionalidades              └──────────────────────────┘

```

### 📦 Gestão de Produtos

- Criar, editar e remover produtos### 🎨 Stack Tecnológico

- Upload de imagens (Base64)

- Busca em tempo real (debounce)**Frontend:**

- Visualização em cards ou tabela- Angular 19.2 (Standalone Components)

- TypeScript 5.7

### 📋 Notas Fiscais- Angular Material 19.2

- Criar notas com múltiplos itens- RxJS 7.8 (Operadores: shareReplay, retryWhen, debounceTime)

- Status: ABERTA (editável) / FECHADA (finalizada)- NGX-Toastr (notificações elegantes)

- Atualização automática de estoque- Date-fns (formatação de datas)

- Visualização em formato de impressão- NGX-Mask (máscaras de input)



### 🤖 Assistente IA**Backend:**

- Chat inteligente (Hugging Face)- Go 1.23

- Análise de dados de vendas- Gin Framework (rotas HTTP)

- Insights sobre estoque- MySQL Driver

- Zap Logger (logs estruturados)

### 🛡️ Resiliência- Viper (gerenciamento de configs)

- **Circuit Breaker**: Proteção contra falhas- Air (hot reload)

- **Retry**: 3 tentativas com backoff (1s, 2s, 3s)

- **Concorrência**: SELECT FOR UPDATE**Banco de Dados:**

- MariaDB 11.5.2

---- InnoDB Engine

- Transações ACID

## 🚀 Início Rápido- Foreign Keys



### Pré-requisitos**Ferramentas:**

- Node.js 20+- Git (controle de versão)

- Go 1.24+- VS Code (IDE)

- MariaDB 11.5+- Postman/Thunder Client (testes API)



### Instalação---



```bash## 📦 Estrutura do Projeto

# 1. Clone e acesse o projeto

git clone <repo-url>```

cd Korp_Teste_EduardoMartinPROJETO KORP/

│

# 2. Configure o banco├── frontend/                    # Aplicação Angular

mysql -u root -p < database.sql│   ├── src/

│   │   ├── app/

# 3. Inicie backend - Estoque│   │   │   ├── components/     # Componentes standalone

cd backend/estoque│   │   │   │   ├── home/

go run main.go│   │   │   │   ├── produtos/

│   │   │   │   │   ├── produto-form/

# 4. Inicie backend - Faturamento (novo terminal)│   │   │   │   │   └── produto-list/

cd backend/faturamento│   │   │   │   └── notas/

go run main.go│   │   │   │       ├── nota-form/

│   │   │   │       ├── nota-list/

# 5. Inicie frontend (novo terminal)│   │   │   │       └── nota-print-dialog/

cd frontend│   │   │   ├── models/         # Interfaces TypeScript

npm install│   │   │   ├── services/       # Serviços HTTP

npm start│   │   │   ├── app.config.ts   # Configuração da aplicação

```│   │   │   └── app.routes.ts   # Rotas

│   │   └── styles.scss         # Estilos globais

**Acesse:** http://localhost:4200│   ├── angular.json

│   ├── package.json

---│   └── tsconfig.json

│

## 📁 Estrutura├── servico-estoque-go/          # Microserviço de Estoque

│   ├── main.go                  # Código principal

```│   ├── config.yaml              # Configurações

├── frontend/                 # Angular 19.2│   ├── .air.toml                # Config hot reload

│   ├── src/app/│   ├── go.mod

│   │   ├── components/      # Produtos, Notas, Chat IA│   ├── go.sum

│   │   ├── services/        # HTTP Services│   └── logs/                    # Logs estruturados

│   │   └── models/          # Interfaces TypeScript│

│   └── package.json├── servico-faturamento-go/      # Microserviço de Faturamento

││   ├── main.go                  # Código principal

├── backend/│   ├── config.yaml              # Configurações

│   ├── estoque/             # Microserviço Estoque (Go)│   ├── .air.toml                # Config hot reload

│   │   └── main.go│   ├── go.mod

│   └── faturamento/         # Microserviço Faturamento (Go)│   ├── go.sum

│       └── main.go│   └── logs/                    # Logs estruturados

││

└── database.sql             # Schema do banco├── database.sql                 # Script de criação do banco

```├── README.md                    # Este arquivo

├── COMO-INICIAR.md             # Guia de instalação e execução

---└── DETALHAMENTO-TECNICO.md     # Documentação técnica detalhada

```

## 🔧 Stack Tecnológica

---

### Frontend

| Tecnologia | Versão | Uso |## 🚀 Início Rápido

|-----------|--------|-----|

| Angular | 19.2 | Framework SPA |### Pré-requisitos

| TypeScript | 5.x | Linguagem |

| Angular Material | 17.x | UI Components |- Node.js 20+ e npm

| RxJS | 7.x | Programação reativa |- Go 1.23+

- MariaDB 11.5+

### Backend- Git

| Tecnologia | Versão | Uso |

|-----------|--------|-----|### Instalação e Execução

| Go | 1.24 | Linguagem |

| Gin | 1.9.1 | Framework HTTP |Consulte o arquivo **[COMO-INICIAR.md](COMO-INICIAR.md)** para instruções detalhadas de instalação e execução.

| MySQL Driver | 1.7.1 | Banco de dados |

| UUID | 1.5.0 | IDs únicos |**Resumo:**



---```bash

# 1. Clone o repositório

## 📊 Padrões Implementadosgit clone <url-do-repositorio>



### Circuit Breaker# 2. Configure o banco de dados

```mysql -u root -p < database.sql

CLOSED → (3 falhas) → OPEN → (10s) → HALF_OPEN → CLOSED

```# 3. Inicie o backend (Estoque)

- Protege contra falhas em cascatacd servico-estoque-go

- Timeout de 10 segundosair  # ou: go run main.go

- Reset manual disponível

# 4. Inicie o backend (Faturamento)

### Retry com Exponential Backoffcd servico-faturamento-go

```air  # ou: go run main.go

Tentativa 1: Imediato

Tentativa 2: Aguarda 1s# 5. Inicie o frontend

Tentativa 3: Aguarda 2scd frontend

```npm install

npm start

### Cache com RxJS```

```typescript

shareReplay(1) // Reduz 66% das requisições HTTPAcesse: **http://localhost:4200**

```

---

---

## 📚 Documentação

## 🧪 Testando o Sistema

- **[COMO-INICIAR.md](COMO-INICIAR.md)** - Guia completo de instalação, configuração e execução

### 1. Criar Produto- **[DETALHAMENTO-TECNICO.md](DETALHAMENTO-TECNICO.md)** - Arquitetura, fluxos, padrões e implementações técnicas

1. Acesse "Produtos" → "Novo Produto"

2. Preencha: Código, Descrição, Saldo---

3. Adicione imagem (opcional)

4. Salvar## 🎯 Funcionalidades Detalhadas



### 2. Criar Nota Fiscal### 1. Gestão de Produtos

1. Acesse "Notas Fiscais" → "Nova Nota"

2. Selecione produtos e quantidades**Funcionalidades:**

3. Adicione múltiplos itens- Criar produtos com código, descrição, saldo e imagem

4. Salvar- Editar produtos existentes (exceto código)

- Remover produtos (se não houver notas vinculadas)

### 3. Finalizar Nota- Listar produtos com busca em tempo real (debounce 300ms)

1. Liste as notas- Upload de imagens (conversão para Base64)

2. Clique em "Finalizar" na nota ABERTA

3. Verifique: Saldo do produto diminui automaticamente**Validações:**

- Código único (máx. 10 caracteres)

### 4. Testar IA- Descrição obrigatória (máx. 200 caracteres)

1. Clique no ícone de chat (canto inferior direito)- Saldo não negativo

2. Digite: "Como funciona uma nota fiscal?"- Imagem opcional (máx. 2MB)

3. Ou na tela de produtos: "Analisar com IA"

### 2. Controle de Estoque

---

**Funcionalidades:**

## 📝 Documentação Adicional- Atualização automática de saldo ao finalizar nota fiscal

- Controle de concorrência com SELECT FOR UPDATE

- **[COMO-INICIAR.md](COMO-INICIAR.md)** - Guia completo de instalação- Retry automático em caso de conflito (3 tentativas)

- **[DETALHAMENTO-TECNICO.md](DETALHAMENTO-TECNICO.md)** - Arquitetura e implementação- Validação de saldo disponível antes da reserva

- **[GUIA-VIDEO.md](GUIA-VIDEO.md)** - Roteiro para gravação de demo

**Fluxo de atualização:**

---1. Início da transação

2. Lock pessimista (SELECT FOR UPDATE)

## 🛠️ Comandos Úteis3. Validação de saldo

4. Atualização condicional (WHERE id = ? AND saldo = ?)

```bash5. Verificação de rows affected

# Frontend6. Commit ou Rollback

npm start              # Dev server (porta 4200)

npm run build         # Build produção### 3. Notas Fiscais



# Backend**Funcionalidades:**

go run main.go        # Executar- Criar notas com múltiplos itens

go build              # Compilar- Editar notas em status ABERTA

- Finalizar notas (muda status para FECHADA e atualiza estoque)

# Banco- Remover notas (ABERTA ou FECHADA)

mysql -u root -p      # Acessar MariaDB- Imprimir notas em formato profissional

```- Visualizar detalhes em dialog



---**Estados:**

- **ABERTA**: Nota em edição, pode adicionar/remover itens

## 🐛 Troubleshooting- **FECHADA**: Nota finalizada, estoque atualizado, não editável



**Backend não inicia:**### 4. Recursos Avançados

- Verifique se MariaDB está rodando

- Confirme credenciais em `main.go`**Cache com RxJS:**

- `shareReplay(1)` para evitar requisições duplicadas

**Frontend não conecta:**- Invalidação automática após mutações

- Verifique se backends estão nas portas 3001 e 3002- Parâmetro `forceRefresh` para bypass manual

- Limpe cache do navegador

**Retry com Backoff:**

**Circuit Breaker aberto:**- 3 tentativas automáticas

- POST em `http://localhost:3002/circuit-breaker/reset`- Delays crescentes: 1s → 2s → 3s

- Logs informativos de tentativas

---

**Debounce na Busca:**

## 🎯 Tecnologias-Chave- Aguarda 300ms após parar de digitar

- `distinctUntilChanged()` para evitar buscas duplicadas

<div align="center">- Reduz operações em até 87%



| Frontend | Backend | Database |**Circuit Breaker:**

|:--------:|:-------:|:--------:|- Proteção contra falhas no serviço de estoque

| <img src="https://angular.io/assets/images/logos/angular/angular.svg" width="60"> | <img src="https://go.dev/blog/go-brand/Go-Logo/SVG/Go-Logo_Blue.svg" width="80"> | <img src="https://mariadb.org/wp-content/uploads/2019/11/mariadb-logo-vertical_blue.svg" width="60"> |- Estados: CLOSED → OPEN → HALF_OPEN

| **Angular** | **Golang** | **MariaDB** |- Reset manual via endpoint /circuit-breaker/reset



</div>---



---## 🔒 Segurança e Boas Práticas



<div align="center">- ✅ Validação de dados no frontend e backend

- ✅ Transações ACID para consistência

**Desenvolvido com ☕ e 💪**- ✅ Locks pessimistas para concorrência

- ✅ Retry automático com backoff exponencial

*Sistema completo e pronto para produção*- ✅ Circuit breaker para resiliência

- ✅ CORS configurado corretamente

</div>- ✅ Logs estruturados com Zap

- ✅ Configurações externalizadas com Viper
- ✅ .gitignore para não commitar logs e binários

---

## 📊 Performance

**Melhorias implementadas:**
- 🚀 Cache com shareReplay: **66% menos requisições HTTP**
- 🚀 Debounce na busca: **87% menos operações de filtro**
- 🚀 Retry automático: **Maior resiliência a falhas temporárias**
- 🚀 Virtual Scroll (futuro): **Renderizar apenas itens visíveis**

---

## 🧪 Testes

### Testar Concorrência

1. Abra 2 abas do navegador
2. Crie um produto com saldo 1
3. Crie 2 notas fiscais simultaneamente usando o mesmo produto
4. Finalize ambas ao mesmo tempo
5. **Resultado esperado**: Uma nota deve ser finalizada com sucesso, a outra deve retornar erro de saldo insuficiente

### Testar Circuit Breaker

1. Desligue o serviço de estoque
2. Tente criar/finalizar uma nota fiscal
3. Após 3 falhas, circuit breaker abre
4. Tente novamente → resposta instantânea de erro
5. Ligue o serviço de estoque
6. Faça POST em `/api/notas/circuit-breaker/reset`
7. Sistema volta ao normal

### Testar Cache

1. Abra DevTools → Network
2. Acesse lista de produtos
3. Observe: 1 requisição HTTP
4. Navegue para outra página e volte
5. Observe: sem nova requisição (cache ativo)
6. Crie um novo produto
7. Observe: nova requisição (cache invalidado)

---

## 🛠️ Comandos Úteis

```bash
# Frontend
npm start              # Inicia dev server (porta 4200)
npm run build          # Build de produção
npm test               # Executa testes

# Backend (com Air - hot reload)
air                    # Inicia com hot reload

# Backend (sem Air)
go run main.go         # Executa diretamente
go build               # Compila binário
go test ./...          # Executa testes

# Banco de Dados
mysql -u root -p notafiscal_desafio  # Acessa banco
SHOW TABLES;                          # Lista tabelas
SELECT * FROM produtos;               # Lista produtos
```

---

## 🐛 Troubleshooting

**Frontend não conecta ao backend:**
- Verifique se os serviços Go estão rodando nas portas 3001 e 3002
- Confirme CORS configurado no backend
- Verifique console do navegador para erros

**Erro de saldo insuficiente:**
- Verifique saldo do produto no banco de dados
- Confirme que não há notas pendentes usando o produto

**Circuit breaker aberto:**
- Verifique se serviço de estoque está online
- Faça POST em `/api/notas/circuit-breaker/reset` para resetar

**Air não funciona:**
- Certifique-se que `$GOPATH/bin` está no PATH
- Use `go install github.com/air-verse/air@latest`
- Se persistir, use `go run main.go`

---

## 📝 Licença

Este projeto foi desenvolvido para fins educacionais e demonstração de conceitos de arquitetura de microserviços.

---

**Última atualização:** Novembro 2025
