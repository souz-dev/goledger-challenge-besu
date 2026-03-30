package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "challenge-besu/gen/pb"
	"challenge-besu/internal/repository"
	"challenge-besu/internal/service"
	grpcserver "challenge-besu/internal/transport/grpc"
)

func main() {
	// Definir endereço de escuta
	address := "localhost:50051"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("falha ao criar listener: %v", err)
	}
	defer listener.Close()

	// Criar servidor gRPC
	grpcServer := grpc.NewServer()

	// CAMADA 1: Inicializar Repository (camada de persistência)
	storageRepo := repository.NewStorageRepository()

	// CAMADA 2: Inicializar Service (camada de lógica de negócio)
	storageService := service.NewStorageService(storageRepo)

	// CAMADA 3: Inicializar gRPC Server (camada de transporte)
	storageServerInstance := grpcserver.NewStorageServer(storageService)

	// CAMADA 4: Registrar o servidor de armazenamento no gRPC
	pb.RegisterStorageServiceServer(grpcServer, storageServerInstance)

	log.Printf("gRPC server listening on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("falha ao servir: %v", err)
	}
}
