package repository

import (
	"context"
	"errors"
	"testing"
)

func TestGetValue_WithoutSave_ReturnsErrValueNotSet(t *testing.T) {
	repo := NewStorageRepository()

	_, err := repo.GetValue(context.Background())
	if !errors.Is(err, ErrValueNotSet) {
		t.Errorf("esperado ErrValueNotSet, got %v", err)
	}
}

func TestSaveAndGetValue_ReturnsCorrectValue(t *testing.T) {
	repo := NewStorageRepository()

	if err := repo.SaveValue(context.Background(), 42); err != nil {
		t.Fatalf("SaveValue inesperado erro: %v", err)
	}

	got, err := repo.GetValue(context.Background())
	if err != nil {
		t.Fatalf("GetValue inesperado erro: %v", err)
	}
	if got != 42 {
		t.Errorf("esperado 42, got %d", got)
	}
}

func TestSaveZero_GetValue_ReturnsZeroWithoutError(t *testing.T) {
	repo := NewStorageRepository()

	if err := repo.SaveValue(context.Background(), 0); err != nil {
		t.Fatalf("SaveValue inesperado erro: %v", err)
	}

	got, err := repo.GetValue(context.Background())
	if err != nil {
		t.Fatalf("esperado nil, got %v", err)
	}
	if got != 0 {
		t.Errorf("esperado 0, got %d", got)
	}
}
