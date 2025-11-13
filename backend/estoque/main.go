package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Produto struct {
	ID           string  `json:"id"`                           // UUID único
	Codigo       string  `json:"codigo" binding:"required"`    // Código do produto (obrigatório)
	Descricao    string  `json:"descricao" binding:"required"` // Descrição do produto (obrigatório)
	Saldo        int     `json:"saldo" binding:"required"`     // Quantidade em estoque (obrigatório)
	ImagemUrl    *string `json:"imagemUrl,omitempty"`          // Imagem em base64 (opcional)
	CriadoEm     string  `json:"criadoEm,omitempty"`           // Data de criação
	AtualizadoEm string  `json:"atualizadoEm,omitempty"`       // Data de atualização
}
type AtualizarSaldoRequest struct {
	Quantidade int `json:"quantidade" binding:"required"`
}
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

var db *sql.DB

func inicializarBanco() error {
	dsn := "root:@tcp(localhost:3306)/notafiscal_desafio?parseTime=true"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("erro ao conectar no banco: %v", err)
	}
	db.SetMaxOpenConns(10)   // Máximo de conexões abertas
	db.SetMaxIdleConns(5)    // Máximo de conexões ociosas
	db.SetConnMaxLifetime(0) // Tempo de vida das conexões (0 = ilimitado)

	log.Println("✅ Conexão com MariaDB estabelecida com sucesso!")
	return nil
}
func listarProdutos(c *gin.Context) {
	log.Println("📦 Listando todos os produtos do banco de dados")
	query := `
		SELECT id, codigo, descricao, saldo, imagem_url, criado_em, atualizado_em 
		FROM produtos 
		ORDER BY codigo
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("❌ Erro ao buscar produtos: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar produtos",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()
	produtos := []Produto{}
	for rows.Next() {
		var p Produto
		err := rows.Scan(
			&p.ID,
			&p.Codigo,
			&p.Descricao,
			&p.Saldo,
			&p.ImagemUrl,
			&p.CriadoEm,
			&p.AtualizadoEm,
		)
		if err != nil {
			log.Printf("❌ Erro ao ler produto: %v", err)
			continue
		}
		produtos = append(produtos, p)
	}

	log.Printf("✅ %d produtos encontrados", len(produtos))
	c.JSON(http.StatusOK, produtos)
}
func buscarProduto(c *gin.Context) {
	id := c.Param("id")
	log.Printf("📦 Buscando produto com ID: %s", id)
	query := `
		SELECT id, codigo, descricao, saldo, imagem_url, criado_em, atualizado_em 
		FROM produtos 
		WHERE id = ?
	`

	var p Produto
	err := db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Codigo,
		&p.Descricao,
		&p.Saldo,
		&p.ImagemUrl,
		&p.CriadoEm,
		&p.AtualizadoEm,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Produto não encontrado",
			Message: fmt.Sprintf("Não existe produto com ID %s", id),
		})
		return
	}

	if err != nil {
		log.Printf("❌ Erro ao buscar produto: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar produto",
			Message: err.Error(),
		})
		return
	}

	log.Printf("✅ Produto encontrado: %s", p.Descricao)
	c.JSON(http.StatusOK, p)
}
func criarProduto(c *gin.Context) {
	log.Println("➕ Criando novo produto")

	var produto Produto
	if err := c.ShouldBindJSON(&produto); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: "Código, descrição e saldo são obrigatórios",
		})
		return
	}
	var count int
	checkQuery := "SELECT COUNT(*) FROM produtos WHERE codigo = ?"
	err := db.QueryRow(checkQuery, produto.Codigo).Scan(&count)
	if err != nil {
		log.Printf("❌ Erro ao verificar código: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao verificar código",
			Message: err.Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Código duplicado",
			Message: fmt.Sprintf("Já existe um produto com o código %s", produto.Codigo),
		})
		return
	}
	produto.ID = uuid.New().String()
	insertQuery := `
		INSERT INTO produtos (id, codigo, descricao, saldo, imagem_url) 
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = db.Exec(insertQuery,
		produto.ID,
		produto.Codigo,
		produto.Descricao,
		produto.Saldo,
		produto.ImagemUrl,
	)

	if err != nil {
		log.Printf("❌ Erro ao criar produto: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao criar produto",
			Message: err.Error(),
		})
		return
	}
	selectQuery := `
		SELECT id, codigo, descricao, saldo, imagem_url, criado_em, atualizado_em 
		FROM produtos 
		WHERE id = ?
	`

	var novoProduto Produto
	err = db.QueryRow(selectQuery, produto.ID).Scan(
		&novoProduto.ID,
		&novoProduto.Codigo,
		&novoProduto.Descricao,
		&novoProduto.Saldo,
		&novoProduto.ImagemUrl,
		&novoProduto.CriadoEm,
		&novoProduto.AtualizadoEm,
	)

	if err != nil {
		log.Printf("❌ Erro ao buscar produto criado: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar produto criado",
			Message: err.Error(),
		})
		return
	}

	log.Printf("✅ Produto criado: %s (ID: %s)", novoProduto.Descricao, novoProduto.ID)
	c.JSON(http.StatusCreated, novoProduto)
}
func atualizarProduto(c *gin.Context) {
	id := c.Param("id")
	log.Printf("✏️ Atualizando produto ID: %s", id)

	var produto Produto
	if err := c.ShouldBindJSON(&produto); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}
	var exists int
	checkQuery := "SELECT COUNT(*) FROM produtos WHERE id = ?"
	err := db.QueryRow(checkQuery, id).Scan(&exists)
	if err != nil || exists == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Produto não encontrado",
			Message: fmt.Sprintf("Não existe produto com ID %s", id),
		})
		return
	}
	updateQuery := `
		UPDATE produtos 
		SET codigo = ?, descricao = ?, saldo = ?, imagem_url = ? 
		WHERE id = ?
	`

	_, err = db.Exec(updateQuery,
		produto.Codigo,
		produto.Descricao,
		produto.Saldo,
		produto.ImagemUrl,
		id,
	)

	if err != nil {
		log.Printf("❌ Erro ao atualizar produto: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao atualizar produto",
			Message: err.Error(),
		})
		return
	}
	selectQuery := `
		SELECT id, codigo, descricao, saldo, imagem_url, criado_em, atualizado_em 
		FROM produtos 
		WHERE id = ?
	`

	var produtoAtualizado Produto
	err = db.QueryRow(selectQuery, id).Scan(
		&produtoAtualizado.ID,
		&produtoAtualizado.Codigo,
		&produtoAtualizado.Descricao,
		&produtoAtualizado.Saldo,
		&produtoAtualizado.ImagemUrl,
		&produtoAtualizado.CriadoEm,
		&produtoAtualizado.AtualizadoEm,
	)

	if err != nil {
		log.Printf("❌ Erro ao buscar produto atualizado: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar produto atualizado",
			Message: err.Error(),
		})
		return
	}

	log.Printf("✅ Produto atualizado: %s", produtoAtualizado.Descricao)
	c.JSON(http.StatusOK, produtoAtualizado)
}
func removerProduto(c *gin.Context) {
	id := c.Param("id")
	log.Printf("🗑️ Removendo produto ID: %s", id)
	var descricao string
	selectQuery := "SELECT descricao FROM produtos WHERE id = ?"
	err := db.QueryRow(selectQuery, id).Scan(&descricao)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Produto não encontrado",
			Message: fmt.Sprintf("Não existe produto com ID %s", id),
		})
		return
	}

	if err != nil {
		log.Printf("❌ Erro ao buscar produto: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar produto",
			Message: err.Error(),
		})
		return
	}
	deleteQuery := "DELETE FROM produtos WHERE id = ?"
	_, err = db.Exec(deleteQuery, id)

	if err != nil {
		if err.Error() != "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Produto vinculado a notas fiscais",
				Message: fmt.Sprintf("O produto \"%s\" possui notas fiscais vinculadas e não pode ser excluído. Para excluir este produto, primeiro exclua as notas fiscais que o utilizam.", descricao),
			})
			return
		}

		log.Printf("❌ Erro ao remover produto: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao remover produto",
			Message: err.Error(),
		})
		return
	}

	log.Printf("✅ Produto removido: %s", descricao)
	c.JSON(http.StatusOK, gin.H{
		"message":   "Produto removido com sucesso",
		"descricao": descricao,
	})
}
func atualizarSaldo(c *gin.Context) {
	id := c.Param("id")
	log.Printf("📊 Atualizando saldo do produto ID: %s", id)

	var req AtualizarSaldoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: "Quantidade é obrigatória",
		})
		return
	}
	maxRetries := 3
	for tentativa := 1; tentativa <= maxRetries; tentativa++ {
		log.Printf("🔄 Tentativa %d de %d para atualizar saldo do produto %s", tentativa, maxRetries, id)
		tx, err := db.Begin()
		if err != nil {
			log.Printf("❌ Erro ao iniciar transação: %v", err)
			if tentativa == maxRetries {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Erro ao processar requisição",
					Message: "Não foi possível iniciar a transação após várias tentativas",
				})
				return
			}
			continue
		}

		var descricao string
		var saldoAtual int
		selectQuery := "SELECT descricao, saldo FROM produtos WHERE id = ? FOR UPDATE"
		err = tx.QueryRow(selectQuery, id).Scan(&descricao, &saldoAtual)

		if err == sql.ErrNoRows {
			tx.Rollback()
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Produto não encontrado",
				Message: fmt.Sprintf("Não existe produto com ID %s", id),
			})
			return
		}

		if err != nil {
			log.Printf("❌ Erro ao buscar produto (tentativa %d): %v", tentativa, err)
			tx.Rollback()
			if tentativa == maxRetries {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Erro ao buscar produto",
					Message: err.Error(),
				})
				return
			}
			continue
		}
		novoSaldo := saldoAtual - req.Quantidade
		if novoSaldo < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Saldo insuficiente",
				Message: fmt.Sprintf("O produto %s possui apenas %d unidades em estoque (tentando reservar %d)", descricao, saldoAtual, req.Quantidade),
			})
			return
		}
		updateQuery := "UPDATE produtos SET saldo = ? WHERE id = ? AND saldo = ?"
		result, err := tx.Exec(updateQuery, novoSaldo, id, saldoAtual)

		if err != nil {
			log.Printf("❌ Erro ao atualizar saldo (tentativa %d): %v", tentativa, err)
			tx.Rollback()
			if tentativa == maxRetries {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Erro ao atualizar saldo",
					Message: err.Error(),
				})
				return
			}
			continue
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("❌ Erro ao verificar linhas afetadas: %v", err)
			tx.Rollback()
			if tentativa == maxRetries {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Erro ao verificar atualização",
					Message: err.Error(),
				})
				return
			}
			continue
		}

		if rowsAffected == 0 {
			log.Printf("⚠️ Conflito de concorrência detectado na tentativa %d - saldo foi modificado por outra transação", tentativa)
			tx.Rollback()
			if tentativa == maxRetries {
				c.JSON(http.StatusConflict, ErrorResponse{
					Error:   "Conflito de concorrência",
					Message: "O saldo foi modificado por outra operação. Por favor, tente novamente.",
				})
				return
			}
			continue
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("❌ Erro ao confirmar transação (tentativa %d): %v", tentativa, err)
			if tentativa == maxRetries {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Erro ao confirmar transação",
					Message: err.Error(),
				})
				return
			}
			continue
		}
		var produtoAtualizado Produto
		finalQuery := `
			SELECT id, codigo, descricao, saldo, imagem_url, criado_em, atualizado_em 
			FROM produtos 
			WHERE id = ?
		`
		err = db.QueryRow(finalQuery, id).Scan(
			&produtoAtualizado.ID,
			&produtoAtualizado.Codigo,
			&produtoAtualizado.Descricao,
			&produtoAtualizado.Saldo,
			&produtoAtualizado.ImagemUrl,
			&produtoAtualizado.CriadoEm,
			&produtoAtualizado.AtualizadoEm,
		)

		if err != nil {
			log.Printf("❌ Erro ao buscar produto atualizado: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro ao buscar produto atualizado",
				Message: err.Error(),
			})
			return
		}

		log.Printf("✅ Saldo atualizado com sucesso! Produto: %s | Saldo anterior: %d | Quantidade reservada: %d | Novo saldo: %d | Tentativa: %d",
			descricao, saldoAtual, req.Quantidade, novoSaldo, tentativa)

		c.JSON(http.StatusOK, produtoAtualizado)
		return
	}
	log.Printf("❌ Falha ao atualizar saldo após %d tentativas", maxRetries)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "Erro ao processar requisição",
		Message: "Não foi possível completar a operação após várias tentativas",
	})
}
func healthCheck(c *gin.Context) {
	err := db.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    "ERROR",
			"service":   "Serviço de Estoque",
			"database":  "MariaDB desconectado",
			"error":     err.Error(),
			"timestamp": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"service":   "Serviço de Estoque",
		"database":  "MariaDB conectado",
		"timestamp": "",
	})
}
func main() {
	log.Println("\n🚀 Iniciando Serviço de Estoque (Golang)...\n")
	err := inicializarBanco()
	if err != nil {
		log.Fatal("❌ Erro ao conectar no banco de dados:", err)
	}
	defer db.Close()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	api := router.Group("/api")
	{
		api.GET("/produtos", listarProdutos)
		api.GET("/produtos/:id", buscarProduto)
		api.POST("/produtos", criarProduto)
		api.PUT("/produtos/:id", atualizarProduto)
		api.DELETE("/produtos/:id", removerProduto)
		api.PUT("/produtos/:id/atualizar-saldo", atualizarSaldo)
	}
	router.GET("/health", healthCheck)
	log.Println("\n✅ Serviço de Estoque rodando na porta 3001")
	log.Println("📍 URL: http://localhost:3001")
	log.Println("🗄️  Banco: MariaDB (notafiscal_desafio)")
	log.Println("\n📝 Rotas disponíveis:")
	log.Println("   GET    /api/produtos")
	log.Println("   GET    /api/produtos/:id")
	log.Println("   POST   /api/produtos")
	log.Println("   PUT    /api/produtos/:id")
	log.Println("   DELETE /api/produtos/:id")
	log.Println("   PUT    /api/produtos/:id/atualizar-saldo")
	log.Println("   GET    /health\n")

	if err := router.Run(":3001"); err != nil {
		log.Fatal("❌ Erro ao iniciar servidor:", err)
	}
}
