package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/roompb"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roomGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
	GetTotalRooms(ctx context.Context) (int, error)
}

type roomServer struct {
	roompb.UnimplementedRoomGrpcServiceServer
	roomGetter roomGetter
}

func NewRoomServer(roomGetter roomGetter) *roomServer {
	return &roomServer{roomGetter: roomGetter}
}

func (s *roomServer) GetRoomHost(ctx context.Context, req *roompb.GetRoomHostRequest) (*roompb.GetRoomHostResponse, error) {
	const op = "grpc.RoomService.GetRoomHost"
	roomID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	room, err := s.roomGetter.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &roompb.GetRoomHostResponse{HostId: room.HostID.String()}, nil
}

func (s *roomServer) GetTotalRooms(ctx context.Context, req *emptypb.Empty) (*roompb.GetTotalRoomsResponse, error) {
	const op = "grpc.RoomService.GetTotalRooms"

	total, err := s.roomGetter.GetTotalRooms(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &roompb.GetTotalRoomsResponse{
		TotalRooms: int32(total),
	}, nil
}
