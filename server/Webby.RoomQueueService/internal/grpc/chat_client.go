package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"webby/room-queue-service/internal/grpc/chatpb"
)

type ChatIDInfo struct {
	success bool
	chatID  uuid.UUID
}

type ChatClient struct {
	client chatpb.ChatServiceClient
	conn   *grpc.ClientConn
}

func NewChatClient(address string) (*ChatClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to chat service: %w", err,
		)
	}

	client := chatpb.NewChatServiceClient(conn)

	return &ChatClient{
		client: client,
		conn:   conn,
	}, nil
}

func (m *ChatClient) Close() error {
	return m.conn.Close()
}

func (m *ChatClient) GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (uuid.UUID, error) {
	const op = "grpc.chat_client.GetChatIDByRoomID"

	resp, err := m.client.GetChatIDByRoomID(ctx, &chatpb.GetChatIDByRoomIDRequest{
		RoomID: roomID.String(),
		UserID: userID.String(),
	})
	if err != nil && !resp.Success {
		return uuid.Nil, fmt.Errorf("%s: failed to get chat id by room id: %w", op, err)
	}

	return uuid.Parse(resp.GetChatID())
}
