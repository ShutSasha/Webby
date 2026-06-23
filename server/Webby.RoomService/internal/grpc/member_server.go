package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/memberpb"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type memberChecker interface {
	GetMemberStatus(ctx context.Context, roomId, userId uuid.UUID) (*models.MemberStatus, error)
}

type memberServer struct {
	memberpb.UnimplementedMemberGrpcServiceServer
	memberChecker memberChecker
}

func NewMemberServer(memberChecker memberChecker) *memberServer {
	return &memberServer{memberChecker: memberChecker}
}

func (s *memberServer) MemberExists(ctx context.Context, req *memberpb.MemberExistsRequest) (*memberpb.MemberExistsResponse, error) {
	const op = "grpc.memberServer.MemberExists"

	roomID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	status, err := s.memberChecker.GetMemberStatus(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &memberpb.MemberExistsResponse{Exists: status.IsMember}, nil
}
