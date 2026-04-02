package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"challenge-besu/internal/blockchain/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BesuClient struct {
	client   *ethclient.Client
	rpcURL   string
	contract *contracts.SimpleStorage
	address  common.Address
	auth     *bind.TransactOpts
	chainID  *big.Int
}

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

func (c *BesuClient) ChainID(ctx context.Context) (*big.Int, error) {
	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}
	return chainID, nil
}

func (c *BesuClient) BlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}
	return blockNumber, nil
}

func (c *BesuClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

func (c *BesuClient) LoadContract(contractAddress string) error {
	if contractAddress == "" {
		return fmt.Errorf("contract address cannot be empty")
	}

	address := common.HexToAddress(contractAddress)
	c.address = address

	contract, err := contracts.NewSimpleStorage(address, c.client)
	if err != nil {
		return fmt.Errorf("failed to instantiate contract at %s: %w", contractAddress, err)
	}

	c.contract = contract
	return nil
}

func (c *BesuClient) GetStorageValue(ctx context.Context) (*big.Int, error) {
	if c.contract == nil {
		return nil, fmt.Errorf("contract not loaded: call LoadContract() first")
	}

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	value, err := c.contract.Get(callOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract.Get(): %w", err)
	}

	return value, nil
}

func (c *BesuClient) LoadContractWithSigner(ctx context.Context, contractAddress, privateKeyHex string) error {
	if contractAddress == "" {
		return fmt.Errorf("contract address cannot be empty")
	}
	if privateKeyHex == "" {
		return fmt.Errorf("private key cannot be empty")
	}

	if err := c.LoadContract(contractAddress); err != nil {
		return err
	}

	chainID, err := c.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}
	c.chainID = chainID

	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	c.auth = auth
	return nil
}

func (c *BesuClient) SetStorageValue(ctx context.Context, value *big.Int) (string, error) {
	if c.contract == nil {
		return "", fmt.Errorf("contract not loaded: call LoadContractWithSigner() first")
	}
	if c.auth == nil {
		return "", fmt.Errorf("signer not configured: call LoadContractWithSigner() first")
	}

	c.auth.Context = ctx

	tx, err := c.contract.Set(c.auth, value)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, c.client, tx)
	if err != nil {
		return "", fmt.Errorf("failed to wait for transaction mining: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return "", fmt.Errorf("transaction failed with status %d (reverted)", receipt.Status)
	}

	return tx.Hash().Hex(), nil
}
