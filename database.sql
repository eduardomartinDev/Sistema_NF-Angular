-- ========================================
-- SCRIPT DE CRIAÇÃO DO BANCO DE DADOS KORP
-- ========================================
-- 
-- Este script cria:
-- 1. Banco de dados 'korp_sistema'
-- 2. Tabela de produtos (com imagem)
-- 3. Tabela de notas fiscais
-- 4. Tabela de itens das notas fiscais
-- 5. Dados iniciais para teste
--
-- COMO EXECUTAR NO HEIDISQL:
-- 1. Abra o HeidiSQL
-- 2. Conecte ao seu servidor MariaDB
-- 3. Vá em "Arquivo" -> "Executar arquivo SQL"
-- 4. Selecione este arquivo (database.sql)
-- 5. Clique em "Executar"
-- ========================================

-- Remove banco se já existir (CUIDADO: apaga todos os dados!)
DROP DATABASE IF EXISTS notafiscal_desafio;

-- Cria o banco de dados
CREATE DATABASE notafiscal_desafio 
  CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;

-- Seleciona o banco para uso
USE notafiscal_desafio;

-- ========================================
-- TABELA: produtos
-- Armazena os produtos do estoque
-- ========================================
CREATE TABLE produtos (
  -- Chave primária (UUID gerado pela aplicação)
  id VARCHAR(36) PRIMARY KEY,
  
  -- Código único do produto (ex: PROD001)
  codigo VARCHAR(50) NOT NULL UNIQUE,
  
  -- Descrição do produto
  descricao VARCHAR(255) NOT NULL,
  
  -- Quantidade disponível em estoque
  saldo INT NOT NULL DEFAULT 0,
  
  -- Imagem do produto em base64 (armazena a string completa: "data:image/png;base64,...")
  -- LONGTEXT suporta até 4GB de texto
  imagem_url LONGTEXT NULL,
  
  -- Data de criação do registro
  criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  -- Data da última atualização
  atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  -- Índices para melhorar performance nas consultas
  INDEX idx_codigo (codigo),
  INDEX idx_descricao (descricao)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ========================================
-- TABELA: notas_fiscais
-- Armazena as notas fiscais emitidas
-- ========================================
CREATE TABLE notas_fiscais (
  -- Chave primária (UUID gerado pela aplicação)
  id VARCHAR(36) PRIMARY KEY,
  
  -- Número sequencial da nota fiscal
  numero INT NOT NULL UNIQUE AUTO_INCREMENT,
  
  -- Status da nota: ABERTA ou FECHADA
  status ENUM('ABERTA', 'FECHADA') NOT NULL DEFAULT 'ABERTA',
  
  -- Data e hora de emissão
  data_emissao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  -- Data de criação do registro
  criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  -- Data da última atualização
  atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  -- Índices
  INDEX idx_numero (numero),
  INDEX idx_status (status),
  INDEX idx_data_emissao (data_emissao)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ========================================
-- TABELA: notas_itens
-- Armazena os itens (produtos) de cada nota fiscal
-- ========================================
CREATE TABLE notas_itens (
  -- Chave primária
  id INT PRIMARY KEY AUTO_INCREMENT,
  
  -- Referência para a nota fiscal (chave estrangeira)
  nota_id VARCHAR(36) NOT NULL,
  
  -- Referência para o produto (chave estrangeira)
  produto_id VARCHAR(36) NOT NULL,
  
  -- Código do produto (denormalizado para histórico)
  produto_codigo VARCHAR(50) NOT NULL,
  
  -- Descrição do produto (denormalizado para histórico)
  produto_descricao VARCHAR(255) NOT NULL,
  
  -- Quantidade do produto nesta nota
  quantidade INT NOT NULL,
  
  -- Data de criação do registro
  criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  -- Chaves estrangeiras com CASCADE
  -- Se deletar uma nota, deleta todos os itens automaticamente
  CONSTRAINT fk_nota 
    FOREIGN KEY (nota_id) 
    REFERENCES notas_fiscais(id) 
    ON DELETE CASCADE,
  
  -- Se deletar um produto, impede a exclusão se houver notas vinculadas
  CONSTRAINT fk_produto 
    FOREIGN KEY (produto_id) 
    REFERENCES produtos(id) 
    ON DELETE RESTRICT,
  
  -- Índices
  INDEX idx_nota_id (nota_id),
  INDEX idx_produto_id (produto_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ========================================
-- DADOS INICIAIS PARA TESTE
-- ========================================

-- Insere produtos de exemplo (sem imagem por enquanto)
INSERT INTO produtos (id, codigo, descricao, saldo) VALUES
  ('1', 'PROD001', 'Notebook Dell Inspiron 15', 10),
  ('2', 'PROD002', 'Mouse Logitech MX Master 3', 50),
  ('3', 'PROD003', 'Teclado Mecânico Keychron K2', 30),
  ('4', 'PROD004', 'Monitor LG 27 polegadas', 15),
  ('5', 'PROD005', 'Webcam Logitech C920', 25);

-- Insere uma nota fiscal de exemplo (FECHADA)
INSERT INTO notas_fiscais (id, numero, status, data_emissao) VALUES
  ('nota-1', 1, 'FECHADA', '2025-11-08 10:30:00');

-- Insere itens da nota fiscal de exemplo
INSERT INTO notas_itens (nota_id, produto_id, produto_codigo, produto_descricao, quantidade) VALUES
  ('nota-1', '1', 'PROD001', 'Notebook Dell Inspiron 15', 2),
  ('nota-1', '2', 'PROD002', 'Mouse Logitech MX Master 3', 5);

-- ========================================
-- VIEWS ÚTEIS PARA CONSULTAS
-- ========================================

-- View que mostra todas as notas com total de itens e quantidade total
CREATE VIEW vw_notas_resumo AS
SELECT 
  nf.id,
  nf.numero,
  nf.status,
  nf.data_emissao,
  COUNT(ni.id) as total_itens,
  COALESCE(SUM(ni.quantidade), 0) as quantidade_total
FROM notas_fiscais nf
LEFT JOIN notas_itens ni ON nf.id = ni.nota_id
GROUP BY nf.id, nf.numero, nf.status, nf.data_emissao
ORDER BY nf.numero DESC;

-- View que mostra os produtos com suas movimentações
CREATE VIEW vw_produtos_movimentacao AS
SELECT 
  p.id,
  p.codigo,
  p.descricao,
  p.saldo as saldo_atual,
  COALESCE(SUM(ni.quantidade), 0) as total_vendido
FROM produtos p
LEFT JOIN notas_itens ni ON p.id = ni.produto_id
GROUP BY p.id, p.codigo, p.descricao, p.saldo;

-- ========================================
-- PROCEDURE PARA ATUALIZAR ESTOQUE
-- ========================================
-- Esta procedure é chamada quando uma nota é impressa/finalizada

DELIMITER //

CREATE PROCEDURE sp_atualizar_estoque_nota(
  IN p_nota_id VARCHAR(36)
)
BEGIN
  -- Declara variáveis para controle
  DECLARE v_erro VARCHAR(255);
  DECLARE v_produto_id VARCHAR(36);
  DECLARE v_quantidade INT;
  DECLARE v_saldo_atual INT;
  DECLARE done INT DEFAULT FALSE;
  
  -- Cursor para iterar sobre os itens da nota
  DECLARE cursor_itens CURSOR FOR 
    SELECT produto_id, quantidade 
    FROM notas_itens 
    WHERE nota_id = p_nota_id;
  
  -- Handler para fim do cursor
  DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
  
  -- Inicia transação
  START TRANSACTION;
  
  -- Abre o cursor
  OPEN cursor_itens;
  
  -- Loop pelos itens
  read_loop: LOOP
    FETCH cursor_itens INTO v_produto_id, v_quantidade;
    
    IF done THEN
      LEAVE read_loop;
    END IF;
    
    -- Verifica saldo disponível
    SELECT saldo INTO v_saldo_atual 
    FROM produtos 
    WHERE id = v_produto_id;
    
    -- Se não tiver saldo suficiente, aborta
    IF v_saldo_atual < v_quantidade THEN
      SET v_erro = CONCAT('Saldo insuficiente para produto ID: ', v_produto_id);
      -- Rollback e retorna erro
      ROLLBACK;
      SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = v_erro;
    END IF;
    
    -- Atualiza o saldo
    UPDATE produtos 
    SET saldo = saldo - v_quantidade 
    WHERE id = v_produto_id;
    
  END LOOP;
  
  -- Fecha o cursor
  CLOSE cursor_itens;
  
  -- Atualiza o status da nota para FECHADA
  UPDATE notas_fiscais 
  SET status = 'FECHADA' 
  WHERE id = p_nota_id;
  
  -- Confirma a transação
  COMMIT;
  
END //

DELIMITER ;

-- ========================================
-- CONSULTAS ÚTEIS PARA VERIFICAÇÃO
-- ========================================

-- Ver todos os produtos
-- SELECT * FROM produtos;

-- Ver todas as notas
-- SELECT * FROM notas_fiscais;

-- Ver todos os itens de notas
-- SELECT * FROM notas_itens;

-- Ver resumo das notas
-- SELECT * FROM vw_notas_resumo;

-- Ver movimentação dos produtos
-- SELECT * FROM vw_produtos_movimentacao;

-- Ver itens de uma nota específica
-- SELECT * FROM notas_itens WHERE nota_id = 'nota-1';

-- ========================================
-- FIM DO SCRIPT
-- ========================================

SELECT 'Banco de dados KORP criado com sucesso!' AS mensagem;
SELECT COUNT(*) AS total_produtos FROM produtos;
SELECT COUNT(*) AS total_notas FROM notas_fiscais;
SELECT COUNT(*) AS total_itens FROM notas_itens;
