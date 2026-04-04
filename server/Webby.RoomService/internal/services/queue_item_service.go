package services

import (
	"context"
	"fmt"
	"webby/internal/apperrors"
	grpcClient "webby/internal/grpc"
	"webby/internal/models"

	"github.com/google/uuid"
)

type QueueItemRepository interface {
	Create(ctx context.Context, item *models.QueueItem) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*models.QueueItem, error)
	ListByRoom(ctx context.Context, roomId uuid.UUID) ([]models.QueueItem, error)
}

type MediaClient interface {
	GetVideo(ctx context.Context, id uuid.UUID) (*grpcClient.VideoInfo, error)
	GetPlaylist(ctx context.Context, id uuid.UUID, page, pageSize int32) (*grpcClient.PlaylistInfo, error)
}

type MemberChecker interface {
	Exists(ctx context.Context, roomId, userId uuid.UUID) (bool, error)
}

type QueueItemService struct {
	repo          QueueItemRepository
	mediaClient   MediaClient
	memberChecker MemberChecker
}

func NewQueueItemService(repo QueueItemRepository, mediaClient MediaClient, memberChecker MemberChecker) *QueueItemService {
	return &QueueItemService{
		repo:          repo,
		mediaClient:   mediaClient,
		memberChecker: memberChecker,
	}
}

type QueueVideoChild struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Thumbnail string    `json:"thumbnail"`
	VideoUrl  string    `json:"videoUrl"`
}

type QueueItemEnriched struct {
	Id            uuid.UUID         `json:"id"`
	EntityId      uuid.UUID         `json:"entityId"`
	EntityType    string            `json:"entityType"`
	Title         string            `json:"title"`
	Thumbnail     string            `json:"thumbnail"`
	VideoUrl      string            `json:"videoUrl"`
	IsActive      bool              `json:"isActive"`
	IsFolder      bool              `json:"isFolder"`
	Position      int               `json:"position"`
	TotalChildren int               `json:"totalChildren"`
	Children      []QueueVideoChild `json:"children,omitempty"`
}

func (s *QueueItemService) ensureMember(ctx context.Context, roomId, userId uuid.UUID) error {
	exists, err := s.memberChecker.Exists(ctx, roomId, userId)
	if err != nil {
		return fmt.Errorf("check membership: %w", err)
	}
	if !exists {
		return apperrors.ErrForbidden
	}
	return nil
}

func (s *QueueItemService) AddToQueue(ctx context.Context, roomId, userId, entityId uuid.UUID, entityType string) (*models.QueueItem, error) {
	if err := s.ensureMember(ctx, roomId, userId); err != nil {
		return nil, err
	}

	if entityType != "video" && entityType != "playlist" {
		return nil, fmt.Errorf("%w: entityType must be 'video' or 'playlist'", apperrors.ErrInvalidInput)
	}

	switch entityType {
	case "video":
		if _, err := s.mediaClient.GetVideo(ctx, entityId); err != nil {
			return nil, fmt.Errorf("%w: video not found", apperrors.ErrNotFound)
		}
	case "playlist":
		if _, err := s.mediaClient.GetPlaylist(ctx, entityId, 1, 1); err != nil {
			return nil, fmt.Errorf("%w: playlist not found", apperrors.ErrNotFound)
		}
	}

	item := &models.QueueItem{
		RoomId:     roomId,
		EntityId:   entityId,
		EntityType: entityType,
		IsActive:   false,
	}

	id, err := s.repo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	item.Id = id
	return item, nil
}

func (s *QueueItemService) GetQueue(ctx context.Context, roomId, userId uuid.UUID, page, limit int) ([]QueueItemEnriched, int, error) {
	if err := s.ensureMember(ctx, roomId, userId); err != nil {
		return nil, 0, err
	}

	items, err := s.repo.ListByRoom(ctx, roomId)
	if err != nil {
		return nil, 0, err
	}

	type itemMeta struct {
		videoCount int
	}
	metas := make([]itemMeta, len(items))
	totalVideos := 0

	for i, item := range items {
		switch item.EntityType {
		case "video":
			metas[i] = itemMeta{videoCount: 1}
		case "playlist":
			pl, err := s.mediaClient.GetPlaylist(ctx, item.EntityId, 1, 1)
			if err != nil {
				metas[i] = itemMeta{videoCount: 0}
			} else {
				metas[i] = itemMeta{videoCount: pl.TotalCount}
			}
		}
		totalVideos += metas[i].videoCount
	}

	offset := (page - 1) * limit
	running := 0
	enriched := make([]QueueItemEnriched, 0)
	remaining := limit

	for i, item := range items {
		vc := metas[i].videoCount
		if vc == 0 {
			running += vc
			continue
		}

		itemEnd := running + vc
		if itemEnd <= offset {
			running = itemEnd
			continue
		}
		if remaining <= 0 {
			break
		}

		entry := QueueItemEnriched{
			Id:         item.Id,
			EntityId:   item.EntityId,
			EntityType: item.EntityType,
			IsActive:   item.IsActive,
			Position:   item.Position,
		}

		switch item.EntityType {
		case "video":
			video, err := s.mediaClient.GetVideo(ctx, item.EntityId)
			if err != nil {
				entry.Title = "Unknown video"
				entry.Thumbnail = ""
			} else {
				entry.Title = video.Title
				entry.Thumbnail = video.Thumbnail
				entry.VideoUrl = video.VideoUrl
			}
			entry.IsFolder = false
			remaining--

		case "playlist":
			subStart := 0
			if offset > running {
				subStart = offset - running
			}
			subCount := vc - subStart
			if subCount > remaining {
				subCount = remaining
			}

			playlist, err := s.mediaClient.GetPlaylist(ctx, item.EntityId, 1, int32(subStart+subCount))
			if err != nil {
				entry.Title = "Unknown playlist"
				entry.Thumbnail = ""
				entry.IsFolder = true
				entry.TotalChildren = vc
			} else {
				entry.Title = playlist.Title
				entry.Thumbnail = playlist.Thumbnail
				entry.IsFolder = true
				entry.TotalChildren = playlist.TotalCount
				vids := playlist.Videos
				if subStart < len(vids) {
					vids = vids[subStart:]
				}
				children := make([]QueueVideoChild, 0, len(vids))
				for _, v := range vids {
					children = append(children, QueueVideoChild{
						Id:        v.Id,
						Title:     v.Title,
						Thumbnail: v.Thumbnail,
						VideoUrl:  v.VideoUrl,
					})
				}
				entry.Children = children
			}
			remaining -= subCount
		}

		enriched = append(enriched, entry)
		running = itemEnd
	}

	return enriched, totalVideos, nil
}

func (s *QueueItemService) DeleteFromQueue(ctx context.Context, itemId, userId uuid.UUID) error {
	item, err := s.repo.GetById(ctx, itemId)
	if err != nil {
		return err
	}

	if err := s.ensureMember(ctx, item.RoomId, userId); err != nil {
		return err
	}

	return s.repo.Delete(ctx, itemId)
}
