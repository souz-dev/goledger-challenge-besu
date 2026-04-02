package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "challenge-besu/gen/pb"
	"challenge-besu/internal/repository"
	"challenge-besu/internal/service"
)

// StorageServer implementa a interface pb.StorageServiceServer.
type StorageServer struct {
	pb.UnimplementedStorageServiceServer
	service service.StorageService
}

// NewStorageServer cria uma nova instância do servidor gRPC.
func NewStorageServer(svc service.StorageService) *StorageServer {
	return &StorageServer{
		service: svc,
	}
}

// toGRPCError traduz erros de domínio em status gRPC.
func toGRPCError(err error) error {
	if errors.Is(err, repository.ErrValueNotSet) {
		return status.Error(codes.NotFound, "value not set")
	}
	return status.Error(codes.Internal, err.Error())
}

// SetValue implementa a RPC para definir um valor.
func (s *StorageServer) SetValue(ctx context.Context, req *pb.SetValueRequest) (*pb.SetValueResponse, error) {
	if err := s.service.SetValue(ctx, req.Value); err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.SetValueResponse{Success: true}, nil
}

// GetValue implementa a RPC para obter um valor.
func (s *StorageServer) GetValue(ctx context.Context, req *pb.GetValueRequest) (*pb.GetValueResponse, error) {
	value, err := s.service.GetValue(ctx)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.GetValueResponse{Value: value}, nil
}

// SyncValue implementa a RPC para sincronizar valores.
// TODO: Implementar quando blockchain for integrado.
func (s *StorageServer) SyncValue(ctx context.Context, req *pb.SyncValueRequest) (*pb.SyncValueResponse, error) {
	// Por enquanto, retorna sucesso sem fazer nada
	// Quando blockchain for implementado, chamará s.service.SyncValue(ctx)
	return &pb.SyncValueResponse{
		Success: true,
		Message: "sync will be implemented with blockchain integration",
	}, nil
}

// CheckValue implementa a RPC para verificar consistência de valores.
func (s *StorageServer) CheckValue(ctx context.Context, req *pb.CheckValueRequest) (*pb.CheckValueResponse, error) {
	equal, err := s.service.CheckValue(ctx, req.Value)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.CheckValueResponse{Equal: equal}, nil
}
