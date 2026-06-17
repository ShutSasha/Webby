package grpc

import (
	"context"
	"fmt"
	"webby/admin-service/internal/grpc/mediapb"
	"webby/admin-service/internal/models"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type mediaClient struct {
	client mediapb.MediaServiceClient
	conn   *grpc.ClientConn
}

func NewMediaClient(address string) (*mediaClient, error) {
	const op = "grpc.NewMediaClient"

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to complaint service: %w", op, err)
	}

	client := mediapb.NewMediaServiceClient(conn)

	return &mediaClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *mediaClient) Close() error {
	return c.conn.Close()
}

func (c *mediaClient) GetVideoByID(ctx context.Context, videoID uuid.UUID) (*models.Video, error) {
	const op = "grpc.mediaClient.GetVideoByID"

	video, err := c.client.GetVideo(ctx, &mediapb.GetVideoRequest{
		Id: "wb_" + videoID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &models.Video{
		ID:    video.Id,
		Title: video.Title,
	}, nil
}

func (c *mediaClient) GetAuthorIDByVideoID(ctx context.Context, videoID string) (uuid.UUID, error) {
	const op = "grpc.mediaClient.GetAuthorIDByVideoID"

	resp, err := c.client.GetAuthorIDByVideoID(ctx, &mediapb.GetAuthorIDByVideoIDRequest{
		VideoID: videoID,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s %w", op, err)
	}

	authorID, err := uuid.Parse(resp.GetAuthorID())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s %w", op, err)
	}

	return authorID, nil
}
