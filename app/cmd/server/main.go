package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "challenge-besu/gen/pb"
	"challenge-besu/internal/blockchain"
	"challenge-besu/internal/database"
	"challenge-besu/internal/repository"
	"challenge-besu/internal/service"
	grpcserver "challenge-besu/internal/transport/grpc"
)

func main() {
	ctx := context.Background()

	// Database configuration
	databaseURL := getEnv("DATABASE_URL", "postgres://admin:admin123@localhost:5433/challenge_besu?sslmode=disable")

	// Connect to PostgreSQL
	pool, err := database.NewPostgresPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("✅ Connected to PostgreSQL")

	// Connect to Besu node
	besuRPCURL := getEnv("BESU_RPC_URL", "http://localhost:8545")
	besuClient, err := blockchain.NewBesuClient(ctx, besuRPCURL)
	if err != nil {
		log.Fatalf("failed to connect to Besu node: %v", err)
	}
	defer besuClient.Close()

	// Validate Besu connection with health checks
	chainID, err := besuClient.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to get chain ID: %v", err)
	}

	blockNumber, err := besuClient.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("failed to get block number: %v", err)
	}

	log.Printf("✅ Connected to Besu (Chain ID: %s, Block: %d)", chainID.String(), blockNumber)

	// Load smart contract with signer
	contractAddress := getEnv("CONTRACT_ADDRESS", "0x42699a7612a82f1d9c36148af9c77354759b210b")
	privateKey := getEnv("PRIVATE_KEY", "0x8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63")

	if err := besuClient.LoadContractWithSigner(ctx, contractAddress, privateKey); err != nil {
		log.Fatalf("failed to load contract with signer: %v", err)
	}
	log.Printf("✅ Contract loaded at %s", contractAddress)

	// Initialize repository
	repo := repository.NewSQLRepository(pool)

	// Initialize service with blockchain integration
	storageService := service.NewStorageService(repo, besuClient)
	log.Println("✅ Service initialized with blockchain integration")

	// gRPC server configuration
	grpcPort := getEnv("GRPC_PORT", "50051")
	address := fmt.Sprintf(":%s", grpcPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10 * 1024 * 1024), // 10MB
	)

	// Register service
	storageServer := grpcserver.NewStorageServer(storageService)
	pb.RegisterStorageServiceServer(grpcServer, storageServer)

	// Enable reflection for grpcurl/Postman
	reflection.Register(grpcServer)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("received signal %v, shutting down gracefully...", sig)
		grpcServer.GracefulStop()
	}()

	// Start server
	log.Printf("🚀 gRPC server listening on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// getEnv returns environment variable or default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
