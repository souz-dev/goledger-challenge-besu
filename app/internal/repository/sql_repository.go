package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// sqlRepository implementa StorageRepository usando PostgreSQL.
// Usa padrão Singleton com id=1 fixo na tabela storage.
type sqlRepository struct {
	pool *pgxpool.Pool
}

// NewSQLRepository cria uma nova instância do repositório SQL.
// O pool deve estar conectado e pronto antes de chamar este construtor.
func NewSQLRepository(pool *pgxpool.Pool) *sqlRepository {
	return &sqlRepository{
		pool: pool,
	}
}

// SaveValue atualiza o valor na tabela storage (sempre id=1).
// Usa UPDATE em vez de INSERT pois a migration garante que o registro existe.
func (r *sqlRepository) SaveValue(ctx context.Context, value int64) error {
	query := `
		UPDATE storage 
		SET value = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = 1
	`

	result, err := r.pool.Exec(ctx, query, value)
	if err != nil {
		return fmt.Errorf("erro ao salvar valor no banco: %w", err)
	}

	// Verificar se realmente atualizou uma linha
	if result.RowsAffected() == 0 {
		return errors.New("nenhuma linha atualizada: tabela storage pode estar vazia")
	}

	return nil
}

// GetValue retorna o valor armazenado na tabela storage (id=1).
// Retorna ErrValueNotSet se o registro não existir (não deve acontecer se migration rodou).
func (r *sqlRepository) GetValue(ctx context.Context) (int64, error) {
	query := `SELECT value FROM storage WHERE id = 1`

	var value int64
	err := r.pool.QueryRow(ctx, query).Scan(&value)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrValueNotSet
		}
		return 0, fmt.Errorf("erro ao buscar valor no banco: %w", err)
	}

	return value, nil
}
