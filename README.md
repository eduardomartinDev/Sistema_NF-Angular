# 📊 Sistema NF - Gestão Completa de Notas Fiscais# 📊 Sistema NF - Gestão Completa de Notas Fiscais# 📊 Sistema NF - Gestão Completa de Notas Fiscais



<div align="center">



**Sistema empresarial com arquitetura de microsserviços**  <div align="center"><div align="center">

**Gestão de produtos, estoque e notas fiscais em tempo real**



<br>

**Sistema empresarial com arquitetura de microsserviços para gestão de produtos, estoque e notas fiscais****Sistema empresarial com arquitetura de microsserviços para gestão de produtos, estoque e notas fiscais**

![Status](https://img.shields.io/badge/status-produção-success?style=for-the-badge)

![Angular](https://img.shields.io/badge/Angular-19.2-DD0031?style=for-the-badge&logo=angular)

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go)

![MariaDB](https://img.shields.io/badge/MariaDB-11.5-003545?style=for-the-badge&logo=mariadb)![Status](https://img.shields.io/badge/status-produção-success?style=for-the-badge)![Status](https://img.shields.io/badge/status-produção-success?style=for-the-badge)

![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript)

![Angular](https://img.shields.io/badge/Angular-19.2-DD0031?style=for-the-badge&logo=angular)![Angular](https://img.shields.io/badge/Angular-19.2-DD0031?style=for-the-badge&logo=angular)

</div>

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go)![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go)

<br>

![MariaDB](https://img.shields.io/badge/MariaDB-11.5-003545?style=for-the-badge&logo=mariadb)![MariaDB](https://img.shields.io/badge/MariaDB-11.5-003545?style=for-the-badge&logo=mariadb)

## 🎯 Visão Geral

![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript)![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript)

Sistema completo desenvolvido com **arquitetura de microsserviços** para gerenciar produtos, controlar estoque em tempo real e emitir notas fiscais. Utiliza tecnologias modernas com backend em Go, frontend em Angular e banco de dados MariaDB.



<br>

</div></div>

## ✨ Funcionalidades Principais



<table>

<tr>---## 🎯 Visão Geral

<td width="50%">



### 📦 Gestão de Produtos

✅ **Cadastro completo** - CRUD de produtos  ## 🎯 Visão GeralSistema completo desenvolvido com **arquitetura de microsserviços** para gerenciar produtos, controlar estoque em tempo real e emitir notas fiscais. Utiliza tecnologias modernas com backend em Go, frontend em Angular e banco de dados MariaDB.

✅ **Upload de imagens** - Base64 (máx. 2MB)  

✅ **Busca inteligente** - Filtro com debounce  

✅ **Visualização flexível** - Cards ou tabela  

Sistema completo desenvolvido com **arquitetura de microsserviços** para gerenciar produtos, controlar estoque em tempo real e emitir notas fiscais. Utiliza tecnologias modernas com backend em Go, frontend em Angular e banco de dados MariaDB.---

### 📋 Notas Fiscais

✅ **Emissão completa** - Múltiplos itens  

✅ **Status dinâmico** - ABERTA / FECHADA  

✅ **Atualização automática** - Estoque em tempo real  ---## ✨ Funcionalidades Principais

✅ **Formato profissional** - Visualização para impressão  



</td>

<td width="50%">## ✨ Funcionalidades Principais### 📦 Gestão de Produtos



### 🛡️ Sistema Resiliente- ✅ **Cadastro completo**: Criar, editar, visualizar e remover produtos

✅ **Circuit Breaker** - Proteção contra falhas  

✅ **Retry automático** - Backoff exponencial  ### 📦 Gestão de Produtos- ✅ **Upload de imagens**: Suporte a Base64 com validação (máx. 2MB)

✅ **Controle de concorrência** - SELECT FOR UPDATE  

✅ **Cache inteligente** - 66% menos requisições  - ✅ **Cadastro completo**: Criar, editar, visualizar e remover produtos- ✅ **Busca inteligente**: Filtro em tempo real com debounce



### 🤖 Assistente IA- ✅ **Upload de imagens**: Suporte a Base64 com validação (máx. 2MB)- ✅ **Visualização flexível**: Cards ou tabela conforme preferência

✅ **Chat inteligente** - Hugging Face API  

✅ **Análise de dados** - Insights de vendas  - ✅ **Busca inteligente**: Filtro em tempo real com debounce

✅ **Processamento natural** - NLP avançado  

- ✅ **Visualização flexível**: Cards ou tabela conforme preferência### 📋 Notas Fiscais

</td>

</tr>- ✅ **Emissão completa**: Criar notas com múltiplos itens

</table>

### 📋 Notas Fiscais- ✅ **Status dinâmico**: ABERTA (editável) / FECHADA (finalizada)

<br>

- ✅ **Emissão completa**: Criar notas com múltiplos itens- ✅ **Atualização automática**: Estoque atualizado em tempo real

## 🏗️ Arquitetura do Sistema

- ✅ **Status dinâmico**: ABERTA (editável) / FECHADA (finalizada)- ✅ **Formato profissional**: Visualização para impressão

```

┌─────────────────────────────────────────────────────────────┐- ✅ **Atualização automática**: Estoque atualizado em tempo real

│                    FRONTEND (Angular)                        │

│         • Componentes Standalone                            │- ✅ **Formato profissional**: Visualização para impressão### 🛡️ Sistema Resiliente

│         • RxJS para reatividade                             │

│         • Angular Material Design                           │- ✅ **Circuit Breaker**: Proteção contra falhas em cascata

└────────────────────┬────────────────────────────────────────┘

                     │ HTTP REST API### 🛡️ Sistema Resiliente- ✅ **Retry automático**: Tentativas com backoff exponencial

     ┌───────────────┼───────────────┐

     ▼               ▼               ▼- ✅ **Circuit Breaker**: Proteção contra falhas em cascata- ✅ **Controle de concorrência**: SELECT FOR UPDATE para transações seguras

┌────────────────┐ ┌────────────────┐ ┌────────────────┐

│   ESTOQUE      │ │  FATURAMENTO   │ │      IA        │- ✅ **Retry automático**: Tentativas com backoff exponencial- ✅ **Cache inteligente**: Redução de 66% nas requisições HTTP

│  Porta: 3001   │ │  Porta: 3002   │ │  Hugging Face  │

│ • CRUD Produtos│ │ • Notas Fiscais│ │ • Chat         │- ✅ **Controle de concorrência**: SELECT FOR UPDATE para transações seguras

│ • Controle     │ │ • Circuit Br.  │ │ • Analytics    │

│   Saldo        │ │ • Validações   │ │                │- ✅ **Cache inteligente**: Redução de 66% nas requisições HTTP### 🤖 Assistente IA

└────────┬───────┘ └───────┬────────┘ └────────────────┘

         │                 │- ✅ **Chat inteligente**: Integração com Hugging Face

         └────────┬─────────┘

                  ▼### 🤖 Assistente IA- ✅ **Análise de dados**: Insights sobre vendas e estoque

         ┌──────────────────┐

         │   MariaDB 11.5   │- ✅ **Chat inteligente**: Integração com Hugging Face- ✅ **Processamento natural**: Compreensão de linguagem natural

         │ • produtos       │

         │ • notasfiscais   │- ✅ **Análise de dados**: Insights sobre vendas e estoque

         │ • itens          │

         └──────────────────┘- ✅ **Processamento natural**: Compreensão de linguagem natural---

```



<br>

---## 🏗️ Arquitetura do Sistema

## 🔧 Stack Tecnológica

┌─────────────────────────────────────────────────────────────┐

| Camada | Tecnologias |

|:------:|-------------|## 🏗️ Arquitetura do Sistema│ FRONTEND (Angular) │

| **Frontend** | Angular 19.2 (Standalone) · TypeScript 5.7 · RxJS 7.8 · Angular Material 19.2 · NGX-Toastr · Date-fns · NGX-Mask |

| **Backend** | Go 1.23 · Gin Framework · Zap Logger · Viper · Air (hot reload) |│ - Componentes Standalone │

| **Banco de Dados** | MariaDB 11.5.2 · InnoDB Engine · Transações ACID |

```│ - RxJS para reatividade │

<br>

┌─────────────────────────────────────────────────────────────┐│ - Angular Material Design │

## 🚀 Início Rápido

│                    FRONTEND (Angular)                        ││ - NGX-Toastr, Date-fns, NGX-Mask │

### Pré-requisitos

- Node.js 20+ e npm│         - Componentes Standalone                            │└────────────────────┬────────────────────────────────────────┘

- Go 1.23+

- MariaDB 11.5+│         - RxJS para reatividade                             ││ HTTP REST API

- Git

│         - Angular Material Design                           │┌───────────────┼───────────────┐

### Instalação

│         - NGX-Toastr, Date-fns, NGX-Mask                   │▼ ▼ ▼

```bash

# 1. Clone o repositório└────────────────────┬────────────────────────────────────────┘┌────────────────┐ ┌────────────────┐ ┌────────────────┐

git clone https://github.com/eduardomartinDev/Korp_Teste_EduardoMartin.git

cd Korp_Teste_EduardoMartin                     │ HTTP REST API│SERVIÇO ESTOQUE│ │SER. FATURAMENTO│ │ ASSISTENTE │



# 2. Configure o banco de dados     ┌───────────────┼───────────────┐│ Porta: 3001 │ │ Porta: 3002 │ │ IA │

mysql -u root -p < database.sql

     ▼               ▼               ▼│ - CRUD Produtos│ │ - Notas Fiscais│ │ - Hugging Face │

# 3. Inicie o serviço de Estoque (Terminal 1)

cd backend/estoque┌────────────────┐ ┌────────────────┐ ┌────────────────┐│ - Controle Saldo│ │ - Circuit Br. │ │ - Analytics │

air  # ou: go run main.go

│SERVIÇO ESTOQUE │ │SER. FATURAMENTO│ │  ASSISTENTE    ││ - SELECT UPDATE │ │ - Validações │ │ - Chat │

# 4. Inicie o serviço de Faturamento (Terminal 2)

cd backend/faturamento│  Porta: 3001   │ │  Porta: 3002   │ │      IA        │└────────┬───────┘ └───────┬────────┘ └────────────────┘

air  # ou: go run main.go

│ - CRUD Produtos│ │ - Notas Fiscais│ │ - Hugging Face ││ │

# 5. Inicie o Frontend (Terminal 3)

cd frontend│ - Controle Saldo│ │ - Circuit Br. │ │ - Analytics    │└────────┬─────────┘

npm install

npm start│ - SELECT UPDATE │ │ - Validações   │ │ - Chat         │▼

```

└────────┬───────┘ └───────┬────────┘ └────────────────┘┌──────────────────┐

**🌐 Acesse:** http://localhost:4200

         │                 ││ MariaDB 11.5 │

<br>

         └────────┬─────────┘│ notafiscal_desafio│

## 📁 Estrutura do Projeto

                  ▼│ - produtos │

```

Korp_Teste_EduardoMartin/         ┌──────────────────┐│ - notasfiscais │

│

├── frontend/                    # Aplicação Angular         │   MariaDB 11.5   ││ - itens │

│   ├── src/app/

│   │   ├── components/         # Componentes standalone         │notafiscal_desafio│└──────────────────┘

│   │   ├── services/           # Serviços HTTP

│   │   └── models/             # Interfaces TypeScript         │   - produtos     │

│   └── package.json

│         │   - notasfiscais │text

├── backend/

│   ├── estoque/                # Microsserviço de Estoque         │   - itens        │

│   │   ├── main.go

│   │   ├── config.yaml         └──────────────────┘### 🔧 Stack Tecnológica

│   │   └── .air.toml

│   │```

│   └── faturamento/            # Microsserviço de Faturamento

│       ├── main.go**Frontend:**

│       ├── config.yaml

│       └── .air.toml### 🔧 Stack Tecnológica- Angular 19.2 (Standalone Components)

│

├── database.sql                # Schema do banco- TypeScript 5.7 + RxJS 7.8

└── README.md

```| Camada | Tecnologias |- Angular Material 19.2



<br>|--------|-------------|- NGX-Toastr, Date-fns, NGX-Mask



## 🔒 Segurança e Boas Práticas| **Frontend** | Angular 19.2 (Standalone), TypeScript 5.7, RxJS 7.8, Angular Material 19.2, NGX-Toastr, Date-fns, NGX-Mask |



- ✅ Validação de dados em frontend e backend| **Backend** | Go 1.23, Gin Framework, Zap Logger, Viper, Air (hot reload) |**Backend:**

- ✅ Transações ACID para consistência

- ✅ Locks pessimistas para controle de concorrência| **Banco de Dados** | MariaDB 11.5.2, InnoDB Engine, Transações ACID |- Go 1.23 + Gin Framework

- ✅ CORS configurado corretamente

- ✅ Logs estruturados com Zap- Zap Logger (logs estruturados)

- ✅ Configurações externalizadas

---- Viper (configurações)

<br>

- Air (hot reload)

## 📊 Performance

## 🚀 Início Rápido

| Otimização | Resultado |

|:-----------|----------:|**Banco de Dados:**

| Cache com shareReplay | 🚀 **-66%** requisições HTTP |

| Debounce na busca | 🚀 **-87%** operações de filtro |### Pré-requisitos- MariaDB 11.5.2

| Retry automático | ✅ Resiliência a falhas |

| Circuit Breaker | ✅ Proteção do sistema |- Node.js 20+ e npm- InnoDB Engine



<br>- Go 1.23+- Transações ACID



## 🧪 Testes do Sistema- MariaDB 11.5+



### ⚡ Testar Concorrência- Git---

1. Crie produto com saldo `1`

2. Tente finalizar `2 notas` simultaneamente

3. ✅ **Resultado:** Uma nota sucede, outra falha por saldo insuficiente

### Instalação e Execução## 🚀 Início Rápido

### 🔌 Testar Circuit Breaker

1. Desligue serviço de estoque

2. Tente operações → Circuit Breaker abre após 3 falhas

3. Ligue serviço e reset via endpoint```bash### Pré-requisitos



### 💾 Testar Cache# 1. Clone o repositório- Node.js 20+ e npm

1. Acesse lista de produtos (1 requisição)

2. Navegue e volte (0 requisições - cache ativo)git clone https://github.com/eduardomartinDev/Korp_Teste_EduardoMartin.git- Go 1.23+

3. Crie produto (cache invalidado automaticamente)

cd Korp_Teste_EduardoMartin- MariaDB 11.5+

<br>

- Git

## 🛠️ Comandos Úteis

# 2. Configure o banco de dados

```bash

# Desenvolvimento Frontendmysql -u root -p < database.sql### Instalação e Execução

npm start              # Servidor dev (porta 4200)

npm run build          # Build produção



# Desenvolvimento Backend# 3. Inicie o serviço de Estoque```bash

air                    # Hot reload

go run main.go         # Execução diretacd backend/estoque# 1. Clone o repositório



# Banco de Dadosair  # ou: go run main.gogit clone <url-do-repositorio>

mysql -u root -p notafiscal_desafio

```cd Korp_Teste_EduardoMartin



<br># 4. Em novo terminal, inicie o serviço de Faturamento



## 🐛 Troubleshootingcd backend/faturamento# 2. Configure o banco de dados



| Problema | Solução |air  # ou: go run main.gomysql -u root -p < database.sql

|----------|---------|

| Frontend não conecta | Verifique serviços nas portas `3001` e `3002` |

| Erro de saldo insuficiente | Confirme saldo disponível no banco |

| Circuit Breaker aberto | POST em `/circuit-breaker/reset` |# 5. Em novo terminal, inicie o Frontend# 3. Inicie o serviço de Estoque

| Air não funciona | Use `go run main.go` como alternativa |

cd frontendcd servico-estoque-go

<br>

npm installair  # ou: go run main.go

## 📚 Documentação

npm start

- 📖 **[COMO-INICIAR.md](COMO-INICIAR.md)** - Guia completo de instalação

- 🔧 **[DETALHAMENTO-TECNICO.md](DETALHAMENTO-TECNICO.md)** - Documentação técnica detalhada```# 4. Em novo terminal, inicie o serviço de Faturamento



<br>cd servico-faturamento-go



---**Acesse:** http://localhost:4200air  # ou: go run main.go



<div align="center">



**Desenvolvido com ☕ e 💪**---# 5. Em novo terminal, inicie o Frontend



Sistema completo e pronto para produçãocd frontend



<br>## 📁 Estrutura do Projetonpm install



*Última atualização: Novembro 2025*npm start



</div>```Acesse: http://localhost:4200


Korp_Teste_EduardoMartin/

│📁 Estrutura do Projeto

├── frontend/                    # Aplicação Angulartext

│   ├── src/app/PROJETO KORP/

│   │   ├── components/         # Componentes standalone│

│   │   │   ├── home/├── frontend/                    # Aplicação Angular

│   │   │   ├── produtos/│   ├── src/app/

│   │   │   └── notas/│   │   ├── components/         # Componentes standalone

│   │   ├── services/           # Serviços HTTP│   │   │   ├── home/

│   │   └── models/             # Interfaces TypeScript│   │   │   ├── produtos/

│   └── package.json│   │   │   └── notas/

││   │   ├── services/           # Serviços HTTP

├── backend/│   │   └── models/             # Interfaces TypeScript

│   ├── estoque/                # Microsserviço de Estoque│   └── package.json

│   │   ├── main.go│

│   │   ├── config.yaml├── servico-estoque-go/         # Microsserviço de Estoque

│   │   └── .air.toml│   ├── main.go

│   ││   ├── config.yaml

│   └── faturamento/            # Microsserviço de Faturamento│   └── .air.toml

│       ├── main.go│

│       ├── config.yaml├── servico-faturamento-go/     # Microsserviço de Faturamento

│       └── .air.toml│   ├── main.go

││   ├── config.yaml

├── database.sql                # Schema do banco│   └── .air.toml

├── README.md│

├── COMO-INICIAR.md             # Guia detalhado├── database.sql                # Schema do banco

└── DETALHAMENTO-TECNICO.md     # Documentação técnica├── README.md

```├── COMO-INICIAR.md             # Guia detalhado

└── DETALHAMENTO-TECNICO.md     # Documentação técnica

---🔒 Segurança e Boas Práticas

✅ Validação de dados em frontend e backend

## 🔒 Segurança e Boas Práticas

✅ Transações ACID para consistência

- ✅ Validação de dados em frontend e backend

- ✅ Transações ACID para consistência✅ Locks pessimistas para controle de concorrência

- ✅ Locks pessimistas para controle de concorrência

- ✅ CORS configurado corretamente✅ CORS configurado corretamente

- ✅ Logs estruturados com Zap

- ✅ Configurações externalizadas✅ Logs estruturados com Zap



---✅ Configurações externalizadas



## 📊 Performance📊 Performance

Otimizações implementadas:

**Otimizações implementadas:**

🚀 Cache com shareReplay: 66% menos requisições HTTP

| Recurso | Melhoria |

|---------|----------|🚀 Debounce na busca: 87% menos operações de filtro

| Cache com shareReplay | 66% menos requisições HTTP |

| Debounce na busca | 87% menos operações de filtro |🚀 Retry automático: Maior resiliência a falhas

| Retry automático | Maior resiliência a falhas |

| Circuit Breaker | Proteção do sistema |🚀 Circuit Breaker: Proteção do sistema



---🧪 Testes do Sistema

Testar Concorrência

## 🧪 Testes do SistemaCrie produto com saldo 1



### Testar ConcorrênciaTente finalizar 2 notas simultaneamente

1. Crie produto com saldo 1

2. Tente finalizar 2 notas simultaneamenteResultado: Uma nota sucede, outra falha por saldo insuficiente

3. **Resultado:** Uma nota sucede, outra falha por saldo insuficiente

Testar Circuit Breaker

### Testar Circuit BreakerDesligue serviço de estoque

1. Desligue serviço de estoque

2. Tente operações → Circuit Breaker abre após 3 falhasTente operações → Circuit Breaker abre após 3 falhas

3. Ligue serviço e reset via endpoint

Ligue serviço e reset via endpoint

### Testar Cache

1. Acesse lista de produtos (1 requisição)Testar Cache

2. Navegue e volte (0 requisições - cache ativo)Acesse lista de produtos (1 requisição)

3. Crie produto (cache invalidado automaticamente)

Navegue e volte (0 requisições - cache ativo)

---

Crie produto (cache invalidado automaticamente)

## 🛠️ Comandos Úteis

🛠️ Comandos Úteis

```bashbash

# Desenvolvimento Frontend# Desenvolvimento Frontend

npm start              # Servidor dev (porta 4200)npm start              # Servidor dev (porta 4200)

npm run build          # Build produçãonpm run build          # Build produção



# Desenvolvimento Backend# Desenvolvimento Backend

air                    # Hot reloadair                    # Hot reload

go run main.go         # Execução diretago run main.go         # Execução direta



# Banco de Dados# Banco de Dados

mysql -u root -p notafiscal_desafiomysql -u root -p notafiscal_desafio

```🐛 Troubleshooting

Problema comum	Solução

---Frontend não conecta	Verifique serviços nas portas 3001/3002

Erro de saldo insuficiente	Confirme saldo disponível no banco

## 🐛 TroubleshootingCircuit Breaker aberto	POST em /circuit-breaker/reset

Air não funciona	Use go run main.go como alternativa

| Problema | Solução |📚 Documentação

|----------|---------|

| Frontend não conecta | Verifique serviços nas portas 3001/3002 |<div align="center">

| Erro de saldo insuficiente | Confirme saldo disponível no banco |Desenvolvido com ☕ e 💪

| Circuit Breaker aberto | POST em `/circuit-breaker/reset` |Sistema completo e pronto para produção

| Air não funciona | Use `go run main.go` como alternativa |

</div>

---Última atualização: Novembro 2025

</div>
