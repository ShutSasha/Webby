package services

import (
	"context"
	"fmt"
	"log/slog"
	"webby/room-queue-service/internal/apperrors"
	"webby/room-queue-service/internal/models"
	"webby/room-queue-service/pkg/logger"

	"github.com/google/uuid"
)

type queueUpdatedEvent struct {
	Positions []int `json:"positions"`
}

type eventEnvelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

const EventTypeQueueUpdated = "QUEUE_UPDATED"

type VideoInfo struct {
	ID        string
	Title     string
	Thumbnail string
	VideoUrl  string
}

type queueItemRepository interface {
	Create(ctx context.Context, item *models.QueueItem) (uuid.UUID, int, error)
	Delete(ctx context.Context, id uuid.UUID) (int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.QueueItem, error)
	ListByRoom(ctx context.Context, roomID uuid.UUID, offset, limit int) ([]models.QueueItem, int, error)
	MoveToTop(ctx context.Context, id uuid.UUID) error
	ActivateVideo(ctx context.Context, roomID, itemID uuid.UUID) (int, int, error)
}

type mediaRetriever interface {
	GetVideo(ctx context.Context, id string) (*VideoInfo, error)
	GetVideosBatch(ctx context.Context, ids []string) ([]VideoInfo, error)
}

type chatRetriever interface {
	GetChatIDByRoomID(ctx context.Context, roomID, userID uuid.UUID) (uuid.UUID, error)
}

type memberChecker interface {
	Exists(ctx context.Context, roomId, userId uuid.UUID) (bool, error)
}

type eventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type service struct {
	repository     queueItemRepository
	mediaRetriever mediaRetriever
	chatRetriever  chatRetriever
	memberChecker  memberChecker
	eventPublisher eventPublisher
}

func New(
	repo queueItemRepository,
	mediaRetriever mediaRetriever,
	chatRetriever chatRetriever,
	memberChecker memberChecker,
	eventPublisher eventPublisher,
) *service {
	return &service{
		repository:     repo,
		mediaRetriever: mediaRetriever,
		chatRetriever:  chatRetriever,
		memberChecker:  memberChecker,
		eventPublisher: eventPublisher,
	}
}

func (s *service) isMember(ctx context.Context, roomID, userID uuid.UUID) error {
	exists, err := s.memberChecker.Exists(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("check membership: %w", err)
	}
	if !exists {
		return apperrors.ErrNotRoomMember
	}
	return nil
}

func (s *service) AddToQueue(ctx context.Context, roomID, userID uuid.UUID, videoID string) (uuid.UUID, int, error) {
	const op = "services.AddToQueue"
	log := logger.FromContext(ctx).With(slog.String("op", op))

	if err := s.isMember(ctx, roomID, userID); err != nil {
		return uuid.Nil, -1, err
	}

	chatID, err := s.chatRetriever.GetChatIDByRoomID(ctx, roomID, userID)
	if err != nil {
		return uuid.Nil, -1, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := s.mediaRetriever.GetVideo(ctx, videoID); err != nil {
		return uuid.Nil, -1, fmt.Errorf("%s: %w", op, apperrors.ErrVideoNotFound)
	}

	item := &models.QueueItem{
		RoomID:  roomID,
		VideoID: videoID,
	}
	id, position, err := s.repository.Create(ctx, item)
	if err != nil {
		return uuid.Nil, -1, err
	}

	envelope := eventEnvelope{
		Type: EventTypeQueueUpdated,
		Payload: queueUpdatedEvent{
			Positions: []int{position},
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.eventPublisher.Publish(ctx, topic, envelope); err != nil {
		log.Error("failed to publish queue update", slog.String("err", err.Error()))
	}
	return id, position, nil
}

func (s *service) GetQueue(
	ctx context.Context,
	roomID, userID uuid.UUID,
	page, limit int,
) ([]models.EnrichedQueueItem, int, error) {
	const op = "services.GetQueue"

	if err := s.isMember(ctx, roomID, userID); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	offset := (page - 1) * limit
	items, total, err := s.repository.ListByRoom(ctx, roomID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	videoIDs := make([]string, 0, len(items))
	for _, item := range items {
		videoIDs = append(videoIDs, item.VideoID)
	}

	videos, err := s.mediaRetriever.GetVideosBatch(ctx, videoIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	videoMap := make(map[string]VideoInfo, len(videos))
	for _, v := range videos {
		videoMap[v.ID] = v
	}

	enrichedItems := make([]models.EnrichedQueueItem, 0, len(items))
	for _, item := range items {
		vidInfo := videoMap[item.VideoID]

		entry := models.EnrichedQueueItem{
			ID:        item.ID,
			VideoID:   item.VideoID,
			IsActive:  item.IsActive,
			Position:  item.Position,
			Title:     vidInfo.Title,
			Thumbnail: vidInfo.Thumbnail,
			VideoUrl:  vidInfo.VideoUrl,
		}
		enrichedItems = append(enrichedItems, entry)
	}

	return enrichedItems, total, nil
}

func (s *service) DeleteFromQueue(ctx context.Context, itemID, userID uuid.UUID) error {
	const op = "services.DeleteFromQueue"
	log := logger.FromContext(ctx).With(slog.String("op", op))

	item, err := s.repository.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed get item by ID: %w", op, err)
	}

	if err := s.isMember(ctx, item.RoomID, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	chatID, err := s.chatRetriever.GetChatIDByRoomID(ctx, item.RoomID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	position, err := s.repository.Delete(ctx, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item: %w", op, err)
	}

	envelope := eventEnvelope{
		Type: EventTypeQueueUpdated,
		Payload: queueUpdatedEvent{
			Positions: []int{position},
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.eventPublisher.Publish(ctx, topic, envelope); err != nil {
		log.Error("failed to publish queue update", slog.String("err", err.Error()))
	}

	return nil
}

// TODO
func (s *service) MoveToTop(ctx context.Context, id uuid.UUID) error {
	return s.repository.MoveToTop(ctx, id)
}

func (s *service) ActivateVideo(ctx context.Context, itemID, userID uuid.UUID) error {
	const op = "services.ActivateVideo"
	log := logger.FromContext(ctx).With(slog.String("op", op))

	item, err := s.repository.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed get queue item by ID: %w", op, err)
	}

	if err := s.isMember(ctx, item.RoomID, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	chatID, err := s.chatRetriever.GetChatIDByRoomID(ctx, item.RoomID, userID)
	if err != nil {
		return fmt.Errorf("%s: failed to get chat id: %w", op, err)
	}

	prevPos, currPos, err := s.repository.ActivateVideo(ctx, item.RoomID, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to activate video in db: %w", op, err)
	}

	if currPos != -1 {
		updatedPositions := make([]int, 0, 2)
		if prevPos != -1 {
			updatedPositions = append(updatedPositions, prevPos)
		}
		updatedPositions = append(updatedPositions, currPos)

		envelope := eventEnvelope{
			Type: EventTypeQueueUpdated,
			Payload: queueUpdatedEvent{
				Positions: updatedPositions,
			},
		}
		topic := fmt.Sprintf("chat:%s", chatID.String())
		if err := s.eventPublisher.Publish(ctx, topic, envelope); err != nil {
			log.Error("failed to publish queue update", slog.String("err", err.Error()))
		}
	}

	return nil
}
