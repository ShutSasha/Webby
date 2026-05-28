package grpc

import (
	"context"
	"fmt"
	"webby/room-category-service/internal/grpc/categorypb"
)

type CategoryService interface {
	Exists(ctx context.Context, name string) (bool, error)
}

type CategoryServer struct {
	categorypb.UnimplementedCategoryGrpcServiceServer
	service CategoryService
}

func NewCategoryServer(service CategoryService) *CategoryServer {
	return &CategoryServer{service: service}
}

func (s *CategoryServer) CategoryExists(
	ctx context.Context,
	req *categorypb.CategoryExistsRequest,
) (*categorypb.CategoryExistsResponse, error) {
	const op = "grpc.CategoryServer.CategoryExists"

	exists, err := s.service.Exists(ctx, req.GetName())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &categorypb.CategoryExistsResponse{Exists: exists}, nil
}
