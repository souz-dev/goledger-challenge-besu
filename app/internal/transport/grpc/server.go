package grpc

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "challenge-besu/gen/pb"
	"challenge-besu/internal/repository"
	"challenge-besu/internal/service"
)

// StorageServer implements pb.StorageServiceServer interface.
type StorageServer struct {
	pb.UnimplementedStorageServiceServer
	service service.StorageService
}

// NewStorageServer creates a new gRPC server instance.
func NewStorageServer(svc service.StorageService) *StorageServer {
	return &StorageServer{
		service: svc,
	}
}

// toGRPCError translates domain errors to gRPC status codes.
func toGRPCError(err error) error {
	if errors.Is(err, repository.ErrValueNotSet) {
		return status.Error(codes.NotFound, "value not set")
	}
	return status.Error(codes.Internal, err.Error())
}

// SetValue implements RPC to write a value to the blockchain.
// Returns transaction hash on success.
func (s *StorageServer) SetValue(ctx context.Context, req *pb.SetValueRequest) (*pb.SetValueResponse, error) {
	txHash, err := s.service.SetValue(ctx, req.Value)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.SetValueResponse{
		Success: true,
		TxHash:  txHash,
	}, nil
}

// GetValue implements RPC to read a value from the blockchain.
func (s *StorageServer) GetValue(ctx context.Context, req *pb.GetValueRequest) (*pb.GetValueResponse, error) {
	value, err := s.service.GetValue(ctx)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.GetValueResponse{Value: value}, nil
}

// SyncValue implements RPC to sync database cache with blockchain state.
func (s *StorageServer) SyncValue(ctx context.Context, req *pb.SyncValueRequest) (*pb.SyncValueResponse, error) {
	// GetValue already syncs database cache as side effect
	_, err := s.service.GetValue(ctx)
	if err != nil {
		return &pb.SyncValueResponse{
			Success: false,
			Message: fmt.Sprintf("sync failed: %v", err),
		}, nil
	}
	return &pb.SyncValueResponse{
		Success: true,
		Message: "database cache synchronized with blockchain",
	}, nil
}

// CheckValue implements RPC to check if value matches stored value.
func (s *StorageServer) CheckValue(ctx context.Context, req *pb.CheckValueRequest) (*pb.CheckValueResponse, error) {
	equal, err := s.service.CheckValue(ctx, req.Value)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.CheckValueResponse{Equal: equal}, nil
}
