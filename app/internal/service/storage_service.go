package service

import (
	"context"
	"fmt"
)

// StorageRepository define o contrato para operações de persistência.
// Interface mantida pois já existe implementação SQL e facilita testes.
type StorageRepository interface {
	SaveValue(ctx context.Context, value int64) error
	GetValue(ctx context.Context) (int64, error)
}

// storageService implementa a lógica de negócio da aplicação.
// Por enquanto, atua como thin layer entre transport e repository.
// Quando blockchain for integrado, coordenará múltiplas fontes de dados.
type storageService struct {
	repo StorageRepository
}

// NewStorageService cria uma nova instância do serviço.
func NewStorageService(repo StorageRepository) StorageService {
	return &storageService{
		repo: repo,
	}
}

// SetValue persiste um valor no banco de dados.
func (s *storageService) SetValue(ctx context.Context, value int64) error {
	return s.repo.SaveValue(ctx, value)
}

// GetValue recupera o valor armazenado no banco de dados.
func (s *storageService) GetValue(ctx context.Context) (int64, error) {
	return s.repo.GetValue(ctx)
}

// CheckValue verifica se o valor informado corresponde ao valor armazenado.
// Retorna true se os valores forem iguais, false caso contrário.
func (s *storageService) CheckValue(ctx context.Context, value int64) (bool, error) {
	storedValue, err := s.repo.GetValue(ctx)
	if err != nil {
		return false, fmt.Errorf("erro ao buscar valor armazenado: %w", err)
	}
	return value == storedValue, nil
}
