package grpc

import (
	"context"
	"fmt"
	"webby/admin-service/internal/apperrors"
	"webby/admin-service/internal/grpc/videopb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type videoClient struct {
	client videopb.VideoGrpcServiceClient
	conn   *grpc.ClientConn
}

func NewVideoClient(address string) (*videoClient, error) {
	const op = "grpc.NewVideoClient"

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to video service: %w", op, err)
	}

	client := videopb.NewVideoGrpcServiceClient(conn)

	return &videoClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *videoClient) Close() error {
	return c.conn.Close()
}

func (c *videoClient) BanVideo(ctx context.Context, videoID uuid.UUID) error {
	const op = "grpc.videoClient.BanVideo"

	_, err := c.client.BanVideo(ctx, &videopb.BanVideoRequest{
		VideoID: videoID.String(),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return fmt.Errorf("%s: %w", op, apperrors.ErrVideoNotFound)
			case codes.InvalidArgument:
				return fmt.Errorf("%s: %w", op, apperrors.ErrVideoAlreadyBanned)
			}
		}

		return fmt.Errorf("%s %w", op, err)
	}

	return nil
}
