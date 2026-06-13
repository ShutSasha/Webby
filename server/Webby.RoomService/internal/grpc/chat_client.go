package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/chatpb"

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

func (c *chatClient) CreateChat(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	const op = "grpc.ChatClient.CreateChat"

	resp, err := c.client.CreateChat(ctx, &chatpb.CreateChatRequest{RoomID: roomID.String()})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	chatID, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s parse chat id: %w", op, err)
	}

	return chatID, nil
}

func (c *chatClient) GetChatByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	const op = "grpc.chatClient.GetChatByRoomID"

	resp, err := c.client.GetChatByRoomID(ctx, &chatpb.GetChatByRoomIdRequest{RoomID: roomID.String()})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	chatID, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s parse chat id: %w", op, err)
	}

	return chatID, nil
}

func (m *chatClient) GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	const op = "grpc.chatClient.GetChatIDByRoomID"

	resp, err := m.client.GetChatIDByRoomID(ctx, &chatpb.GetChatIDByRoomIDRequest{
		RoomID: roomID.String(),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return uuid.Parse(resp.GetChatID())
}

func (c *chatClient) AddChatMember(ctx context.Context, chatID, userID uuid.UUID) error {
	const op = "grpc.chatClient.AddChatMember"

	if _, err := c.client.AddChatMember(ctx, &chatpb.AddChatMemberRequest{
		ChatID: chatID.String(),
		UserID: userID.String(),
	}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
