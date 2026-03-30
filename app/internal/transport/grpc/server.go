package grpc

import (
	"context"

	pb "challenge-besu/gen/pb"
	"challenge-besu/internal/service"
)

// StorageServer implementa a interface pb.StorageServiceServer.
type StorageServer struct {
	pb.UnimplementedStorageServiceServer
	service service.StorageService
}

// NewStorageServer cria uma nova instância do servidor gRPC com dependency injection.
func NewStorageServer(svc service.StorageService) *StorageServer {
	return &StorageServer{
		service: svc,
	}
}

// SetValue implementa a RPC para definir um valor no armazenamento.
// Delega a operação para a camada de serviço.
func (s *StorageServer) SetValue(ctx context.Context, req *pb.SetValueRequest) (*pb.SetValueResponse, error) {
	if err := s.service.SetValue(ctx, req.Value); err != nil {
		return nil, err
	}
	return &pb.SetValueResponse{Success: true}, nil
}

// GetValue implementa a RPC para obter um valor do armazenamento.
// Delega a operação para a camada de serviço.
func (s *StorageServer) GetValue(ctx context.Context, req *pb.GetValueRequest) (*pb.GetValueResponse, error) {
	value, err := s.service.GetValue(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetValueResponse{Value: value}, nil
}

// SyncValue implementa a RPC para sincronizar um valor com a blockchain.
// Delega a operação para a camada de serviço.
func (s *StorageServer) SyncValue(ctx context.Context, req *pb.SyncValueRequest) (*pb.SyncValueResponse, error) {
	if err := s.service.SyncValue(ctx); err != nil {
		return nil, err
	}
	return &pb.SyncValueResponse{Success: true}, nil
}

// CheckValue implementa a RPC para verificar a igualdade de um valor.
// Delega a operação para a camada de serviço.
func (s *StorageServer) CheckValue(ctx context.Context, req *pb.CheckValueRequest) (*pb.CheckValueResponse, error) {
	equal, err := s.service.CheckValue(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.CheckValueResponse{Equal: equal}, nil
}
