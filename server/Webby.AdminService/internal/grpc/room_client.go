package grpc

import (
	"context"
	"fmt"
	"webby/admin-service/internal/grpc/roompb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roomClient struct {
	client roompb.RoomGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewRoomClient(address string) (*roomClient, error) {
	const op = "grpc.NewRoomClient"

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to room service: %w", op, err)
	}

	client := roompb.NewRoomGrpcServiceClient(conn)

	return &roomClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *roomClient) Close() error {
	return c.conn.Close()
}

func (c *roomClient) GetTotalRooms(ctx context.Context) (int, error) {
	const op = "grpc.roomClient.GetTotalRooms"

	resp, err := c.client.GetTotalRooms(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return int(resp.GetTotalRooms()), nil
}
