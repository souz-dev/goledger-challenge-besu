package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// BesuClient encapsula a conexão RPC com o nó Hyperledger Besu.
// Responsável por comunicação low-level com a blockchain via JSON-RPC.
type BesuClient struct {
	client *ethclient.Client
	rpcURL string
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
