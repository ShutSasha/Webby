package grpc

import (
	"context"
	"fmt"
	"webby-room-queue/internal/grpc/mediapb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type VideoInfo struct {
	Id        uuid.UUID
	Title     string
	Thumbnail string
	VideoUrl  string
}

type PlaylistInfo struct {
	Id         uuid.UUID
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

func entityTypeToVideoType(entityType string) mediapb.VideoType {
	switch entityType {
	case "youtube":
		return mediapb.VideoType_YOUTUBE
	case "twitch":
		return mediapb.VideoType_TWITCH
	default:
		return mediapb.VideoType_VIDEO
	}
}

func (m *MediaClient) GetVideo(
	ctx context.Context, id uuid.UUID, entityType string,
) (*VideoInfo, error) {
	resp, err := m.client.GetVideo(
		ctx, &mediapb.GetVideoRequest{
			Id:        id.String(),
			VideoType: entityTypeToVideoType(entityType),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get video %s: %w", id.String(), err)
	}

	videoId, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, fmt.Errorf("parse video id: %w", err)
	}

	return &VideoInfo{
		Id:        videoId,
		Title:     resp.Title,
		Thumbnail: resp.Thumbnail,
		VideoUrl:  resp.VideoUrl,
	}, nil
}

func (m *MediaClient) GetPlaylist(
	ctx context.Context, id uuid.UUID, page, pageSize int32,
) (*PlaylistInfo, error) {
	resp, err := m.client.GetPlaylist(ctx, &mediapb.GetPlaylistRequest{
		Id:       id.String(),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("get playlist %s: %w", id.String(), err)
	}

	playlistId, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, fmt.Errorf("parse playlist id: %w", err)
	}

	videos := make([]VideoInfo, 0, len(resp.Videos))
	for _, v := range resp.Videos {
		vid, err := uuid.Parse(v.Id)
		if err != nil {
			return nil, fmt.Errorf("parse video id in playlist: %w", err)
		}
		videos = append(videos, VideoInfo{
			Id:        vid,
			Title:     v.Title,
			Thumbnail: v.Thumbnail,
			VideoUrl:  v.VideoUrl,
		})
	}

	return &PlaylistInfo{
		Id:         playlistId,
		Title:      resp.Title,
		Thumbnail:  resp.Thumbnail,
		TotalCount: int(resp.TotalCount),
		Videos:     videos,
	}, nil
}
