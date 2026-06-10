package grpc

import (
	"context"
	"webby/room-queue-service/internal/grpc/queuepb"

	"github.com/google/uuid"
)

type queueService interface {
	MoveToTop(ctx context.Context, id uuid.UUID) error
}

type queueServer struct {
	queuepb.UnimplementedQueueGrpcServiceServer
	service queueService
}

func NewQueueServer(service queueService) *queueServer {
	return &queueServer{service: service}
}

func (s *queueServer) MoveToTop(ctx context.Context, req *queuepb.MoveToTopRequest) (*queuepb.MoveToTopResponse, error) {
	const op = "grpc.queueServer.MoveToTop"
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := s.service.MoveToTop(ctx, id); err != nil {
		return nil, err
	}

	return &queuepb.MoveToTopResponse{Success: true}, nil
}
