package service

import (
	"context"
)

// StorageService define o contrato para operações de armazenamento na blockchain.
type StorageService interface {
	SetValue(ctx context.Context, value int64) error
	GetValue(ctx context.Context) (int64, error)
	SyncValue(ctx context.Context) error
	CheckValue(ctx context.Context, value int64) (bool, error)
}
