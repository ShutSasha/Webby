package services

import (
	"context"
	"fmt"
	"webby/room-queue-service/internal/apperrors"
	grpcClient "webby/room-queue-service/internal/grpc"
	"webby/room-queue-service/internal/models"

	"github.com/google/uuid"
)

type QueueItemRepository interface {
	Create(ctx context.Context, item *models.QueueItem) (uuid.UUID, int, error)
	Delete(ctx context.Context, id uuid.UUID) (int, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.QueueItem, error)
	ListByRoom(ctx context.Context, roomId uuid.UUID, offset, limit int) ([]models.QueueItem, int, error)
	MoveToTop(ctx context.Context, id uuid.UUID) error
}

type MediaClient interface {
	GetVideo(ctx context.Context, id string) (*grpcClient.VideoInfo, error)
	GetVideosBatch(ctx context.Context, ids []string) ([]grpcClient.VideoInfo, error)
}

type ChatClient interface {
	GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (uuid.UUID, error)
}

type MemberChecker interface {
	Exists(ctx context.Context, roomId, userId uuid.UUID) (bool, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type Service struct {
	repo          QueueItemRepository
	mediaClient   MediaClient
	chatClient    ChatClient
	memberChecker MemberChecker
	publisher     EventPublisher
}

func New(
	repo QueueItemRepository,
	mediaClient MediaClient,
	chatClient ChatClient,
	memberChecker MemberChecker,
	publisher EventPublisher,
) *Service {
	return &Service{
		repo:          repo,
		mediaClient:   mediaClient,
		chatClient:    chatClient,
		memberChecker: memberChecker,
		publisher:     publisher,
	}
}

type QueueItemEnriched struct {
	Id        uuid.UUID `json:"id"`
	VideoID   uuid.UUID `json:"videoId"`
	VideoType string    `json:"videoType"`
	Title     string    `json:"title"`
	Thumbnail string    `json:"thumbnail"`
	VideoUrl  string    `json:"videoUrl"`
	IsActive  bool      `json:"isActive"`
	Position  int       `json:"position"`
}

func (s *Service) ensureMember(
	ctx context.Context, roomId, userId uuid.UUID,
) error {
	exists, err := s.memberChecker.Exists(ctx, roomId, userId)
	if err != nil {
		return fmt.Errorf("check membership: %w", err)
	}
	if !exists {
		return apperrors.ErrForbidden
	}
	return nil
}

func (s *Service) AddToQueue(
	ctx context.Context,
	roomID, userID uuid.UUID,
	videoID string,
) (uuid.UUID, int, error) {
	const op = "services.AddToQueue"

	if err := s.ensureMember(ctx, roomID, userID); err != nil {
		return uuid.Nil, -1, err
	}

	if _, err := s.mediaClient.GetVideo(ctx, videoID); err != nil {
		return uuid.Nil, -1, fmt.Errorf("%s: %w", op, apperrors.ErrVideoNotFound)
	}

	item := &models.QueueItem{
		RoomID:  roomID,
		VideoID: videoID,
	}
	id, position, err := s.repo.Create(ctx, item)
	if err != nil {
		return uuid.Nil, -1, err
	}

	chatID, err := s.chatClient.GetChatIDByRoomID(ctx, roomID, userID)
	if err != nil {
		return uuid.Nil, -1, fmt.Errorf("%s: %w", op, err)
	}

	envelope := struct {
		Type    string `json:"type"`
		Payload any    `json:"payload"`
	}{
		Type:    "QUEUE_UPDATED",
		Payload: map[string]any{"position": position},
	}
	if err := s.publisher.Publish(ctx, "chat:"+chatID.String(), envelope); err != nil {
		return uuid.Nil, -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, position, nil
}

func (s *Service) GetQueue(
	ctx context.Context,
	roomID, userID uuid.UUID,
	page, limit int,
) ([]models.EnrichedQueueItem, int, error) {
	const op = "services.GetQueue"

	if err := s.ensureMember(ctx, roomID, userID); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	offset := (page - 1) * limit
	items, total, err := s.repo.ListByRoom(ctx, roomID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	videoIDs := make([]string, 0, len(items))
	for _, item := range items {
		videoIDs = append(videoIDs, item.VideoID)
	}

	videos, err := s.mediaClient.GetVideosBatch(ctx, videoIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	enrichedItems := make([]models.EnrichedQueueItem, 0, len(items))
	for i := range items {
		entry := models.EnrichedQueueItem{
			ID:        items[i].ID,
			VideoID:   items[i].VideoID,
			IsActive:  items[i].IsActive,
			Position:  items[i].Position,
			Title:     videos[i].Title,
			Thumbnail: videos[i].Thumbnail,
			VideoUrl:  videos[i].VideoUrl,
		}
		enrichedItems = append(enrichedItems, entry)
	}

	return enrichedItems, total, nil
}

func (s *Service) DeleteFromQueue(ctx context.Context, itemID, userID uuid.UUID) error {
	const op = "services.DeleteFromQueue"

	item, err := s.repo.GetById(ctx, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed get item by ID: %w", op, err)
	}

	if err := s.ensureMember(ctx, item.RoomID, userID); err != nil {
		return fmt.Errorf("%s: failed ensure user is room member: %w", op, err)
	}

	position, err := s.repo.Delete(ctx, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item: %w", op, err)
	}

	chatID, err := s.chatClient.GetChatIDByRoomID(ctx, item.RoomID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	envelope := struct {
		Type    string `json:"type"`
		Payload any    `json:"payload"`
	}{
		Type:    "QUEUE_UPDATED",
		Payload: map[string]any{"position": position},
	}
	if err := s.publisher.Publish(ctx, "chat:"+chatID.String(), envelope); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) MoveToTop(ctx context.Context, id uuid.UUID) error {
	return s.repo.MoveToTop(ctx, id)
}
