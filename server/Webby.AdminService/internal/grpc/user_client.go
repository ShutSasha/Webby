package grpc

import (
	"context"
	"fmt"
	"webby/admin-service/internal/grpc/userpb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var months = map[int32]string{
	1: "Jan", 2: "Feb", 3: "Mar",
	4: "Apr", 5: "May", 6: "Jun",
	7: "Jul", 8: "Aug", 9: "Sep",
	10: "Oct", 11: "Nov", 12: "Dec",
}

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

func (c *userClient) GetMonthlyRegistrations(ctx context.Context) (map[string]int, error) {
	const op = "grpc.userClient.GetMonthlyRegistrations"

	stats, err := c.client.GetMonthlyRegistrations(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	registerStats := make(map[string]int)
	for month, registrations := range stats.GetRegistrations() {
		registerStats[months[month]] = int(registrations)
	}

	return registerStats, nil
}

func (c *userClient) GetMonthlySubscriptions(ctx context.Context) (map[string]int, error) {
	const op = "grpc.userClient.GetMonthlySubscriptions"

	stats, err := c.client.GetMonthlySubscriptions(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	subsStat := make(map[string]int)
	for month, subs := range stats.GetSubscriptions() {
		subsStat[months[month]] = int(subs)
	}

	return subsStat, nil
}

func (c *userClient) GetTotalRegistrations(ctx context.Context) (int, error) {
	const op = "grpc.userClient.GetTotalRegistrations"

	resp, err := c.client.GetTotalRegistrations(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return int(resp.GetTotalRegistrations()), nil
}

func (c *userClient) GetMonthRevenue(ctx context.Context) (int, error) {
	const op = "grpc.userClient.GetMonthRevenue"

	resp, err := c.client.GetMonthRevenue(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return int(resp.GetMonthRevenue()), nil
}
