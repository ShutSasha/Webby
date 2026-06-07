package grpc

import (
	"context"
	"fmt"
	"webby/chat-service/internal/grpc/userpb"
	"webby/chat-service/internal/models"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client userpb.UserGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(address string) (*UserClient, error) {
	const op = "grpc.NewUserClient"

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to member service: %w", op, err)
	}

	client := userpb.NewUserGrpcServiceClient(conn)

	return &UserClient{
		client: client,
		conn:   conn,
	}, nil
}

func (u *UserClient) Close() error {
	return u.conn.Close()
}

func (u *UserClient) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.Sender, error) {
	const op = "grpc.UserClient.GetUserByID"

	resp, err := u.client.GetUserById(
		ctx, &userpb.GetUserRequest{
			UserId: userID.String(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s, user %s: %w", op, userID.String(), err)
	}

	return &models.Sender{
		ID:        userID,
		Username:  resp.Username,
		AvatarURL: resp.AvatarUrl,
	}, nil
}
