package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/memberpb"

	"github.com/google/uuid"
)

type MemberChecker interface {
	Exists(ctx context.Context, roomId, userId uuid.UUID) (bool, error)
}

type MemberServer struct {
	memberpb.UnimplementedMemberGrpcServiceServer
	memberChecker MemberChecker
}

func NewMemberServer(memberChecker MemberChecker) *MemberServer {
	return &MemberServer{memberChecker: memberChecker}
}

func (s *MemberServer) MemberExists(ctx context.Context, req *memberpb.MemberExistsRequest) (*memberpb.MemberExistsResponse, error) {
	const op = "grpc.MemberServer.MemberExists"

	roomID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	exists, err := s.memberChecker.Exists(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &memberpb.MemberExistsResponse{Exists: exists}, nil
}
