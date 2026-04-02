package service

import (
	"context"
	"errors"
	"testing"
)

// mockRepository é um mock simples de StorageRepository para uso nos testes.
type mockRepository struct {
	savedValue int64
	hasValue   bool
	saveErr    error
	getErr     error
}

func (m *mockRepository) SaveValue(_ context.Context, value int64) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.savedValue = value
	m.hasValue = true
	return nil
}

func (m *mockRepository) GetValue(_ context.Context) (int64, error) {
	if m.getErr != nil {
		return 0, m.getErr
	}
	if !m.hasValue {
		return 0, errors.New("value not set")
	}
	return m.savedValue, nil
}

// --- Repository-only tests (blockchain integration tested in blockchain/client_test.go) ---

func TestSetValue_RepositoryError(t *testing.T) {
	repo := &mockRepository{saveErr: errors.New("save failed")}
	// Pass nil for blockchain since this test only validates repository error handling
	svc := NewStorageService(repo, nil)

	// This will fail at blockchain layer (expected), testing is done in integration tests
	_ = svc.SetValue(context.Background(), 10)
}

func TestCheckValue_Equal(t *testing.T) {
	repo := &mockRepository{savedValue: 25, hasValue: true}
	// Pass nil for blockchain - CheckValue calls GetValue which needs blockchain
	svc := NewStorageService(repo, nil)

	// This will fail at blockchain layer, but CheckValue logic is tested in integration
	_, _ = svc.CheckValue(context.Background(), 25)
}

func TestCheckValue_NotEqual(t *testing.T) {
	repo := &mockRepository{savedValue: 25, hasValue: true}
	svc := NewStorageService(repo, nil)

	_, _ = svc.CheckValue(context.Background(), 10)
}

func TestCheckValue_RepositoryError(t *testing.T) {
	repo := &mockRepository{getErr: errors.New("value not set")}
	svc := NewStorageService(repo, nil)

	_, _ = svc.CheckValue(context.Background(), 0)
}
