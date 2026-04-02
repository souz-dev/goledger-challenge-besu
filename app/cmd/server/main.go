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

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	pool, err := database.NewPostgresPool(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("✅ Connected to PostgreSQL")

	besuClient, err := blockchain.NewBesuClient(ctx, cfg.Besu.RPCURL)
	if err != nil {
		log.Fatalf("failed to connect to Besu node: %v", err)
	}
	defer besuClient.Close()

	chainID, err := besuClient.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to get chain ID: %v", err)
	}

	blockNumber, err := besuClient.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("failed to get block number: %v", err)
	}

	log.Printf("✅ Connected to Besu (Chain ID: %s, Block: %d)", chainID.String(), blockNumber)

	if err := besuClient.LoadContractWithSigner(ctx, cfg.Besu.ContractAddress, cfg.Besu.PrivateKey); err != nil {
		log.Fatalf("failed to load contract with signer: %v", err)
	}
	log.Printf("✅ Contract loaded at %s", cfg.Besu.ContractAddress)

	repo := repository.NewSQLRepository(pool)
	storageService := service.NewStorageService(repo, besuClient)
	log.Println("✅ Service initialized with blockchain integration")

	address := fmt.Sprintf(":%s", cfg.GRPCPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10 * 1024 * 1024),
	)

	storageServer := grpcserver.NewStorageServer(storageService)
	pb.RegisterStorageServiceServer(grpcServer, storageServer)
	reflection.Register(grpcServer)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("received signal %v, shutting down gracefully...", sig)
		grpcServer.GracefulStop()
	}()

	log.Printf("🚀 gRPC server listening on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
