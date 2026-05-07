package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/categorypb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CategoryClient struct {
	client categorypb.CategoryGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewCategoryClient(address string) (*CategoryClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to category service: %w", err)
	}

	client := categorypb.NewCategoryGrpcServiceClient(conn)

	return &CategoryClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *CategoryClient) Close() error {
	return c.conn.Close()
}

func (c *CategoryClient) Exists(ctx context.Context, name string) (bool, error) {
	resp, err := c.client.CategoryExists(ctx, &categorypb.CategoryExistsRequest{
		Name: name,
	})
	if err != nil {
		return false, fmt.Errorf("category exists check for %q: %w", name, err)
	}

	return resp.GetExists(), nil
}
