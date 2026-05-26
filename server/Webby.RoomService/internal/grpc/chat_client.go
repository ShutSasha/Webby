package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/chatpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatClient struct {
	client chatpb.ChatGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewChatClient(address string) (*ChatClient, error) {
	const op = "grpc.NewChatClient"

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to chat service: %w", op, err)
	}

	client := chatpb.NewChatGrpcServiceClient(conn)

	return &ChatClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *ChatClient) Close() error {
	return c.conn.Close()
}

func (c *ChatClient) CreateChat(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	const op = "grpc.ChatClient.CreateChat"

	resp, err := c.client.CreateChat(ctx, &chatpb.CreateChatRequest{
		RoomID: roomID.String(),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s for room %s: %w", op, roomID, err)
	}

	chatID, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s parse chat id: %w", op, err)
	}

	return chatID, nil
}

func (c *ChatClient) GetChatByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	const op = "grpc.ChatClient.GetChatByRoomID"

	resp, err := c.client.GetChatByRoomID(ctx, &chatpb.GetChatByRoomIdRequest{
		RoomID: roomID.String(),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s for room %s: %w", op, roomID, err)
	}

	chatID, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s parse chat id: %w", op, err)
	}

	return chatID, nil
}

func (m *ChatClient) GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (uuid.UUID, error) {
	const op = "grpc.chat_client.GetChatIDByRoomID"

	resp, err := m.client.GetChatIDByRoomID(ctx, &chatpb.GetChatIDByRoomIDRequest{
		RoomID: roomID.String(),
		UserID: userID.String(),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: failed to get chat id by room id: %w", op, err)
	}

	return uuid.Parse(resp.GetChatID())
}

func (c *ChatClient) AddChatMember(ctx context.Context, chatID, userID uuid.UUID) error {
	const op = "grpc.ChatClient.AddChatMember"

	if _, err := c.client.AddChatMember(ctx, &chatpb.AddChatMemberRequest{
		ChatID: chatID.String(),
		UserID: userID.String(),
	}); err != nil {
		return fmt.Errorf("%s (chat=%s, user=%s): %w", op, chatID, userID, err)
	}

	return nil
}
