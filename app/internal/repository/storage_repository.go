package repository

import (
	"context"
)

// storageRepository é uma implementação temporária que mantém estado em memória.
// Será substituída por uma implementação real com banco de dados ou blockchain.
type storageRepository struct {
	value int64

	// TODO: Adicionar campos de conexão com banco de dados
	// db *sql.DB
	// client *blockchain.Client
}

// NewStorageRepository cria uma nova instância do repositório de armazenamento.
func NewStorageRepository() *storageRepository {
	return &storageRepository{}
}

// SaveValue implementa o método da interface StorageRepository.
// Armazena o valor em memória para uso temporário/desenvolvimento.
func (r *storageRepository) SaveValue(ctx context.Context, value int64) error {
	r.value = value
	// TODO: Implementar persistência real
	// Exemplos futuros:
	// - Inserir em banco de dados (PostgreSQL, etc)
	// - Fazer transação na blockchain (Besu)
	// - Validar e sincronizar estado
	return nil
}

// GetValue implementa o método da interface StorageRepository.
// Retorna o valor armazenado em memória.
func (r *storageRepository) GetValue(ctx context.Context) (int64, error) {
	return r.value, nil
	// TODO: Implementar recuperação real
	// Exemplos futuros:
	// - Consultar banco de dados
	// - Ler estado da blockchain
	// - Cache distribuído
}
