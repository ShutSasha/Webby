package grpc

import (
	"context"
	"fmt"
	"webby/room-queue-service/internal/grpc/mediapb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type VideoInfo struct {
	ID        string
	Title     string
	Thumbnail string
	VideoUrl  string
}

type MediaClient struct {
	client mediapb.MediaServiceClient
	conn   *grpc.ClientConn
}

func NewMediaClient(address string) (*MediaClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to media service: %w", err,
		)
	}

	client := mediapb.NewMediaServiceClient(conn)

	return &MediaClient{
		client: client,
		conn:   conn,
	}, nil
}

func (m *MediaClient) Close() error {
	return m.conn.Close()
}

func (m *MediaClient) GetVideo(ctx context.Context, id string) (*VideoInfo, error) {
	const op = "grpc.media_client.GetVideo"

	resp, err := m.client.GetVideo(ctx, &mediapb.GetVideoRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &VideoInfo{
		ID:        resp.Id,
		Title:     resp.Title,
		Thumbnail: resp.Thumbnail,
		VideoUrl:  resp.VideoUrl,
	}, nil
}

func (m *MediaClient) GetVideosBatch(ctx context.Context, ids []string) ([]VideoInfo, error) {
	const op = "grpc.media_client.GetVideosBatch"

	resp, err := m.client.GetVideosBatch(ctx, &mediapb.GetVideosBatchRequest{Ids: ids})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	videos := make([]VideoInfo, 0, len(resp.Videos))
	for _, video := range resp.GetVideos() {
		vidInfo := VideoInfo{
			ID:        video.Id,
			Title:     video.Title,
			Thumbnail: video.Thumbnail,
			VideoUrl:  video.VideoUrl,
		}

		videos = append(videos, vidInfo)
	}
	return videos, nil
}
