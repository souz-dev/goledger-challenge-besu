package repository

import (
	"context"
)

// storageRepository é uma implementação temporária/fake do StorageRepository.
// Será substituída por uma implementação real com banco de dados ou blockchain.
type storageRepository struct {
	// TODO: Adicionar campos de conexão com banco de dados
	// db *sql.DB
	// client *blockchain.Client
}

// NewStorageRepository cria uma nova instância do repositório de armazenamento.
func NewStorageRepository() *storageRepository {
	return &storageRepository{}
}

// SaveValue implementa o método da interface StorageRepository.
// Por enquanto, apenas simula a operação sem persistência real.
func (r *storageRepository) SaveValue(ctx context.Context, value int64) error {
	// TODO: Implementar persistência real
	// Exemplos futuros:
	// - Inserir em banco de dados (PostgreSQL, etc)
	// - Fazer transação na blockchain (Besu)
	// - Validar e sincronizar estado
	return nil
}

// GetValue implementa o método da interface StorageRepository.
// Por enquanto, retorna um valor padrão sem recuperação real.
func (r *storageRepository) GetValue(ctx context.Context) (int64, error) {
	// TODO: Implementar recuperação real
	// Exemplos futuros:
	// - Consultar banco de dados
	// - Ler estado da blockchain
	// - Cache em memória
	return 0, nil
}
