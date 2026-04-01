package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool cria um pool de conexões com PostgreSQL.
// O pool gerencia automaticamente conexões, retries e health checks.
//
// Parâmetros:
//   - ctx: contexto para controle de timeout/cancelamento
//   - databaseURL: string de conexão no formato postgres://user:pass@host:port/dbname
//
// Retorna erro se não conseguir conectar ou fazer ping no banco.
func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	// Configuração do pool (usa defaults do pgx que são otimizados)
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear database URL: %w", err)
	}

	// Ajustes opcionais do pool (podem ser customizados conforme carga)
	config.MaxConns = 10       // Máximo de conexões simultâneas
	config.MinConns = 2        // Mínimo de conexões idle
	config.MaxConnLifetime = 0 // Sem limite de tempo de vida
	config.MaxConnIdleTime = 0 // Sem timeout de idle
	// HealthCheckPeriod usa o default do pgx (1 minuto)

	// Criar pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar pool de conexões: %w", err)
	}

	// Testar conectividade
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	return pool, nil
}
