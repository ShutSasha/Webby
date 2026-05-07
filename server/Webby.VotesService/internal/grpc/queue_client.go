package grpc

import (
	"context"
	"fmt"
	"webby/vote-service/internal/grpc/queuepb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QueueClient struct {
	client queuepb.QueueGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewQueueClient(address string) (*QueueClient, error) {
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

	return &QueueClient{
		client: client,
		conn:   conn,
	}, nil
}

func (q *QueueClient) Close() error {
	return q.conn.Close()
}

func (q *QueueClient) MoveToTop(
	ctx context.Context, id uuid.UUID,
) error {
	_, err := q.client.MoveToTop(
		ctx, &queuepb.MoveToTopRequest{
			Id: id.String(),
		},
	)
	if err != nil {
		return fmt.Errorf("move to top %s: %w", id.String(), err)
	}

	return nil
}
