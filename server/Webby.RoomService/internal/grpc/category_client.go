package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/categorypb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type categoryClient struct {
	client categorypb.CategoryGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewCategoryClient(address string) (*categoryClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to category service: %w", err)
	}

	client := categorypb.NewCategoryGrpcServiceClient(conn)

	return &categoryClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *categoryClient) Close() error {
	return c.conn.Close()
}

func (c *categoryClient) Exists(ctx context.Context, name string) (bool, error) {
	const op = "grpc.categoryClient.Exists"

	resp, err := c.client.CategoryExists(ctx, &categorypb.CategoryExistsRequest{Name: name})
	if err != nil {
		return false, fmt.Errorf("%s for %q: %w", op, name, err)
	}

	return resp.GetExists(), nil
}
