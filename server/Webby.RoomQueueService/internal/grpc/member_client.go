package grpc

import (
	"context"
	"fmt"
	"webby/room-queue-service/internal/grpc/memberpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type memberClient struct {
	client memberpb.MemberGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewMemberClient(address string) (*memberClient, error) {
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

	return &memberClient{
		client: client,
		conn:   conn,
	}, nil
}

func (m *memberClient) Close() error {
	return m.conn.Close()
}

func (m *memberClient) Exists(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
	const op = "grpc.memberClient.Exists"

	resp, err := m.client.MemberExists(
		ctx, &memberpb.MemberExistsRequest{
			RoomId: roomID.String(),
			UserId: userID.String(),
		},
	)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetExists(), nil
}
