package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/notificationpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type notificationClient struct {
	client notificationpb.NotificationGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewNotificationClient(address string) (*notificationClient, error) {
	const op = "grpc.NewNotificationClient"

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to notification service: %w", op, err)
	}

	client := notificationpb.NewNotificationGrpcServiceClient(conn)

	return &notificationClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *notificationClient) Close() error {
	return c.conn.Close()
}

func (c *notificationClient) SendNotificationToUser(ctx context.Context, userID, roomID uuid.UUID) error {
	const op = "grpc.NotificationClient.SendNotificationToUser"

	_, err := c.client.SendNotificationToUser(ctx, &notificationpb.CreateNotificationRequest{
		UserId:           userID.String(),
		Title:            "Room invitation",
		Message:          "You were invited to the room",
		TargetType:       notificationpb.GrpcNotificationTargetType_GRPC_NOTIFICATION_TARGET_TYPE_ROOM,
		TargetIdentifier: roomID.String(),
	})
	if err != nil {
		return fmt.Errorf("%s for room %s: %w", op, roomID, err)
	}

	return nil
}
