package service

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"challenge-besu/internal/blockchain"
)

type StorageRepository interface {
	SaveValue(ctx context.Context, value int64) error
	GetValue(ctx context.Context) (int64, error)
}

type storageService struct {
	repo       StorageRepository
	blockchain *blockchain.BesuClient
}

func NewStorageService(repo StorageRepository, besuClient *blockchain.BesuClient) StorageService {
	return &storageService{
		repo:       repo,
		blockchain: besuClient,
	}
}

func (s *storageService) SetValue(ctx context.Context, value int64) (string, error) {
	bigValue := big.NewInt(value)

	txHash, err := s.blockchain.SetStorageValue(ctx, bigValue)
	if err != nil {
		return "", fmt.Errorf("failed to write value to blockchain: %w", err)
	}

	log.Printf("blockchain transaction mined: %s (value=%d)", txHash, value)

	if err := s.repo.SaveValue(ctx, value); err != nil {
		log.Printf("warning: failed to cache value in database: %v", err)
	}

	return txHash, nil
}

func (s *storageService) GetValue(ctx context.Context) (int64, error) {
	bigValue, err := s.blockchain.GetStorageValue(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to read value from blockchain: %w", err)
	}

	value := bigValue.Int64()

	if err := s.repo.SaveValue(ctx, value); err != nil {
		log.Printf("warning: failed to update database cache: %v", err)
	}

	return value, nil
}

func (s *storageService) CheckValue(ctx context.Context, value int64) (bool, error) {
	storedValue, err := s.GetValue(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get stored value: %w", err)
	}
	return value == storedValue, nil
}
