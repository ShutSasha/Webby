package grpc

import (
	"context"
	"fmt"
	"webby/vote-service/internal/grpc/chatpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type chatClient struct {
	client chatpb.ChatGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewChatClient(address string) (*chatClient, error) {
	const op = "grpc.NewChatClient"

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to chat service: %w", op, err)
	}

	client := chatpb.NewChatGrpcServiceClient(conn)

	return &chatClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *chatClient) Close() error {
	return c.conn.Close()
}

func (m *chatClient) GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (uuid.UUID, error) {
	const op = "grpc.chatClient.GetChatIDByRoomID"

	resp, err := m.client.GetChatIDByRoomID(ctx, &chatpb.GetChatIDByRoomIDRequest{
		RoomID: roomID.String(),
		UserID: userID.String(),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return uuid.Parse(resp.GetChatID())
}
