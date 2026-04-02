package blockchain

import (
	"context"
	"math/big"
	"testing"
	"time"
)

// TestGetStorageValue_Integration testa leitura real do contrato deployado no Besu.
// REQUISITOS:
// - Besu rodando em localhost:8545
// - Contrato SimpleStorage deployado em 0x42699a7612a82f1d9c36148af9c77354759b210b
// - Valor já setado no contrato via cast send (999)
func TestGetStorageValue_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ARRANGE: Conectar ao Besu
	client, err := NewBesuClient(ctx, "http://localhost:8545")
	if err != nil {
		t.Fatalf("Failed to connect to Besu: %v", err)
	}
	defer client.Close()

	// ACT: Carregar contrato
	contractAddr := "0x42699a7612a82f1d9c36148af9c77354759b210b"
	err = client.LoadContract(contractAddr)
	if err != nil {
		t.Fatalf("Failed to load contract: %v", err)
	}

	// ACT: Ler valor do contrato
	value, err := client.GetStorageValue(ctx)
	if err != nil {
		t.Fatalf("Failed to get storage value: %v", err)
	}

	// ASSERT: Validar valor esperado (999 foi setado anteriormente)
	expectedValue := big.NewInt(999)
	if value.Cmp(expectedValue) != 0 {
		t.Errorf("Expected value %s, got %s", expectedValue.String(), value.String())
	}

	t.Logf("✅ Successfully read value from contract: %s", value.String())
}

// TestGetStorageValue_WithoutLoadContract verifica erro quando contrato não foi carregado.
func TestGetStorageValue_WithoutLoadContract(t *testing.T) {
	ctx := context.Background()

	// ARRANGE: Cliente sem contrato carregado
	client := &BesuClient{}

	// ACT: Tentar ler valor sem carregar contrato
	_, err := client.GetStorageValue(ctx)

	// ASSERT: Deve retornar erro
	if err == nil {
		t.Error("Expected error when contract not loaded, got nil")
	}

	if err.Error() != "contract not loaded: call LoadContract() first" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestLoadContract_EmptyAddress verifica validação de endereço vazio.
func TestLoadContract_EmptyAddress(t *testing.T) {
	client := &BesuClient{}

	err := client.LoadContract("")
	if err == nil {
		t.Error("Expected error for empty address, got nil")
	}

	if err.Error() != "contract address cannot be empty" {
		t.Errorf("Unexpected error message: %v", err)
	}
}
