package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type CircuitState string

const (
	StateClosed   CircuitState = "CLOSED"    // Circuito fechado (funcionando)
	StateOpen     CircuitState = "OPEN"      // Circuito aberto (muitas falhas)
	StateHalfOpen CircuitState = "HALF_OPEN" // Circuito meio aberto (testando)
)

type CircuitBreaker struct {
	State        CircuitState  // Estado atual
	Failures     int           // Número de falhas consecutivas
	LastFailTime time.Time     // Horário da última falha
	Threshold    int           // Limite de falhas para abrir o circuito
	Timeout      time.Duration // Tempo para tentar fechar o circuito
}
type NotaFiscal struct {
	ID           string     `json:"id"`                     // UUID único
	Numero       int        `json:"numero"`                 // Número sequencial da nota
	Status       string     `json:"status"`                 // ABERTA ou FECHADA
	DataEmissao  *time.Time `json:"dataEmissao,omitempty"`  // Data de finalização
	CriadoEm     string     `json:"criadoEm,omitempty"`     // Data de criação
	AtualizadoEm string     `json:"atualizadoEm,omitempty"` // Data de atualização
	Itens        []NotaItem `json:"itens,omitempty"`        // Itens da nota (opcional)
}
type NotaItem struct {
	ID               string `json:"id"`                            // UUID único
	NotaId           string `json:"notaId"`                        // FK para nota fiscal
	ProdutoId        string `json:"produtoId" binding:"required"`  // FK para produto
	ProdutoCodigo    string `json:"produtoCodigo"`                 // Código do produto
	ProdutoDescricao string `json:"produtoDescricao"`              // Descrição do produto
	Quantidade       int    `json:"quantidade" binding:"required"` // Quantidade vendida
	CriadoEm         string `json:"criadoEm,omitempty"`            // Data de criação
}
type CriarNotaRequest struct {
	Itens []NotaItem `json:"itens" binding:"required,dive"`
}
type AtualizarNotaRequest struct {
	Itens []NotaItem `json:"itens" binding:"required,dive"`
}
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
type Produto struct {
	ID        string `json:"id"`
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
}

var (
	db                *sql.DB
	circuitBreaker    *CircuitBreaker
	estoqueServiceURL = "http://localhost:3001" // URL do serviço de estoque
)

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
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	log.Println("✅ Conexão com MariaDB estabelecida com sucesso!")
	return nil
}
func inicializarCircuitBreaker() {
	circuitBreaker = &CircuitBreaker{
		State:     StateClosed,
		Failures:  0,
		Threshold: 3,                // Abre após 3 falhas consecutivas
		Timeout:   10 * time.Second, // Tenta fechar após 10 segundos
	}
	log.Println("✅ Circuit Breaker inicializado (Threshold: 3, Timeout: 10s)")
}
func (cb *CircuitBreaker) recordSuccess() {
	cb.Failures = 0
	if cb.State == StateHalfOpen {
		cb.State = StateClosed
		log.Println("🔵 Circuit Breaker: Estado alterado para CLOSED")
	}
}
func (cb *CircuitBreaker) recordFailure() {
	cb.Failures++
	cb.LastFailTime = time.Now()

	if cb.Failures >= cb.Threshold {
		cb.State = StateOpen
		log.Printf("🔴 Circuit Breaker: Estado alterado para OPEN (%d falhas)", cb.Failures)
	}
}
func (cb *CircuitBreaker) canCall() bool {
	if cb.State == StateClosed {
		return true
	}

	if cb.State == StateOpen {
		if time.Since(cb.LastFailTime) >= cb.Timeout {
			cb.State = StateHalfOpen
			log.Println("🟡 Circuit Breaker: Estado alterado para HALF_OPEN (testando)")
			return true
		}
		return false
	}
	return true
}
func buscarProdutoEstoque(produtoId string) (*Produto, error) {
	if !circuitBreaker.canCall() {
		return nil, fmt.Errorf("circuit breaker está OPEN - serviço de estoque indisponível")
	}
	maxRetries := 3
	retryDelays := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}

	var lastErr error

	for tentativa := 0; tentativa < maxRetries; tentativa++ {
		if tentativa > 0 {
			delay := retryDelays[tentativa-1]
			log.Printf("⏳ Retry %d/%d após %v", tentativa, maxRetries-1, delay)
			time.Sleep(delay)
		}
		url := fmt.Sprintf("%s/api/produtos/%s", estoqueServiceURL, produtoId)
		resp, err := http.Get(url)

		if err != nil {
			lastErr = fmt.Errorf("erro na requisição: %v", err)
			log.Printf("❌ Tentativa %d falhou: %v", tentativa+1, err)
			continue
		}

		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var produto Produto
			err = json.NewDecoder(resp.Body).Decode(&produto)
			if err != nil {
				lastErr = fmt.Errorf("erro ao decodificar resposta: %v", err)
				continue
			}
			circuitBreaker.recordSuccess()
			return &produto, nil
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("produto não encontrado")
		}
		body, _ := io.ReadAll(resp.Body)
		lastErr = fmt.Errorf("erro HTTP %d: %s", resp.StatusCode, string(body))
		log.Printf("❌ Tentativa %d falhou: Status %d", tentativa+1, resp.StatusCode)
	}
	circuitBreaker.recordFailure()
	return nil, fmt.Errorf("todas as tentativas falharam: %v", lastErr)
}
func atualizarSaldoEstoque(produtoId string, quantidade int) error {
	if !circuitBreaker.canCall() {
		return fmt.Errorf("circuit breaker está OPEN - serviço de estoque indisponível")
	}
	maxRetries := 3
	retryDelays := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}

	var lastErr error

	for tentativa := 0; tentativa < maxRetries; tentativa++ {
		if tentativa > 0 {
			delay := retryDelays[tentativa-1]
			log.Printf("⏳ Retry %d/%d após %v", tentativa, maxRetries-1, delay)
			time.Sleep(delay)
		}
		payload := map[string]int{"quantidade": quantidade}
		jsonData, _ := json.Marshal(payload)
		url := fmt.Sprintf("%s/api/produtos/%s/atualizar-saldo", estoqueServiceURL, produtoId)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("erro ao criar requisição: %v", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)

		if err != nil {
			lastErr = fmt.Errorf("erro na requisição: %v", err)
			log.Printf("❌ Tentativa %d falhou: %v", tentativa+1, err)
			continue
		}

		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			circuitBreaker.recordSuccess()
			return nil
		}
		body, _ := io.ReadAll(resp.Body)
		lastErr = fmt.Errorf("erro HTTP %d: %s", resp.StatusCode, string(body))
		log.Printf("❌ Tentativa %d falhou: Status %d", tentativa+1, resp.StatusCode)
		if resp.StatusCode == http.StatusBadRequest {
			return lastErr
		}
	}
	circuitBreaker.recordFailure()
	return fmt.Errorf("todas as tentativas falharam: %v", lastErr)
}
func listarNotas(c *gin.Context) {
	log.Println("📄 Listando todas as notas fiscais")

	query := `
		SELECT id, numero, status, data_emissao, criado_em, atualizado_em 
		FROM notas_fiscais 
		ORDER BY numero DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("❌ Erro ao buscar notas: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar notas",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	notas := []NotaFiscal{}

	for rows.Next() {
		var n NotaFiscal
		err := rows.Scan(
			&n.ID,
			&n.Numero,
			&n.Status,
			&n.DataEmissao,
			&n.CriadoEm,
			&n.AtualizadoEm,
		)
		if err != nil {
			log.Printf("❌ Erro ao ler nota: %v", err)
			continue
		}
		itensQuery := `
			SELECT id, nota_id, produto_id, produto_codigo, produto_descricao, quantidade, criado_em 
			FROM notas_itens 
			WHERE nota_id = ?
		`

		itensRows, err := db.Query(itensQuery, n.ID)
		if err != nil {
			log.Printf("❌ Erro ao buscar itens da nota %s: %v", n.ID, err)
			n.Itens = []NotaItem{}
		} else {
			n.Itens = []NotaItem{}
			for itensRows.Next() {
				var item NotaItem
				err := itensRows.Scan(
					&item.ID,
					&item.NotaId,
					&item.ProdutoId,
					&item.ProdutoCodigo,
					&item.ProdutoDescricao,
					&item.Quantidade,
					&item.CriadoEm,
				)
				if err != nil {
					log.Printf("❌ Erro ao ler item: %v", err)
					continue
				}
				n.Itens = append(n.Itens, item)
			}
			itensRows.Close()
		}

		notas = append(notas, n)
	}

	log.Printf("✅ %d notas encontradas", len(notas))
	c.JSON(http.StatusOK, notas)
}
func buscarNota(c *gin.Context) {
	id := c.Param("id")
	log.Printf("📄 Buscando nota com ID: %s", id)
	notaQuery := `
		SELECT id, numero, status, data_emissao, criado_em, atualizado_em 
		FROM notas_fiscais 
		WHERE id = ?
	`

	var nota NotaFiscal
	err := db.QueryRow(notaQuery, id).Scan(
		&nota.ID,
		&nota.Numero,
		&nota.Status,
		&nota.DataEmissao,
		&nota.CriadoEm,
		&nota.AtualizadoEm,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Nota não encontrada",
			Message: fmt.Sprintf("Não existe nota com ID %s", id),
		})
		return
	}

	if err != nil {
		log.Printf("❌ Erro ao buscar nota: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar nota",
			Message: err.Error(),
		})
		return
	}
	itensQuery := `
		SELECT id, nota_id, produto_id, produto_codigo, produto_descricao, quantidade, criado_em 
		FROM notas_itens 
		WHERE nota_id = ?
	`

	rows, err := db.Query(itensQuery, id)
	if err != nil {
		log.Printf("❌ Erro ao buscar itens: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar itens",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	nota.Itens = []NotaItem{}

	for rows.Next() {
		var item NotaItem
		err := rows.Scan(
			&item.ID,
			&item.NotaId,
			&item.ProdutoId,
			&item.ProdutoCodigo,
			&item.ProdutoDescricao,
			&item.Quantidade,
			&item.CriadoEm,
		)
		if err != nil {
			log.Printf("❌ Erro ao ler item: %v", err)
			continue
		}
		nota.Itens = append(nota.Itens, item)
	}

	log.Printf("✅ Nota encontrada: #%d (%s) com %d itens", nota.Numero, nota.Status, len(nota.Itens))
	c.JSON(http.StatusOK, nota)
}
func criarNota(c *gin.Context) {
	log.Println("➕ Criando nova nota fiscal")

	var req CriarNotaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: "É necessário informar pelo menos um item",
		})
		return
	}

	if len(req.Itens) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: "A nota deve ter pelo menos um item",
		})
		return
	}
	for i := range req.Itens {
		produto, err := buscarProdutoEstoque(req.Itens[i].ProdutoId)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Produto não encontrado",
				Message: fmt.Sprintf("Erro ao buscar produto %s: %v", req.Itens[i].ProdutoId, err),
			})
			return
		}
		req.Itens[i].ProdutoCodigo = produto.Codigo
		req.Itens[i].ProdutoDescricao = produto.Descricao
	}
	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Erro ao iniciar transação: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao criar nota",
			Message: err.Error(),
		})
		return
	}
	notaId := uuid.New().String()
	insertNotaQuery := `
		INSERT INTO notas_fiscais (id, status) 
		VALUES (?, 'ABERTA')
	`

	_, err = tx.Exec(insertNotaQuery, notaId)
	if err != nil {
		tx.Rollback()
		log.Printf("❌ Erro ao criar nota: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao criar nota",
			Message: err.Error(),
		})
		return
	}
	insertItemQuery := `
		INSERT INTO notas_itens (nota_id, produto_id, produto_codigo, produto_descricao, quantidade) 
		VALUES (?, ?, ?, ?, ?)
	`

	for _, item := range req.Itens {
		_, err = tx.Exec(insertItemQuery,
			notaId,
			item.ProdutoId,
			item.ProdutoCodigo,
			item.ProdutoDescricao,
			item.Quantidade,
		)

		if err != nil {
			tx.Rollback()
			log.Printf("❌ Erro ao inserir item: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro ao inserir item",
				Message: err.Error(),
			})
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("❌ Erro ao confirmar transação: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao confirmar transação",
			Message: err.Error(),
		})
		return
	}
	var notaCriada NotaFiscal
	selectQuery := `
		SELECT id, numero, status, data_emissao, criado_em, atualizado_em 
		FROM notas_fiscais 
		WHERE id = ?
	`
	err = db.QueryRow(selectQuery, notaId).Scan(
		&notaCriada.ID,
		&notaCriada.Numero,
		&notaCriada.Status,
		&notaCriada.DataEmissao,
		&notaCriada.CriadoEm,
		&notaCriada.AtualizadoEm,
	)

	if err != nil {
		log.Printf("❌ Erro ao buscar nota criada: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar nota criada",
			Message: err.Error(),
		})
		return
	}
	itensQuery := `
		SELECT id, nota_id, produto_id, produto_codigo, produto_descricao, quantidade, criado_em 
		FROM notas_itens 
		WHERE nota_id = ?
	`

	rows, err := db.Query(itensQuery, notaId)
	if err == nil {
		defer rows.Close()
		notaCriada.Itens = []NotaItem{}

		for rows.Next() {
			var item NotaItem
			err := rows.Scan(
				&item.ID,
				&item.NotaId,
				&item.ProdutoId,
				&item.ProdutoCodigo,
				&item.ProdutoDescricao,
				&item.Quantidade,
				&item.CriadoEm,
			)
			if err == nil {
				notaCriada.Itens = append(notaCriada.Itens, item)
			}
		}
	}

	log.Printf("✅ Nota criada: #%d com %d itens", notaCriada.Numero, len(notaCriada.Itens))
	c.JSON(http.StatusCreated, notaCriada)
}
func atualizarNota(c *gin.Context) {
	id := c.Param("id")
	log.Printf("✏️ Atualizando nota ID: %s", id)

	var req AtualizarNotaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: "É necessário informar pelo menos um item",
		})
		return
	}

	if len(req.Itens) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: "A nota deve ter pelo menos um item",
		})
		return
	}
	var status string
	checkQuery := "SELECT status FROM notas_fiscais WHERE id = ?"
	err := db.QueryRow(checkQuery, id).Scan(&status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Nota não encontrada",
			Message: fmt.Sprintf("Não existe nota com ID %s", id),
		})
		return
	}

	if err != nil {
		log.Printf("❌ Erro ao buscar nota: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar nota",
			Message: err.Error(),
		})
		return
	}

	if status != "ABERTA" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Nota já finalizada",
			Message: "Não é possível alterar uma nota com status FECHADA",
		})
		return
	}
	for i := range req.Itens {
		produto, err := buscarProdutoEstoque(req.Itens[i].ProdutoId)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Produto não encontrado",
				Message: fmt.Sprintf("Erro ao buscar produto %s: %v", req.Itens[i].ProdutoId, err),
			})
			return
		}

		req.Itens[i].ProdutoCodigo = produto.Codigo
		req.Itens[i].ProdutoDescricao = produto.Descricao
	}
	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Erro ao iniciar transação: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao atualizar nota",
			Message: err.Error(),
		})
		return
	}
	deleteItensQuery := "DELETE FROM notas_itens WHERE nota_id = ?"
	_, err = tx.Exec(deleteItensQuery, id)
	if err != nil {
		tx.Rollback()
		log.Printf("❌ Erro ao remover itens antigos: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao remover itens antigos",
			Message: err.Error(),
		})
		return
	}
	insertItemQuery := `
		INSERT INTO notas_itens (nota_id, produto_id, produto_codigo, produto_descricao, quantidade) 
		VALUES (?, ?, ?, ?, ?)
	`

	for _, item := range req.Itens {
		_, err = tx.Exec(insertItemQuery,
			id,
			item.ProdutoId,
			item.ProdutoCodigo,
			item.ProdutoDescricao,
			item.Quantidade,
		)

		if err != nil {
			tx.Rollback()
			log.Printf("❌ Erro ao inserir item: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Erro ao inserir item",
				Message: err.Error(),
			})
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("❌ Erro ao confirmar transação: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao confirmar transação",
			Message: err.Error(),
		})
		return
	}
	var notaAtualizada NotaFiscal
	selectQuery := `
		SELECT id, numero, status, data_emissao, criado_em, atualizado_em 
		FROM notas_fiscais 
		WHERE id = ?
	`
	err = db.QueryRow(selectQuery, id).Scan(
		&notaAtualizada.ID,
		&notaAtualizada.Numero,
		&notaAtualizada.Status,
		&notaAtualizada.DataEmissao,
		&notaAtualizada.CriadoEm,
		&notaAtualizada.AtualizadoEm,
	)

	if err != nil {
		log.Printf("❌ Erro ao buscar nota atualizada: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar nota atualizada",
			Message: err.Error(),
		})
		return
	}
	itensQuery := `
		SELECT id, nota_id, produto_id, produto_codigo, produto_descricao, quantidade, criado_em 
		FROM notas_itens 
		WHERE nota_id = ?
	`

	rows, err := db.Query(itensQuery, id)
	if err == nil {
		defer rows.Close()
		notaAtualizada.Itens = []NotaItem{}

		for rows.Next() {
			var item NotaItem
			err := rows.Scan(
				&item.ID,
				&item.NotaId,
				&item.ProdutoId,
				&item.ProdutoCodigo,
				&item.ProdutoDescricao,
				&item.Quantidade,
				&item.CriadoEm,
			)
			if err == nil {
				notaAtualizada.Itens = append(notaAtualizada.Itens, item)
			}
		}
	}

	log.Printf("✅ Nota atualizada: #%d com %d itens", notaAtualizada.Numero, len(notaAtualizada.Itens))
	c.JSON(http.StatusOK, notaAtualizada)
}
func removerNota(c *gin.Context) {
	id := c.Param("id")
	log.Printf("🗑️ Removendo nota ID: %s", id)
	var numero int
	var status string
	checkQuery := "SELECT numero, status FROM notas_fiscais WHERE id = ?"
	err := db.QueryRow(checkQuery, id).Scan(&numero, &status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Nota não encontrada",
			Message: fmt.Sprintf("Não existe nota com ID %s", id),
		})
		return
	}

	if err != nil {
		log.Printf("❌ Erro ao buscar nota: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar nota",
			Message: err.Error(),
		})
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Erro ao iniciar transação: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao remover nota",
			Message: err.Error(),
		})
		return
	}
	deleteItensQuery := "DELETE FROM notas_itens WHERE nota_id = ?"
	_, err = tx.Exec(deleteItensQuery, id)
	if err != nil {
		tx.Rollback()
		log.Printf("❌ Erro ao remover itens: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao remover itens",
			Message: err.Error(),
		})
		return
	}
	deleteNotaQuery := "DELETE FROM notas_fiscais WHERE id = ?"
	_, err = tx.Exec(deleteNotaQuery, id)
	if err != nil {
		tx.Rollback()
		log.Printf("❌ Erro ao remover nota: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao remover nota",
			Message: err.Error(),
		})
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("❌ Erro ao confirmar transação: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao confirmar transação",
			Message: err.Error(),
		})
		return
	}

	log.Printf("✅ Nota removida: #%d", numero)
	c.JSON(http.StatusOK, gin.H{
		"message": "Nota fiscal removida com sucesso",
		"numero":  numero,
	})
}
func finalizarNota(c *gin.Context) {
	id := c.Param("id")
	log.Printf("🔒 Finalizando nota ID: %s", id)
	var numero int
	var status string
	checkQuery := "SELECT numero, status FROM notas_fiscais WHERE id = ?"
	err := db.QueryRow(checkQuery, id).Scan(&numero, &status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Nota não encontrada",
			Message: fmt.Sprintf("Não existe nota com ID %s", id),
		})
		return
	}

	if err != nil {
		log.Printf("❌ Erro ao buscar nota: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar nota",
			Message: err.Error(),
		})
		return
	}

	if status != "ABERTA" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Nota já finalizada",
			Message: fmt.Sprintf("A nota #%d já está finalizada", numero),
		})
		return
	}
	itensQuery := "SELECT produto_id, quantidade FROM notas_itens WHERE nota_id = ?"
	rows, err := db.Query(itensQuery, id)
	if err != nil {
		log.Printf("❌ Erro ao buscar itens: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar itens",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	type ItemParaAtualizar struct {
		ProdutoId  string
		Quantidade int
	}

	itensParaAtualizar := []ItemParaAtualizar{}

	for rows.Next() {
		var item ItemParaAtualizar
		err := rows.Scan(&item.ProdutoId, &item.Quantidade)
		if err != nil {
			log.Printf("❌ Erro ao ler item: %v", err)
			continue
		}
		itensParaAtualizar = append(itensParaAtualizar, item)
	}

	if len(itensParaAtualizar) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Nota sem itens",
			Message: "Não é possível finalizar uma nota sem itens",
		})
		return
	}
	log.Printf("📦 Atualizando saldo de %d produtos...", len(itensParaAtualizar))
	for _, item := range itensParaAtualizar {
		err := atualizarSaldoEstoque(item.ProdutoId, item.Quantidade)
		if err != nil {
			log.Printf("❌ Erro ao atualizar saldo: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Erro ao atualizar estoque",
				Message: fmt.Sprintf("Erro ao atualizar saldo do produto: %v", err),
			})
			return
		}
	}
	updateQuery := `
		UPDATE notas_fiscais 
		SET status = 'FECHADA', data_emissao = NOW() 
		WHERE id = ?
	`

	_, err = db.Exec(updateQuery, id)
	if err != nil {
		log.Printf("❌ Erro ao finalizar nota: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao finalizar nota",
			Message: err.Error(),
		})
		return
	}
	var notaFinalizada NotaFiscal
	selectQuery := `
		SELECT id, numero, status, data_emissao, criado_em, atualizado_em 
		FROM notas_fiscais 
		WHERE id = ?
	`
	err = db.QueryRow(selectQuery, id).Scan(
		&notaFinalizada.ID,
		&notaFinalizada.Numero,
		&notaFinalizada.Status,
		&notaFinalizada.DataEmissao,
		&notaFinalizada.CriadoEm,
		&notaFinalizada.AtualizadoEm,
	)

	if err != nil {
		log.Printf("❌ Erro ao buscar nota finalizada: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao buscar nota finalizada",
			Message: err.Error(),
		})
		return
	}
	itensFinalizadosQuery := `
		SELECT id, nota_id, produto_id, produto_codigo, produto_descricao, quantidade, criado_em 
		FROM notas_itens 
		WHERE nota_id = ?
	`

	rowsFinalizados, err := db.Query(itensFinalizadosQuery, id)
	if err == nil {
		defer rowsFinalizados.Close()
		notaFinalizada.Itens = []NotaItem{}

		for rowsFinalizados.Next() {
			var item NotaItem
			err := rowsFinalizados.Scan(
				&item.ID,
				&item.NotaId,
				&item.ProdutoId,
				&item.ProdutoCodigo,
				&item.ProdutoDescricao,
				&item.Quantidade,
				&item.CriadoEm,
			)
			if err == nil {
				notaFinalizada.Itens = append(notaFinalizada.Itens, item)
			}
		}
	}

	log.Printf("✅ Nota finalizada: #%d com %d itens", notaFinalizada.Numero, len(notaFinalizada.Itens))
	c.JSON(http.StatusOK, notaFinalizada)
}
func healthCheck(c *gin.Context) {
	err := db.Ping()
	dbStatus := "MariaDB conectado"
	if err != nil {
		dbStatus = "MariaDB desconectado: " + err.Error()
	}
	circuitStatus := fmt.Sprintf("Circuit Breaker: %s (Falhas: %d)", circuitBreaker.State, circuitBreaker.Failures)

	c.JSON(http.StatusOK, gin.H{
		"status":         "OK",
		"service":        "Serviço de Faturamento",
		"database":       dbStatus,
		"circuitBreaker": circuitStatus,
		"estoqueService": estoqueServiceURL,
		"timestamp":      time.Now().Format("2006-01-02 15:04:05"),
	})
}
func circuitBreakerStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"state":        circuitBreaker.State,
		"failures":     circuitBreaker.Failures,
		"threshold":    circuitBreaker.Threshold,
		"timeout":      circuitBreaker.Timeout.String(),
		"lastFailTime": circuitBreaker.LastFailTime.Format("2006-01-02 15:04:05"),
	})
}
func resetCircuitBreaker(c *gin.Context) {
	circuitBreaker.State = StateClosed
	circuitBreaker.Failures = 0
	log.Println("🔄 Circuit Breaker resetado manualmente")

	c.JSON(http.StatusOK, gin.H{
		"message": "Circuit Breaker resetado com sucesso",
		"state":   circuitBreaker.State,
	})
}

var (
	huggingfaceAPIKey string
	huggingfaceModel  string
)

type ChatIARequest struct {
	Pergunta string `json:"pergunta" binding:"required"`
}
type IAResponse struct {
	Resposta string `json:"resposta"`
}
type HFChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type HFChatRequest struct {
	Model       string          `json:"model"`
	Messages    []HFChatMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type HFChatChoice struct {
	Message HFChatMessage `json:"message"`
}

type HFChatResponse struct {
	Choices []HFChatChoice `json:"choices"`
}

func chamarHuggingFaceChat(mensagens []HFChatMessage) (string, error) {
	apiURL := "https://router.huggingface.co/v1/chat/completions"

	requestBody := HFChatRequest{
		Model:       huggingfaceModel,
		Messages:    mensagens,
		Temperature: 0.7,
		MaxTokens:   500,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao criar JSON: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+huggingfaceAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erro HTTP %d: %s", resp.StatusCode, string(body))
	}

	var chatResp HFChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("nenhuma resposta gerada")
}
func chatIA(c *gin.Context) {
	var req ChatIARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Dados inválidos",
			Message: "Pergunta é obrigatória",
		})
		return
	}
	var totalProdutos, totalNotas, notasAbertas, notasFechadas int
	db.QueryRow("SELECT COUNT(*) FROM produtos").Scan(&totalProdutos)
	db.QueryRow("SELECT COUNT(*) FROM notas_fiscais").Scan(&totalNotas)
	db.QueryRow("SELECT COUNT(*) FROM notas_fiscais WHERE status = 'ABERTA'").Scan(&notasAbertas)
	db.QueryRow("SELECT COUNT(*) FROM notas_fiscais WHERE status = 'FECHADA'").Scan(&notasFechadas)
	rows, _ := db.Query(`
		SELECT codigo, descricao, saldo 
		FROM produtos 
		WHERE saldo < 10 
		ORDER BY saldo ASC 
		LIMIT 5
	`)
	defer rows.Close()

	var produtosBaixoEstoque []string
	for rows.Next() {
		var codigo, descricao string
		var saldo int
		rows.Scan(&codigo, &descricao, &saldo)
		produtosBaixoEstoque = append(produtosBaixoEstoque, fmt.Sprintf("%s (%s): %d unidades", descricao, codigo, saldo))
	}
	sistemaMensagem := fmt.Sprintf(`Você é um assistente virtual de um sistema de gestão de notas fiscais e estoque.

CONTEXTO ATUAL:
- Produtos: %d
- Notas Fiscais: %d (Abertas: %d, Fechadas: %d)
- Produtos com estoque baixo: %s

Responda de forma clara e profissional.`,
		totalProdutos, totalNotas, notasAbertas, notasFechadas,
		strings.Join(produtosBaixoEstoque, ", "))

	mensagens := []HFChatMessage{
		{Role: "system", Content: sistemaMensagem},
		{Role: "user", Content: req.Pergunta},
	}

	resposta, err := chamarHuggingFaceChat(mensagens)
	if err != nil {
		log.Printf("❌ Erro ao chamar Hugging Face: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao processar pergunta",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, IAResponse{
		Resposta: resposta,
	})
}
func analisarIA(c *gin.Context) {
	var totalProdutos, totalNotas int
	db.QueryRow("SELECT COUNT(*) FROM produtos").Scan(&totalProdutos)
	db.QueryRow("SELECT COUNT(*) FROM notas_fiscais WHERE status = 'FECHADA'").Scan(&totalNotas)
	rows, _ := db.Query(`
		SELECT p.codigo, p.descricao, COALESCE(SUM(ni.quantidade), 0) as total_vendido
		FROM produtos p
		LEFT JOIN notas_itens ni ON p.id = ni.produto_id
		LEFT JOIN notas_fiscais nf ON ni.nota_id = nf.id AND nf.status = 'FECHADA'
		GROUP BY p.id, p.codigo, p.descricao
		ORDER BY total_vendido DESC
		LIMIT 5
	`)
	defer rows.Close()

	var maisVendidos []string
	for rows.Next() {
		var codigo, descricao string
		var totalVendido int
		rows.Scan(&codigo, &descricao, &totalVendido)
		maisVendidos = append(maisVendidos, fmt.Sprintf("%s: %d vendas", descricao, totalVendido))
	}
	rows2, _ := db.Query(`
		SELECT codigo, descricao, saldo 
		FROM produtos 
		WHERE saldo < 10 
		ORDER BY saldo ASC 
		LIMIT 5
	`)
	defer rows2.Close()

	var estoqueBaixo []string
	for rows2.Next() {
		var codigo, descricao string
		var saldo int
		rows2.Scan(&codigo, &descricao, &saldo)
		estoqueBaixo = append(estoqueBaixo, fmt.Sprintf("%s: %d unidades", descricao, saldo))
	}
	prompt := fmt.Sprintf(`Analise os seguintes dados de um sistema de gestão:

DADOS:
- Total de Produtos: %d
- Total de Notas Fiscais: %d

TOP 5 MAIS VENDIDOS:
%s

ESTOQUE BAIXO:
%s

Faça uma análise profissional com:
1. Resumo executivo
2. Insights principais  
3. Recomendações estratégicas`,
		totalProdutos, totalNotas,
		strings.Join(maisVendidos, "\n"),
		strings.Join(estoqueBaixo, "\n"))

	mensagens := []HFChatMessage{
		{Role: "system", Content: "Você é um analista de dados especializado em gestão de estoque e vendas."},
		{Role: "user", Content: prompt},
	}

	resposta, err := chamarHuggingFaceChat(mensagens)
	if err != nil {
		log.Printf("❌ Erro ao chamar Hugging Face: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro ao analisar dados",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, IAResponse{
		Resposta: resposta,
	})
}
func inicializarHuggingFace() {
	huggingfaceAPIKey = os.Getenv("HUGGINGFACE_API_KEY")
	if huggingfaceAPIKey == "" {
		log.Fatal("❌ HUGGINGFACE_API_KEY não configurada. Configure a variável de ambiente.")
	}

	huggingfaceModel = "Qwen/Qwen2.5-72B-Instruct"
	log.Println("🤖 Hugging Face AI configurado:", huggingfaceModel)
}
func main() {
	log.Println("\n🚀 Iniciando Serviço de Faturamento (Golang)...\n")
	err := inicializarBanco()
	if err != nil {
		log.Fatal("❌ Erro ao conectar no banco de dados:", err)
	}
	defer db.Close()
	inicializarCircuitBreaker()
	inicializarHuggingFace()
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
		api.GET("/notas", listarNotas)
		api.GET("/notas/:id", buscarNota)
		api.POST("/notas", criarNota)
		api.PUT("/notas/:id", atualizarNota)
		api.DELETE("/notas/:id", removerNota)
		api.POST("/notas/:id/finalizar", finalizarNota)
		api.POST("/ia/chat", chatIA)
		api.POST("/ia/analisar", analisarIA)
	}
	router.GET("/health", healthCheck)
	router.GET("/circuit-breaker", circuitBreakerStatus)
	router.POST("/circuit-breaker/reset", resetCircuitBreaker)
	log.Println("\n✅ Serviço de Faturamento rodando na porta 3002")
	log.Println("📍 URL: http://localhost:3002")
	log.Println("🗄️  Banco: MariaDB (notafiscal_desafio)")
	log.Println("🔌 Conectado ao Serviço de Estoque:", estoqueServiceURL)
	log.Println("\n📝 Rotas disponíveis:")
	log.Println("   GET    /api/notas")
	log.Println("   GET    /api/notas/:id")
	log.Println("   POST   /api/notas")
	log.Println("   PUT    /api/notas/:id")
	log.Println("   DELETE /api/notas/:id")
	log.Println("   POST   /api/notas/:id/finalizar")
	log.Println("   POST   /api/ia/chat              🤖 Chat com IA")
	log.Println("   POST   /api/ia/analisar          📊 Análise de dados")
	log.Println("   GET    /health")
	log.Println("   GET    /circuit-breaker")
	log.Println("   POST   /circuit-breaker/reset\n")

	if err := router.Run(":3002"); err != nil {
		log.Fatal("❌ Erro ao iniciar servidor:", err)
	}
}
