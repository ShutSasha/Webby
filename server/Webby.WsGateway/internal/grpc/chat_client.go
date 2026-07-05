package grpc

import (
	"context"
	"fmt"

	"webby/wsgateway/internal/grpc/chatpb"

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

func (c *chatClient) IsChatRelatedToRoom(ctx context.Context, chatID uuid.UUID) (bool, error) {
	const op = "grpc.chatClient.IsChatRelatedToRoom"

	resp, err := c.client.IsChatRelatedToRoom(ctx, &chatpb.IsChatRelatedToRoomRequest{
		ChatID: chatID.String(),
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.IsRelated, nil
}
