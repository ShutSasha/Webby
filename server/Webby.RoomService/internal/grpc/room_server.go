package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/roompb"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type RoomGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
}

type RoomServer struct {
	roompb.UnimplementedRoomGrpcServiceServer
	roomGetter RoomGetter
}

func NewRoomServer(roomGetter RoomGetter) *RoomServer {
	return &RoomServer{roomGetter: roomGetter}
}

func (s *RoomServer) GetRoomHost(ctx context.Context, req *roompb.GetRoomHostRequest) (*roompb.GetRoomHostResponse, error) {
	const op = "grpc.RoomService.GetRoomHost"
	roomID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	room, err := s.roomGetter.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &roompb.GetRoomHostResponse{
		HostId: room.HostID.String(),
	}, nil
}
