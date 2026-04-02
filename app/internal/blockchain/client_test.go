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

	// ASSERT: Validar que retornou um valor válido (não nil)
	if value == nil {
		t.Error("Expected non-nil value from contract")
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

// TestSetStorageValue_Integration testa escrita real no contrato deployado no Besu.
// REQUISITOS:
// - Besu rodando em localhost:8545
// - Contrato SimpleStorage deployado em 0x42699a7612a82f1d9c36148af9c77354759b210b
// - Private key configurada (Alice do genesis.json)
func TestSetStorageValue_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ARRANGE: Conectar ao Besu
	client, err := NewBesuClient(ctx, "http://localhost:8545")
	if err != nil {
		t.Fatalf("Failed to connect to Besu: %v", err)
	}
	defer client.Close()

	// ARRANGE: Carregar contrato com signer
	contractAddr := "0x42699a7612a82f1d9c36148af9c77354759b210b"
	privateKey := "0x8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"
	err = client.LoadContractWithSigner(ctx, contractAddr, privateKey)
	if err != nil {
		t.Fatalf("Failed to load contract with signer: %v", err)
	}

	// ARRANGE: Ler valor original para restaurar depois
	originalValue, err := client.GetStorageValue(ctx)
	if err != nil {
		t.Logf("Warning: could not read original value: %v", err)
		originalValue = big.NewInt(999) // fallback
	}

	// ACT: Escrever novo valor no contrato
	newValue := big.NewInt(123456)
	txHash, err := client.SetStorageValue(ctx, newValue)
	if err != nil {
		t.Fatalf("Failed to set storage value: %v", err)
	}

	t.Logf("✅ Transaction mined: %s", txHash)

	// ASSERT: Ler valor e verificar que foi atualizado
	readValue, err := client.GetStorageValue(ctx)
	if err != nil {
		t.Fatalf("Failed to read value after write: %v", err)
	}

	if readValue.Cmp(newValue) != 0 {
		t.Errorf("Expected value %s after write, got %s", newValue.String(), readValue.String())
	}

	t.Logf("✅ Successfully wrote and verified value: %s", readValue.String())

	// CLEANUP: Restaurar valor original
	_, err = client.SetStorageValue(ctx, originalValue)
	if err != nil {
		t.Logf("Warning: failed to restore original value: %v", err)
	} else {
		t.Logf("✅ Restored original value: %s", originalValue.String())
	}
}

// TestSetStorageValue_WithoutSigner verifica erro quando signer não foi configurado.
func TestSetStorageValue_WithoutSigner(t *testing.T) {
	ctx := context.Background()

	// ARRANGE: Cliente sem signer configurado
	client := &BesuClient{}

	// ACT: Tentar escrever sem configurar signer
	_, err := client.SetStorageValue(ctx, big.NewInt(100))

	// ASSERT: Deve retornar erro
	if err == nil {
		t.Error("Expected error when signer not configured, got nil")
	}

	if err.Error() != "contract not loaded: call LoadContractWithSigner() first" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestLoadContractWithSigner_EmptyPrivateKey verifica validação de private key vazia.
func TestLoadContractWithSigner_EmptyPrivateKey(t *testing.T) {
	ctx := context.Background()
	client := &BesuClient{}

	err := client.LoadContractWithSigner(ctx, "0x42699a7612a82f1d9c36148af9c77354759b210b", "")
	if err == nil {
		t.Error("Expected error for empty private key, got nil")
	}

	if err.Error() != "private key cannot be empty" {
		t.Errorf("Unexpected error message: %v", err)
	}
}
