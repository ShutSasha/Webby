package grpc

import (
	"context"
	"webby-room-queue/internal/grpc/queuepb"

	"github.com/google/uuid"
)

type QueueService interface {
	MoveToTop(ctx context.Context, id uuid.UUID) error
}

type QueueServer struct {
	queuepb.UnimplementedQueueGrpcServiceServer
	service QueueService
}

func NewQueueServer(service QueueService) *QueueServer {
	return &QueueServer{service: service}
}

func (s *QueueServer) MoveToTop(
	ctx context.Context, req *queuepb.MoveToTopRequest,
) (*queuepb.MoveToTopResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	if err := s.service.MoveToTop(ctx, id); err != nil {
		return nil, err
	}

	return &queuepb.MoveToTopResponse{Success: true}, nil
}
