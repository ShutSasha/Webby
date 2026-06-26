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

func (m *mediaClient) GetVideosBatch(ctx context.Context, ids []string) (map[string]string, error) {
	const op = "grpc.mediaClient.GetVideosBatch"

	if len(ids) == 0 {
		return nil, nil
	}

	resp, err := m.client.GetVideosBatch(ctx, &mediapb.GetVideosBatchRequest{Ids: ids})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	videoMap := make(map[string]*mediapb.VideoResponse, len(resp.GetVideos()))
	for _, v := range resp.GetVideos() {
		videoMap[v.Id] = v
	}

	unavailableMap := make(map[string]struct{}, len(resp.GetUnavailableVideoIds()))
	for _, id := range resp.GetUnavailableVideoIds() {
		unavailableMap[id] = struct{}{}
	}

	resultVideos := make(map[string]string, len(ids))
	for _, id := range ids {
		if _, unavailable := unavailableMap[id]; unavailable {
			resultVideos[id] = "Unavailable video"
			continue
		}

		if v, ok := videoMap[id]; ok {
			resultVideos[id] = v.Title
		}
	}

	return resultVideos, nil
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
