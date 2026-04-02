package repository

import (
	"context"
	"errors"
	"sync"
)

var ErrValueNotSet = errors.New("value not set")

type storageRepository struct {
	mu       sync.RWMutex
	value    int64
	hasValue bool
}

func NewStorageRepository() *storageRepository {
	return &storageRepository{}
}

func (r *storageRepository) SaveValue(ctx context.Context, value int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.value = value
	r.hasValue = true
	return nil
}

func (r *storageRepository) GetValue(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if !r.hasValue {
		return 0, ErrValueNotSet
	}
	return r.value, nil
}
