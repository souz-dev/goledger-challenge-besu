package service

import (
	"context"
)

// StorageService define o contrato para operações de armazenamento.
type StorageService interface {
	SetValue(ctx context.Context, value int64) (string, error)
	GetValue(ctx context.Context) (int64, error)
	CheckValue(ctx context.Context, value int64) (bool, error)
}
