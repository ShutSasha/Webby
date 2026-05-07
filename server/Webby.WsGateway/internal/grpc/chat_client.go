package clients

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"webby/wsgateway/internal/grpc/chatpb"
)

type ChatClient struct {
	conn *grpc.ClientConn
	c    chatpb.ChatGrpcServiceClient
}

func NewChatClient(addr string) (*ChatClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ChatClient{conn: conn, c: chatpb.NewChatGrpcServiceClient(conn)}, nil
}

func (c *ChatClient) Close() error { return c.conn.Close() }

func (c *ChatClient) SaveMessage(ctx context.Context, chatID, userID, content string) (string, error) {
	resp, err := c.c.SaveMessage(ctx, &chatpb.SaveMessageRequest{
		ChatId:  chatID,
		UserId:  userID,
		Content: content,
	})
	if err != nil {
		return "", err
	}
	return resp.GetMessageId(), nil
}
