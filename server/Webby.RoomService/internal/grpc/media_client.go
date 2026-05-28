package grpc

import (
	"context"
	"fmt"
	"webby/room-service/internal/grpc/mediapb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type VideoInfo struct {
	ID        uuid.UUID
	Title     string
	Thumbnail string
	VideoUrl  string
}

type PlaylistInfo struct {
	ID         uuid.UUID
	Title      string
	Thumbnail  string
	TotalCount int
	Videos     []VideoInfo
}

type MediaClient struct {
	client mediapb.MediaServiceClient
	conn   *grpc.ClientConn
}

func NewMediaClient(address string) (*MediaClient, error) {
	const op = "grpc.NewMediaClient"

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to media service: %w", op, err)
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

func (m *MediaClient) GetVideo(ctx context.Context, id uuid.UUID) (*VideoInfo, error) {
	const op = "grpc.MediaClient.GetVideo"

	resp, err := m.client.GetVideo(ctx, &mediapb.GetVideoRequest{Id: id.String()})
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", op, id.String(), err)
	}

	videoId, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, fmt.Errorf("%s parse video id: %w", op, err)
	}

	return &VideoInfo{
		ID:        videoId,
		Title:     resp.Title,
		Thumbnail: resp.Thumbnail,
		VideoUrl:  resp.VideoUrl,
	}, nil
}
