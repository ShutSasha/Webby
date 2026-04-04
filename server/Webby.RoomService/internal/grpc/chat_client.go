package grpc

import (
	"context"
	"fmt"
	"webby/internal/grpc/chatpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatClient struct {
	client chatpb.ChatGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewChatClient(address string) (*ChatClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chat service: %w", err)
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

func (c *ChatClient) CreateChat(ctx context.Context, roomId uuid.UUID) (uuid.UUID, error) {
	resp, err := c.client.CreateChat(ctx, &chatpb.CreateChatRequest{
		RoomId: roomId.String(),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create chat for room %s: %w", roomId, err)
	}

	chatId, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse chat id: %w", err)
	}

	return chatId, nil
}

func (c *ChatClient) GetChatByRoomId(ctx context.Context, roomId uuid.UUID) (uuid.UUID, error) {
	resp, err := c.client.GetChatByRoomId(ctx, &chatpb.GetChatByRoomIdRequest{
		RoomId: roomId.String(),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("get chat for room %s: %w", roomId, err)
	}

	chatId, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse chat id: %w", err)
	}

	return chatId, nil
}

func (c *ChatClient) AddChatMember(ctx context.Context, chatId, userId uuid.UUID) error {
	_, err := c.client.AddChatMember(ctx, &chatpb.AddChatMemberRequest{
		ChatId: chatId.String(),
		UserId: userId.String(),
	})
	if err != nil {
		return fmt.Errorf("add chat member (chat=%s, user=%s): %w", chatId, userId, err)
	}

	return nil
}
