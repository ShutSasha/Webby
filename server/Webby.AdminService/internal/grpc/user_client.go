package grpc

import (
	"context"
	"fmt"
	"webby/admin-service/internal/grpc/userpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userClient struct {
	client userpb.UserGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(address string) (*userClient, error) {
	const op = "grpc.NewUserClient"

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to user service: %w", op, err)
	}

	client := userpb.NewUserGrpcServiceClient(conn)

	return &userClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *userClient) Close() error {
	return c.conn.Close()
}

func (c *userClient) BanUser(ctx context.Context, userID uuid.UUID) error {
	const op = "grpc.userClient.BanUser"

	_, err := c.client.BanUser(ctx, &userpb.BanUserRequest{
		UserID: userID.String(),
	})
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}
