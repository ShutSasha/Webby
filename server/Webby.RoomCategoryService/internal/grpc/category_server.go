package grpc

import (
	"context"
	"fmt"
	"webby/room-category-service/internal/grpc/categorypb"
)

type categoryService interface {
	Exists(ctx context.Context, name string) (bool, error)
}

type categoryServer struct {
	categorypb.UnimplementedCategoryGrpcServiceServer
	service categoryService
}

func NewCategoryServer(service categoryService) *categoryServer {
	return &categoryServer{service: service}
}

func (s *categoryServer) CategoryExists(
	ctx context.Context,
	req *categorypb.CategoryExistsRequest,
) (*categorypb.CategoryExistsResponse, error) {
	const op = "grpc.categoryServer.CategoryExists"

	exists, err := s.service.Exists(ctx, req.GetName())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &categorypb.CategoryExistsResponse{Exists: exists}, nil
}
