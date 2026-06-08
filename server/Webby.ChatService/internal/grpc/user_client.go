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

type userClient struct {
	client userpb.UserGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(address string) (*userClient, error) {
	const op = "grpc.NewUserClient"

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to member service: %w", op, err)
	}

	client := userpb.NewUserGrpcServiceClient(conn)

	return &userClient{
		client: client,
		conn:   conn,
	}, nil
}

func (u *userClient) Close() error {
	return u.conn.Close()
}

func (u *userClient) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.Sender, error) {
	const op = "grpc.userClient.GetUserByID"

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

func (u *userClient) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]models.Sender, error) {
	const op = "grpc.userClient.GetUsersByIDs"

	userIDsStr := make([]string, len(userIDs))
	for i, userID := range userIDs {
		userIDsStr[i] = userID.String()
	}

	resp, err := u.client.GetUsersByIds(ctx, &userpb.GetUsersRequest{
		UserIds: userIDsStr,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	senders := make(map[uuid.UUID]models.Sender, len(resp.Users))
	for i, user := range resp.Users {
		userID := userIDs[i]
		if _, exists := senders[userID]; !exists {
			senders[userID] = models.Sender{
				ID:        userID,
				Username:  user.Username,
				AvatarURL: user.AvatarUrl,
			}
		}
	}

	return senders, nil
}
