package service

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"challenge-besu/internal/blockchain"
)

// StorageRepository define o contrato para operações de persistência.
// Interface mantida pois já existe implementação SQL e facilita testes.
type StorageRepository interface {
	SaveValue(ctx context.Context, value int64) error
	GetValue(ctx context.Context) (int64, error)
}

// storageService implementa a lógica de negócio da aplicação.
// Coordena operações entre blockchain (source of truth) e database (cache).
type storageService struct {
	repo       StorageRepository
	blockchain *blockchain.BesuClient
}

// NewStorageService cria uma nova instância do serviço.
func NewStorageService(repo StorageRepository, besuClient *blockchain.BesuClient) StorageService {
	return &storageService{
		repo:       repo,
		blockchain: besuClient,
	}
}

// SetValue writes value to blockchain (source of truth) and caches in database.
// Flow: blockchain write → database save
func (s *storageService) SetValue(ctx context.Context, value int64) error {
	// Convert int64 to *big.Int for blockchain
	bigValue := big.NewInt(value)

	// Write to blockchain (source of truth)
	txHash, err := s.blockchain.SetStorageValue(ctx, bigValue)
	if err != nil {
		return fmt.Errorf("failed to write value to blockchain: %w", err)
	}

	log.Printf("blockchain transaction mined: %s (value=%d)", txHash, value)

	// Cache in database for fast reads (best effort)
	if err := s.repo.SaveValue(ctx, value); err != nil {
		// Log error but don't fail - blockchain is source of truth
		log.Printf("warning: failed to cache value in database: %v", err)
	}

	return nil
}

// GetValue reads from blockchain (source of truth) and updates database cache.
// Flow: blockchain read → database update (optional)
func (s *storageService) GetValue(ctx context.Context) (int64, error) {
	// Read from blockchain (source of truth)
	bigValue, err := s.blockchain.GetStorageValue(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to read value from blockchain: %w", err)
	}

	// Convert *big.Int to int64
	value := bigValue.Int64()

	// Update database cache (best effort)
	if err := s.repo.SaveValue(ctx, value); err != nil {
		// Log error but don't fail - blockchain is source of truth
		log.Printf("warning: failed to update database cache: %v", err)
	}

	return value, nil
}

// CheckValue verifica se o valor informado corresponde ao valor armazenado.
// Retorna true se os valores forem iguais, false caso contrário.
func (s *storageService) CheckValue(ctx context.Context, value int64) (bool, error) {
	storedValue, err := s.GetValue(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get stored value: %w", err)
	}
	return value == storedValue, nil
}
