package grpc

import (
	"context"
	"fmt"
	"webby/vote-service/internal/grpc/memberpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MemberClient struct {
	client memberpb.MemberGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewMemberClient(address string) (*MemberClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to member service: %w", err,
		)
	}

	client := memberpb.NewMemberGrpcServiceClient(conn)

	return &MemberClient{
		client: client,
		conn:   conn,
	}, nil
}

func (m *MemberClient) Close() error {
	return m.conn.Close()
}

func (m *MemberClient) Exists(
	ctx context.Context, roomId, userId uuid.UUID,
) (bool, error) {
	resp, err := m.client.MemberExists(
		ctx, &memberpb.MemberExistsRequest{
			RoomId: roomId.String(),
			UserId: userId.String(),
		},
	)
	if err != nil {
		return false, fmt.Errorf(
			"member exists check for room %s, user %s: %w",
			roomId.String(), userId.String(), err,
		)
	}

	return resp.GetExists(), nil
}
