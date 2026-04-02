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
	"challenge-besu/internal/config"
	"challenge-besu/internal/database"
	"challenge-besu/internal/repository"
	"challenge-besu/internal/service"
	grpcserver "challenge-besu/internal/transport/grpc"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Connect to PostgreSQL
	pool, err := database.NewPostgresPool(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("✅ Connected to PostgreSQL")

	// Connect to Besu node
	besuClient, err := blockchain.NewBesuClient(ctx, cfg.Besu.RPCURL)
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
	if err := besuClient.LoadContractWithSigner(ctx, cfg.Besu.ContractAddress, cfg.Besu.PrivateKey); err != nil {
		log.Fatalf("failed to load contract with signer: %v", err)
	}
	log.Printf("✅ Contract loaded at %s", cfg.Besu.ContractAddress)

	// Initialize repository
	repo := repository.NewSQLRepository(pool)

	// Initialize service with blockchain integration
	storageService := service.NewStorageService(repo, besuClient)
	log.Println("✅ Service initialized with blockchain integration")

	// gRPC server configuration
	address := fmt.Sprintf(":%s", cfg.GRPCPort)

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
