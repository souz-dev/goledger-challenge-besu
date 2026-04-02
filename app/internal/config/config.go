package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database DatabaseConfig
	Besu     BesuConfig
	GRPCPort string
}

type DatabaseConfig struct {
	URL string
}

type BesuConfig struct {
	RPCURL          string
	ContractAddress string
	PrivateKey      string
}

func Load() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://admin:admin123@localhost:5433/challenge_besu?sslmode=disable"),
		},
		Besu: BesuConfig{
			RPCURL:          getEnv("BESU_RPC_URL", "http://localhost:8545"),
			ContractAddress: getEnv("CONTRACT_ADDRESS", "0x42699a7612a82f1d9c36148af9c77354759b210b"),
			PrivateKey:      getEnv("PRIVATE_KEY", "0x8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"),
		},
		GRPCPort: getEnv("GRPC_PORT", "50051"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.Besu.RPCURL == "" {
		return fmt.Errorf("BESU_RPC_URL is required")
	}
	if c.Besu.ContractAddress == "" {
		return fmt.Errorf("CONTRACT_ADDRESS is required")
	}
	if c.Besu.PrivateKey == "" {
		return fmt.Errorf("PRIVATE_KEY is required")
	}
	if c.GRPCPort == "" {
		return fmt.Errorf("GRPC_PORT is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
