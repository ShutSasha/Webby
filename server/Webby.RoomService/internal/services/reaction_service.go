package services

import (
	"context"
	"fmt"
	"log/slog"
	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"
	"webby/room-service/pkg/logger"

	"github.com/google/uuid"
)

const (
	eventRoomPointsUpdated = "ROOM_POINTS_UPDATED"
	eventReactionSent      = "REACTION_SENT"
)

type reactionRepository interface {
	List(ctx context.Context) ([]models.Reaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Reaction, error)
}

type publisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type pointsUpdater interface {
	AddPointsBulk(ctx context.Context, roomID uuid.UUID, userIDs []uuid.UUID, pointsToAdd int) (map[uuid.UUID]int, error)
	GetMemberPoints(ctx context.Context, roomID, userID uuid.UUID) (int, error)
}

type reactionService struct {
	reactionRepo    reactionRepository
	memberChecker   memberChecker
	chatIDRetriever chatIDRetriever
	pointsUpdater   pointsUpdater
	publisher       publisher
}

func NewReactionService(reactionRepo reactionRepository, memberChecker memberChecker, chatIDRetriever chatIDRetriever, pointsUpdater pointsUpdater, publisher publisher) *reactionService {
	return &reactionService{
		reactionRepo:    reactionRepo,
		memberChecker:   memberChecker,
		chatIDRetriever: chatIDRetriever,
		pointsUpdater:   pointsUpdater,
		publisher:       publisher,
	}
}

func (s *reactionService) List(ctx context.Context) ([]models.Reaction, error) {
	const op = "services.reactionService.List"

	reactions, err := s.reactionRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return reactions, nil
}

func (s *reactionService) UseReaction(ctx context.Context, roomID, userID, reactionID uuid.UUID) error {
	const op = "reactionService.UseReaction"
	log := logger.FromContext(ctx).With("op", op)

	status, err := s.memberChecker.GetMemberStatus(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !status.IsMember {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNotMember)
	}
	if status.IsBanned {
		return fmt.Errorf("%s: %w", op, apperrors.ErrBanned)
	}

	chatID, err := s.chatIDRetriever.GetChatIDByRoomID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	reaction, err := s.reactionRepo.GetByID(ctx, reactionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userPoints, err := s.pointsUpdater.GetMemberPoints(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if userPoints - reaction.Cost < 0 {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNoPoints)
	}

	userPointsMap, err := s.pointsUpdater.AddPointsBulk(ctx, roomID, []uuid.UUID{userID}, -reaction.Cost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	updatePointsEnvelope := eventEnvelope[map[string]any]{
		Type: eventRoomPointsUpdated,
		Payload: map[string]any{
			"added_points": -reaction.Cost,
			"totals":       userPointsMap,
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	err = s.publisher.Publish(ctx, topic, updatePointsEnvelope)
	if err != nil {
		log.Error("failed to publish update points", slog.String("err", err.Error()))
	}

	reactionSentEnvelope := eventEnvelope[*models.Reaction]{
		Type: eventReactionSent,
		Payload: reaction,
	}
	err = s.publisher.Publish(ctx, topic, reactionSentEnvelope)
	if err != nil {
		log.Error("failed to publish reaction sent", slog.String("err", err.Error()))
	}

	return nil
}
