package repository

import (
	"context"
	"errors"
	"sync"
)

// ErrValueNotSet é retornado quando GetValue é chamado antes de qualquer SaveValue.
var ErrValueNotSet = errors.New("value not set")

// storageRepository é uma implementação em memória com suporte a concorrência.
// Protegida por mutex para uso seguro em ambiente gRPC concorrente.
type storageRepository struct {
	mu       sync.RWMutex
	value    int64
	hasValue bool
}

// NewStorageRepository cria uma nova instância do repositório de armazenamento.
func NewStorageRepository() *storageRepository {
	return &storageRepository{}
}

// SaveValue armazena o valor em memória de forma segura para concorrência.
func (r *storageRepository) SaveValue(ctx context.Context, value int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.value = value
	r.hasValue = true
	return nil
}

// GetValue retorna o valor armazenado em memória de forma segura para concorrência.
// Retorna ErrValueNotSet se nenhum valor foi salvo ainda.
func (r *storageRepository) GetValue(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if !r.hasValue {
		return 0, ErrValueNotSet
	}
	return r.value, nil
}
