package grpc

import (
	"context"
	"fmt"
	"webby/room-queue-service/internal/grpc/queuepb"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type nextMover interface {
	MakeNext(ctx context.Context, roomID, queueItemID uuid.UUID) error
}

type queueServer struct {
	queuepb.UnimplementedQueueGrpcServiceServer
	nextMover nextMover
}

func NewQueueServer(nextMover nextMover) *queueServer {
	return &queueServer{nextMover: nextMover}
}

func (s *queueServer) MakeNext(ctx context.Context, req *queuepb.MakeNextRequest) (*emptypb.Empty, error) {
	const op = "grpc.queueServer.MakeNext"

	roomID, err := uuid.Parse(req.GetRoomID())
	if err != nil {
		return nil, fmt.Errorf("%s: parse roomID %w", op, err)
	}

	queueItemID, err := uuid.Parse(req.GetQueueItemID())
	if err != nil {
		return nil, fmt.Errorf("%s: parse queueItemID %w", op, err)
	}

	err = s.nextMover.MakeNext(ctx, roomID, queueItemID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &emptypb.Empty{}, nil
}
