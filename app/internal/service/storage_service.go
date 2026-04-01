package service

import (
	"context"
)

// StorageRepository define o contrato para operações de persistência de armazenamento.
// Esta interface permite que a camada de serviço seja independente da implementação
// de persistência (banco de dados, blockchain, etc).
type StorageRepository interface {
	SaveValue(ctx context.Context, value int64) error
	GetValue(ctx context.Context) (int64, error)
}

// storageService é a implementação concreta de StorageService.
// Encapsula a lógica de negócio e coordena as operações de armazenamento.
type storageService struct {
	repo StorageRepository
}

// NewStorageService cria uma nova instância do serviço de armazenamento com injeção de dependência.
// O repositório é passado como parâmetro, permitindo inversão de controle e facilitando testes.
func NewStorageService(repo StorageRepository) StorageService {
	return &storageService{
		repo: repo,
	}
}

// SetValue salva um valor usando o repositório.
// Retorna erro se a operação falhar.
func (s *storageService) SetValue(ctx context.Context, value int64) error {
	return s.repo.SaveValue(ctx, value)
}

// GetValue recupera um valor usando o repositório.
// Retorna o valor e erro se a operação falhar.
func (s *storageService) GetValue(ctx context.Context) (int64, error) {
	return s.repo.GetValue(ctx)
}

// SyncValue implementa a sincronização de valor com a blockchain.
// Por enquanto, apenas retorna nil (nenhuma operação realizada).
// Será expandido quando a camada de blockchain for implementada.
func (s *storageService) SyncValue(ctx context.Context) error {
	// TODO: Implementar sincronização com blockchain
	return nil
}

// CheckValue verifica a igualdade entre um valor recebido e o valor salvo no repositório.
// Retorna false + erro se nenhum valor foi salvo ainda.
func (s *storageService) CheckValue(ctx context.Context, value int64) (bool, error) {
	savedValue, err := s.repo.GetValue(ctx)
	if err != nil {
		return false, err
	}
	return value == savedValue, nil
}
