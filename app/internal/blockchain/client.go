package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"challenge-besu/internal/blockchain/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BesuClient encapsula a conexão RPC com o nó Hyperledger Besu.
// Responsável por comunicação low-level com a blockchain via JSON-RPC
// e interação com smart contracts deployados.
type BesuClient struct {
	client   *ethclient.Client
	rpcURL   string
	contract *contracts.SimpleStorage // Binding do smart contract (opcional)
	address  common.Address           // Endereço do contrato (se carregado)
}

// NewBesuClient cria uma nova conexão com o nó Besu.
// Retorna erro se não conseguir estabelecer conexão RPC.
func NewBesuClient(ctx context.Context, rpcURL string) (*BesuClient, error) {
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Besu node at %s: %w", rpcURL, err)
	}

	return &BesuClient{
		client: client,
		rpcURL: rpcURL,
	}, nil
}

// ChainID retorna o identificador da rede blockchain.
// Usado como health check para validar que a conexão está viva.
// Para Besu local do desafio, deve retornar 1337.
func (c *BesuClient) ChainID(ctx context.Context) (*big.Int, error) {
	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}
	return chainID, nil
}

// BlockNumber retorna o número do bloco mais recente.
// Útil para verificar que o nó está sincronizado e processando blocos.
func (c *BesuClient) BlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}
	return blockNumber, nil
}

// Close encerra a conexão com o nó Besu.
// Deve ser chamado no shutdown da aplicação.
func (c *BesuClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

// LoadContract carrega o binding do smart contract SimpleStorage.
// Deve ser chamado após NewBesuClient se houver necessidade de interagir com o contrato.
// contractAddress: endereço hexadecimal do contrato deployado (ex: 0x42699...)
func (c *BesuClient) LoadContract(contractAddress string) error {
	if contractAddress == "" {
		return fmt.Errorf("contract address cannot be empty")
	}

	// Converter string hexadecimal para common.Address
	address := common.HexToAddress(contractAddress)
	c.address = address

	// Criar binding do contrato usando o ABI gerado
	contract, err := contracts.NewSimpleStorage(address, c.client)
	if err != nil {
		return fmt.Errorf("failed to instantiate contract at %s: %w", contractAddress, err)
	}

	c.contract = contract
	return nil
}

// GetStorageValue lê o valor armazenado no smart contract (método get()).
// Esta é uma operação read-only (view call) que NÃO consome gas.
// Retorna erro se o contrato não foi carregado via LoadContract().
func (c *BesuClient) GetStorageValue(ctx context.Context) (*big.Int, error) {
	if c.contract == nil {
		return nil, fmt.Errorf("contract not loaded: call LoadContract() first")
	}

	// Configurar opções de chamada com context (permite timeout/cancelamento)
	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	// Chamar método get() do contrato (view call, sem transação)
	value, err := c.contract.Get(callOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract.Get(): %w", err)
	}

	return value, nil
}
