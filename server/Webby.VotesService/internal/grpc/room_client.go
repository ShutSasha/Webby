package grpc

import (
	"context"
	"fmt"
	"webby/vote-service/internal/grpc/roompb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RoomClient struct {
	client roompb.RoomGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewRoomClient(address string) (*RoomClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to room service: %w", err,
		)
	}

	client := roompb.NewRoomGrpcServiceClient(conn)

	return &RoomClient{
		client: client,
		conn:   conn,
	}, nil
}

func (r *RoomClient) Close() error {
	return r.conn.Close()
}

func (r *RoomClient) GetRoomHost(
	ctx context.Context, roomId uuid.UUID,
) (uuid.UUID, error) {
	resp, err := r.client.GetRoomHost(
		ctx, &roompb.GetRoomHostRequest{
			RoomId: roomId.String(),
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"get room host for room %s: %w",
			roomId.String(), err,
		)
	}

	hostId, err := uuid.Parse(resp.GetHostId())
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"parse host id: %w", err,
		)
	}

	return hostId, nil
}
