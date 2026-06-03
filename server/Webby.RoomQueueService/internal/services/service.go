package services

import (
	"context"
	"fmt"
	"log/slog"
	"webby/room-queue-service/internal/apperrors"
	grpcClient "webby/room-queue-service/internal/grpc"
	"webby/room-queue-service/internal/models"
	"webby/room-queue-service/pkg/logger"

	"github.com/google/uuid"
)

type QueueUpdatedEvent struct {
	Positions []int `json:"positions"`
}

type EventEnvelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

const EventTypeQueueUpdated = "QUEUE_UPDATED"

type QueueItemRepository interface {
	Create(ctx context.Context, item *models.QueueItem) (uuid.UUID, int, error)
	Delete(ctx context.Context, id uuid.UUID) (int, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.QueueItem, error)
	ListByRoom(
		ctx context.Context,
		roomId uuid.UUID,
		offset, limit int,
	) ([]models.QueueItem, int, error)
	MoveToTop(ctx context.Context, id uuid.UUID) error
	ActivateVideo(ctx context.Context, roomID, itemID uuid.UUID) (int, int, error)
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

func (s *Service) isMember(
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
	log := logger.FromContext(ctx).With(slog.String("op", op))

	if err := s.isMember(ctx, roomID, userID); err != nil {
		return uuid.Nil, -1, err
	}

	chatID, err := s.chatClient.GetChatIDByRoomID(ctx, roomID, userID)
	if err != nil {
		return uuid.Nil, -1, fmt.Errorf("%s: %w", op, err)
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

	envelope := EventEnvelope{
		Type: EventTypeQueueUpdated,
		Payload: QueueUpdatedEvent{
			Positions: []int{position},
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.publisher.Publish(ctx, topic, envelope); err != nil {
		log.Error("failed to publish queue update", slog.Any("err", err))
	}
	return id, position, nil
}

func (s *Service) GetQueue(
	ctx context.Context,
	roomID, userID uuid.UUID,
	page, limit int,
) ([]models.EnrichedQueueItem, int, error) {
	const op = "services.GetQueue"

	if err := s.isMember(ctx, roomID, userID); err != nil {
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

	videoMap := make(map[string]grpcClient.VideoInfo, len(videos))
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

func (s *Service) DeleteFromQueue(ctx context.Context, itemID, userID uuid.UUID) error {
	const op = "services.DeleteFromQueue"
	log := logger.FromContext(ctx).With(slog.String("op", op))

	item, err := s.repo.GetById(ctx, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed get item by ID: %w", op, err)
	}

	if err := s.isMember(ctx, item.RoomID, userID); err != nil {
		return fmt.Errorf("%s: failed ensure user is room member: %w", op, err)
	}

	chatID, err := s.chatClient.GetChatIDByRoomID(ctx, item.RoomID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	position, err := s.repo.Delete(ctx, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item: %w", op, err)
	}

	envelope := EventEnvelope{
		Type: EventTypeQueueUpdated,
		Payload: QueueUpdatedEvent{
			Positions: []int{position},
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.publisher.Publish(ctx, topic, envelope); err != nil {
		log.Error("failed to publish queue update", slog.Any("err", err))
	}

	return nil
}

// TODO
func (s *Service) MoveToTop(ctx context.Context, id uuid.UUID) error {
	return s.repo.MoveToTop(ctx, id)
}

func (s *Service) ActivateVideo(ctx context.Context, itemID, userID uuid.UUID) error {
	const op = "services.ActivateVideo"
	log := logger.FromContext(ctx).With(slog.String("op", op))

	item, err := s.repo.GetById(ctx, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed get queue item by ID: %w", op, err)
	}

	if err := s.isMember(ctx, item.RoomID, userID); err != nil {
		return fmt.Errorf("%s: failed ensure user is room member: %w", op, err)
	}

	chatID, err := s.chatClient.GetChatIDByRoomID(ctx, item.RoomID, userID)
	if err != nil {
		return fmt.Errorf("%s: failed to get chat id: %w", op, err)
	}

	prevPos, currPos, err := s.repo.ActivateVideo(ctx, item.RoomID, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to activate video in db: %w", op, err)
	}

	if currPos != -1 {
		updatedPositions := make([]int, 0, 2)
		if prevPos != -1 {
			updatedPositions = append(updatedPositions, prevPos)
		}
		updatedPositions = append(updatedPositions, currPos)

		envelope := EventEnvelope{
			Type: EventTypeQueueUpdated,
			Payload: QueueUpdatedEvent{
				Positions: updatedPositions,
			},
		}
		topic := fmt.Sprintf("chat:%s", chatID.String())
		if err := s.publisher.Publish(ctx, topic, envelope); err != nil {
			log.Error("failed to publish queue update", slog.Any("err", err))
		}
	}

	return nil
}
