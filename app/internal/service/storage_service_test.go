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

// --- SetValue ---

func TestSetValue_Success(t *testing.T) {
	repo := &mockRepository{}
	svc := NewStorageService(repo)

	err := svc.SetValue(context.Background(), 42)
	if err != nil {
		t.Fatalf("esperado nil, got %v", err)
	}
	if repo.savedValue != 42 {
		t.Errorf("esperado savedValue=42, got %d", repo.savedValue)
	}
}

func TestSetValue_RepositoryError(t *testing.T) {
	repo := &mockRepository{saveErr: errors.New("save failed")}
	svc := NewStorageService(repo)

	err := svc.SetValue(context.Background(), 10)
	if err == nil {
		t.Fatal("esperado erro, got nil")
	}
}

// --- GetValue ---

func TestGetValue_Success(t *testing.T) {
	repo := &mockRepository{savedValue: 99, hasValue: true}
	svc := NewStorageService(repo)

	got, err := svc.GetValue(context.Background())
	if err != nil {
		t.Fatalf("esperado nil, got %v", err)
	}
	if got != 99 {
		t.Errorf("esperado 99, got %d", got)
	}
}

func TestGetValue_RepositoryError(t *testing.T) {
	repo := &mockRepository{getErr: errors.New("value not set")}
	svc := NewStorageService(repo)

	_, err := svc.GetValue(context.Background())
	if err == nil {
		t.Fatal("esperado erro, got nil")
	}
}

// --- CheckValue ---

func TestCheckValue_Equal(t *testing.T) {
	repo := &mockRepository{savedValue: 25, hasValue: true}
	svc := NewStorageService(repo)

	equal, err := svc.CheckValue(context.Background(), 25)
	if err != nil {
		t.Fatalf("esperado nil, got %v", err)
	}
	if !equal {
		t.Error("esperado true, got false")
	}
}

func TestCheckValue_NotEqual(t *testing.T) {
	repo := &mockRepository{savedValue: 25, hasValue: true}
	svc := NewStorageService(repo)

	equal, err := svc.CheckValue(context.Background(), 10)
	if err != nil {
		t.Fatalf("esperado nil, got %v", err)
	}
	if equal {
		t.Error("esperado false, got true")
	}
}

func TestCheckValue_RepositoryError(t *testing.T) {
	repo := &mockRepository{getErr: errors.New("value not set")}
	svc := NewStorageService(repo)

	equal, err := svc.CheckValue(context.Background(), 0)
	if err == nil {
		t.Fatal("esperado erro, got nil")
	}
	if equal {
		t.Error("esperado false quando há erro, got true")
	}
}

func TestCheckValue_ZeroAfterSet(t *testing.T) {
	repo := &mockRepository{}
	svc := NewStorageService(repo)

	_ = svc.SetValue(context.Background(), 0)

	equal, err := svc.CheckValue(context.Background(), 0)
	if err != nil {
		t.Fatalf("esperado nil, got %v", err)
	}
	if !equal {
		t.Error("esperado true para SetValue(0) + CheckValue(0), got false")
	}
}
