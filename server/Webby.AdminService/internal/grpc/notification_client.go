package grpc

import (
	"context"
	"fmt"
	"webby/admin-service/internal/grpc/notificationpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var complaintTypes = map[string]notificationpb.GrpcNotificationTargetType{
	"User":  notificationpb.GrpcNotificationTargetType_GRPC_NOTIFICATION_TARGET_TYPE_USER,
	"Video": notificationpb.GrpcNotificationTargetType_GRPC_NOTIFICATION_TARGET_TYPE_VIDEO,
}

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

func (c *notificationClient) SendNotificationToUser(ctx context.Context, userID, targetID uuid.UUID, title, message string, complaintType string) error {
	const op = "grpc.notificationClient.SendNotificationToUser"

	_, err := c.client.SendNotificationToUser(ctx, &notificationpb.CreateNotificationRequest{
		UserId:           userID.String(),
		Title:            title,
		Message:          message,
		TargetType:       complaintTypes[complaintType],
		TargetIdentifier: targetID.String(),
	})
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}
