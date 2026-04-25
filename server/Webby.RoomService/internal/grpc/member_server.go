package grpc

import (
	"context"
	"webby/internal/grpc/memberpb"

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

func (s *MemberServer) MemberExists(
	ctx context.Context, req *memberpb.MemberExistsRequest,
) (*memberpb.MemberExistsResponse, error) {
	roomId, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}

	exists, err := s.memberChecker.Exists(ctx, roomId, userId)
	if err != nil {
		return nil, err
	}

	return &memberpb.MemberExistsResponse{Exists: exists}, nil
}
