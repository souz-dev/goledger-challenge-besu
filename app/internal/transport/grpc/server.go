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

type StorageServer struct {
	pb.UnimplementedStorageServiceServer
	service service.StorageService
}

func NewStorageServer(svc service.StorageService) *StorageServer {
	return &StorageServer{
		service: svc,
	}
}

func toGRPCError(err error) error {
	if errors.Is(err, repository.ErrValueNotSet) {
		return status.Error(codes.NotFound, "value not set")
	}
	return status.Error(codes.Internal, err.Error())
}

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

func (s *StorageServer) GetValue(ctx context.Context, req *pb.GetValueRequest) (*pb.GetValueResponse, error) {
	value, err := s.service.GetValue(ctx)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.GetValueResponse{Value: value}, nil
}

func (s *StorageServer) SyncValue(ctx context.Context, req *pb.SyncValueRequest) (*pb.SyncValueResponse, error) {
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

func (s *StorageServer) CheckValue(ctx context.Context, req *pb.CheckValueRequest) (*pb.CheckValueResponse, error) {
	equal, err := s.service.CheckValue(ctx, req.Value)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &pb.CheckValueResponse{Equal: equal}, nil
}
