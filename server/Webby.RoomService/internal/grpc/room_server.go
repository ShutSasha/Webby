package grpc

import (
	"context"
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
	roomID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, err
	}
	room, err := s.roomGetter.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return &roompb.GetRoomHostResponse{
		HostId: room.HostID.String(),
	}, nil
}
