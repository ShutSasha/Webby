package grpc

import (
	"context"
	"fmt"
	"webby/room-queue-service/internal/grpc/mediapb"
	"webby/room-queue-service/pkg/logger"

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
	log := logger.FromContext(ctx).With("op", op)

	resp, err := m.client.GetVideo(ctx, &mediapb.GetVideoRequest{Id: id})
	if err != nil {
		log.Error("get video from media client", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("get video from media flient resp", "resp", resp)

	return &VideoInfo{
		ID:        resp.Id,
		Title:     resp.Title,
		Thumbnail: resp.Thumbnail,
		VideoUrl:  resp.VideoUrl,
	}, nil
}

func (m *MediaClient) GetVideosBatch(ctx context.Context, ids []string) ([]VideoInfo, error) {
	const op = "grpc.media_client.GetVideosBatch"
	log := logger.FromContext(ctx).With("op", op)

	resp, err := m.client.GetVideosBatch(ctx, &mediapb.GetVideosBatchRequest{Ids: ids})
	if err != nil {
		log.Error("get videos batch from media client", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("get videos batch from media flient resp: %v", "resp", resp)

	videoMap := make(map[string]*mediapb.VideoResponse, len(resp.GetVideos()))
	for _, v := range resp.GetVideos() {
		videoMap[v.Id] = v
	}

	unavailableMap := make(map[string]struct{}, len(resp.GetUnavailableVideoIds()))
	for _, id := range resp.GetUnavailableVideoIds() {
		unavailableMap[id] = struct{}{}
	}

	resultVideos := make([]VideoInfo, len(ids))
	for i, id := range ids {
		if _, unavailable := unavailableMap[id]; unavailable {
			resultVideos[i] = VideoInfo{ID: id, Title: "Unavailable video"}
			continue
		}

		if v, ok := videoMap[id]; ok {
			resultVideos[i] = VideoInfo{
				ID:        v.Id,
				Title:     v.Title,
				Thumbnail: v.Thumbnail,
				VideoUrl:  v.VideoUrl,
			}
		}
	}

	return resultVideos, nil
}
