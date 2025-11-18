# 📊 Sistema NF - Gestão Completa de Notas Fiscais

**Sistema empresarial com arquitetura de microsserviços**  
**Gestão de produtos, estoque e notas fiscais em tempo real**

<div align="center">

![Status](https://img.shields.io/badge/status-produção-success?style=for-the-badge)
![Angular](https://img.shields.io/badge/Angular-19.2-DD0031?style=for-the-badge&logo=angular)
![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go)
![MariaDB](https://img.shields.io/badge/MariaDB-11.5-003545?style=for-the-badge&logo=mariadb)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript)

</div>

---

## 🎯 Visão Geral

Sistema completo desenvolvido com **arquitetura de microsserviços** para gerenciar produtos, controlar estoque em tempo real e emitir notas fiscais.  
Utiliza tecnologias modernas com **backend em Go**, **frontend em Angular** e **banco de dados MariaDB**.

---

## ✨ Funcionalidades Principais

### 📦 Gestão de Produtos
- ✅ **Cadastro completo** — CRUD de produtos  
- ✅ **Upload de imagens** — Base64 (máx. 2MB)  
- ✅ **Busca inteligente** — Filtro com debounce  
- ✅ **Visualização flexível** — Cards ou tabela  

### 📋 Notas Fiscais
- ✅ **Emissão completa** — Múltiplos itens  
- ✅ **Status dinâmico** — ABERTA / FECHADA  
- ✅ **Atualização automática** — Estoque em tempo real  
- ✅ **Formato profissional** — Visualização para impressão  

### 🛡️ Sistema Resiliente
- ✅ **Circuit Breaker** — Proteção contra falhas  
- ✅ **Retry automático** — Backoff exponencial  
- ✅ **Controle de concorrência** — SELECT FOR UPDATE  
- ✅ **Cache inteligente** — Redução de 66% nas requisições  

### 🤖 Assistente IA
- ✅ **Chat inteligente** — Integração com Hugging Face API  
- ✅ **Análise de dados** — Insights de vendas  
- ✅ **Processamento natural** — NLP avançado  

---

## 🚀 Início Rápido

### Pré-requisitos
- Node.js 20+
- Go 1.23+
- MariaDB 11.5+
- Git

### Instalação

```bash
# 1. Clone o repositório
git clone https://github.com/eduardomartinDev/SISTEMA_NF.git
cd SISTEMA_NF

# 2. Configure o banco de dados
mysql -u root -p < database.sql

# 3. Inicie o serviço de Estoque (Terminal 1)
cd backend/estoque
air  # ou: go run main.go

# 4. Inicie o serviço de Faturamento (Terminal 2)
cd backend/faturamento
air  # ou: go run main.go

# 5. Inicie o Frontend (Terminal 3)
cd frontend
npm install
npm start

# 6. Acesso
Acesse: http://localhost:4200

```

### Estrutura do Projeto

```bash
Korp_Teste_EduardoMartin/
├── frontend/                    # Aplicação Angular
│   ├── src/app/
│   │   ├── components/          # Componentes standalone
│   │   ├── services/            # Serviços HTTP
│   │   └── models/              # Interfaces TypeScript
│   └── package.json
│
├── backend/
│   ├── estoque/                 # Microsserviço de Estoque
│   │   ├── main.go
│   │   ├── config.yaml
│   │   └── .air.toml
│   │
│   └── faturamento/             # Microsserviço de Faturamento
│       ├── main.go
│       ├── config.yaml
│       └── .air.toml
│
├── database.sql                 # Schema do banco
└── README.md
```

## 🔒 Segurança e Boas Práticas

- ✅ **Validação de dados** em frontend e backend  
- ✅ **Transações ACID** para consistência  
- ✅ **Locks pessimistas** para concorrência  
- ✅ **CORS** configurado corretamente  
- ✅ **Logs estruturados** com Zap  
- ✅ **Configurações externalizadas**

---

## 📊 Performance

```bash
| Otimização | Resultado |
|:-----------|:----------|
| Cache com `shareReplay` | 🚀 **-66%** requisições HTTP |
| Debounce na busca | 🚀 **-87%** operações de filtro |
| Retry automático | ✅ Resiliência a falhas |
| Circuit Breaker | ✅ Proteção do sistema |
```

---

## 🧪 Testes do Sistema

### ⚡ Testar Concorrência
1. Crie produto com saldo `1`  
2. Tente finalizar `2 notas` simultaneamente  
3. ✅ **Resultado:** Uma nota sucede, outra falha por saldo insuficiente  

### 🔌 Testar Circuit Breaker
1. Desligue o serviço de estoque  
2. Tente operações → Circuit Breaker abre após **3 falhas**  
3. Ligue o serviço e faça **reset via endpoint**

### 💾 Testar Cache
1. Acesse lista de produtos (**1 requisição**)  
2. Navegue e volte (**0 requisições — cache ativo**)  
3. Crie produto (**cache invalidado automaticamente**)

---

## 🛠️ Comandos Úteis

```bash
# Frontend
npm start       # Servidor de desenvolvimento
npm run build   # Build de produção

# Backend
air             # Hot reload
go run main.go  # Execução direta

# Banco de Dados
mysql -u root -p notafiscal_desafio
```

## 🐛 Troubleshooting

| Problema | Solução |
|:----------|:----------|
| **Frontend não conecta** | Verifique se os serviços estão ativos nas portas `3001` e `3002`. |
| **Erro de saldo insuficiente** | Confirme o saldo disponível no banco de dados antes de finalizar a nota. |
| **Circuit Breaker aberto** | Execute um `POST` em `/circuit-breaker/reset` para reativar o serviço. |
| **Air não funciona** | Use `go run main.go` como alternativa de execução. |

-- Última atualização: Novembro 2025