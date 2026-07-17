package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"webby/room-queue-service/internal/grpc/chatpb"
)

type chatIDInfo struct {
	success bool
	chatID  uuid.UUID
}

type chatClient struct {
	client chatpb.ChatGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewChatClient(address string) (*chatClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to chat service: %w", err,
		)
	}

	client := chatpb.NewChatGrpcServiceClient(conn)

	return &chatClient{
		client: client,
		conn:   conn,
	}, nil
}

func (m *chatClient) Close() error {
	return m.conn.Close()
}

func (m *chatClient) GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	const op = "grpc.chatClient.GetChatIDByRoomID"

	resp, err := m.client.GetChatIDByRoomID(ctx, &chatpb.GetChatIDByRoomIDRequest{
		RoomID: roomID.String(),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: failed to get chat id by room id: %w", op, err)
	}

	return uuid.Parse(resp.GetChatID())
}
