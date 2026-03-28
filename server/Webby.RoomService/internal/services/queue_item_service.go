package services

import (
	"context"
	"fmt"
	"time"
	"webby/internal/apperrors"
	grpcClient "webby/internal/grpc"
	"webby/internal/models"

	"github.com/google/uuid"
)

type QueueItemRepository interface {
	Create(item *models.QueueItem) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	GetById(id uuid.UUID) (*models.QueueItem, error)
	ListByRoom(roomId uuid.UUID) ([]models.QueueItem, error)
}

type MediaClient interface {
	GetVideo(ctx context.Context, id uuid.UUID) (*grpcClient.VideoInfo, error)
	GetPlaylist(ctx context.Context, id uuid.UUID) (*grpcClient.PlaylistInfo, error)
}

type MemberChecker interface {
	Exists(roomId, userId uuid.UUID) (bool, error)
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
	Id         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Thumbnail  string    `json:"thumbnail"`
	VideoUrl   string    `json:"videoUrl"`
	PreviewUrl string    `json:"previewUrl"`
}

type QueueItemEnriched struct {
	Id         uuid.UUID         `json:"id"`
	EntityId   uuid.UUID         `json:"entityId"`
	EntityType string            `json:"entityType"`
	Title      string            `json:"title"`
	Thumbnail  string            `json:"thumbnail"`
	VideoUrl   string            `json:"videoUrl"`
	PreviewUrl string            `json:"previewUrl"`
	IsActive   bool              `json:"isActive"`
	IsFolder   bool              `json:"isFolder"`
	Children   []QueueVideoChild `json:"children,omitempty"`
}

func (s *QueueItemService) ensureMember(roomId, userId uuid.UUID) error {
	exists, err := s.memberChecker.Exists(roomId, userId)
	if err != nil {
		return fmt.Errorf("check membership: %w", err)
	}
	if !exists {
		return apperrors.ErrForbidden
	}
	return nil
}

func (s *QueueItemService) AddToQueue(ctx context.Context, roomId, userId, entityId uuid.UUID, entityType string) (*models.QueueItem, error) {
	if err := s.ensureMember(roomId, userId); err != nil {
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
		if _, err := s.mediaClient.GetPlaylist(ctx, entityId); err != nil {
			return nil, fmt.Errorf("%w: playlist not found", apperrors.ErrNotFound)
		}
	}

	item := &models.QueueItem{
		RoomId:     roomId,
		EntityId:   entityId,
		EntityType: entityType,
		IsActive:   false,
		CreatedAt:  time.Now(),
	}

	id, err := s.repo.Create(item)
	if err != nil {
		return nil, err
	}

	item.Id = id
	return item, nil
}

func (s *QueueItemService) GetQueue(ctx context.Context, roomId, userId uuid.UUID) ([]QueueItemEnriched, error) {
	if err := s.ensureMember(roomId, userId); err != nil {
		return nil, err
	}

	items, err := s.repo.ListByRoom(roomId)
	if err != nil {
		return nil, err
	}

	enriched := make([]QueueItemEnriched, 0, len(items))
	for _, item := range items {
		entry := QueueItemEnriched{
			Id:         item.Id,
			EntityId:   item.EntityId,
			EntityType: item.EntityType,
			IsActive:   item.IsActive,
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
				entry.PreviewUrl = video.PreviewUrl
			}
			entry.IsFolder = false
		case "playlist":
			playlist, err := s.mediaClient.GetPlaylist(ctx, item.EntityId)
			if err != nil {
				entry.Title = "Unknown playlist"
				entry.Thumbnail = ""
				entry.IsFolder = true
			} else {
				entry.Title = playlist.Title
				entry.Thumbnail = playlist.Thumbnail
				entry.IsFolder = true
				children := make([]QueueVideoChild, 0, len(playlist.Videos))
				for _, v := range playlist.Videos {
					children = append(children, QueueVideoChild{
						Id:         v.Id,
						Title:      v.Title,
						Thumbnail:  v.Thumbnail,
						VideoUrl:   v.VideoUrl,
						PreviewUrl: v.PreviewUrl,
					})
				}
				entry.Children = children
			}
		}

		enriched = append(enriched, entry)
	}

	return enriched, nil
}

func (s *QueueItemService) DeleteFromQueue(ctx context.Context, itemId, userId uuid.UUID) error {
	item, err := s.repo.GetById(itemId)
	if err != nil {
		return err
	}

	if err := s.ensureMember(item.RoomId, userId); err != nil {
		return err
	}

	return s.repo.Delete(itemId)
}
