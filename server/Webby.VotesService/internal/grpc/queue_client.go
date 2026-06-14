package grpc

import (
	"context"
	"fmt"
	"webby/vote-service/internal/grpc/queuepb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type queueClient struct {
	client queuepb.QueueGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewQueueClient(address string) (*queueClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to queue service: %w", err,
		)
	}

	client := queuepb.NewQueueGrpcServiceClient(conn)

	return &queueClient{
		client: client,
		conn:   conn,
	}, nil
}

func (q *queueClient) Close() error {
	return q.conn.Close()
}

func (q *queueClient) MakeNext(ctx context.Context, roomID, queueItemID uuid.UUID) error {
	const op = "queueClient.MakeNext"

	_, err := q.client.MakeNext(ctx, &queuepb.MakeNextRequest{
		RoomID:      roomID.String(),
		QueueItemID: queueItemID.String(),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
