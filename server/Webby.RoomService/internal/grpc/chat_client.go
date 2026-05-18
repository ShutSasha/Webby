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
		RoomId: roomID.String(),
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

	resp, err := c.client.GetChatByRoomId(ctx, &chatpb.GetChatByRoomIdRequest{
		RoomId: roomID.String(),
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

func (c *ChatClient) AddChatMember(ctx context.Context, chatID, userID uuid.UUID) error {
	const op = "grpc.ChatClient.AddChatMember"

	_, err := c.client.AddChatMember(ctx, &chatpb.AddChatMemberRequest{
		ChatId: chatID.String(),
		UserId: userID.String(),
	})
	if err != nil {
		return fmt.Errorf("%s (chat=%s, user=%s): %w", op, chatID, userID, err)
	}

	return nil
}
