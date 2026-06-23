package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"
	"webby/room-service/pkg/logger"

	"github.com/google/uuid"
)

type chatIDRetriever interface {
	GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
}

type eventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type timecodesRepository interface {
	RetrieveTimecodes(ctx context.Context, roomID, syncID uuid.UUID) (map[string]int, error)
	SetTimecode(ctx context.Context, userID, roomID, syncID uuid.UUID, timecode int) error
}

type memberChecker interface {
	GetMemberStatus(ctx context.Context, roomID, userID uuid.UUID) (*models.MemberStatus, error)
}

type synchronizeService struct {
	chatIDRetriever chatIDRetriever
	publisher       eventPublisher
	timecodesRepo   timecodesRepository
	memberChecker   memberChecker
}

func NewSynchronizeService(
	chatClient chatIDRetriever,
	publisher eventPublisher,
	timecodesRepo timecodesRepository,
	memberChecker memberChecker,
) *synchronizeService {
	return &synchronizeService{
		chatIDRetriever: chatClient,
		publisher:       publisher,
		timecodesRepo:   timecodesRepo,
		memberChecker:   memberChecker,
	}
}

type reportPayload struct {
	SyncID string `json:"synchronizeId"`
}

type synchronizePayload struct {
	Timecode int `json:"timecode"`
}

type eventEnvelope[T any] struct {
	Type    string `json:"type"`
	Payload T      `json:"payload"`
}

const (
	eventTypeReportTimecode = "REPORT_TIMECODE"
	eventTypeSynchronize    = "SYNCHRONIZE"
	syncSkipDuration        = 1 * time.Second
)

func (s *synchronizeService) Synchronize(ctx context.Context, userID, roomID uuid.UUID) error {
	const op = "service.synchronizeService.Synchronize"
	log := logger.FromContext(ctx).With(slog.String("op", op))

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

	syncID := uuid.New()
	reportEnvelope := eventEnvelope[reportPayload]{
		Type: eventTypeReportTimecode,
		Payload: reportPayload{
			SyncID: syncID.String(),
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.publisher.Publish(ctx, topic, reportEnvelope); err != nil {
		log.Error("failed to publish queue report", slog.String("err", err.Error()))
	}

	bgCtx := context.WithoutCancel(ctx)

	go func(asyncCtx context.Context, rID uuid.UUID, sID uuid.UUID, chatTopic string) {
		bgLog := logger.FromContext(asyncCtx).With(slog.String("op", op+"_async"))

		time.Sleep(syncSkipDuration)

		timecodes, err := s.timecodesRepo.RetrieveTimecodes(asyncCtx, rID, sID)
		if err != nil {
			bgLog.Error("failed to retrieve timecodes", slog.String("err", err.Error()))
			return
		}

		timecode := s.selectMaxTimecode(timecodes)

		synchEnvelope := eventEnvelope[synchronizePayload]{
			Type: eventTypeSynchronize,
			Payload: synchronizePayload{
				Timecode: timecode + int(syncSkipDuration.Seconds()),
			},
		}
		if err := s.publisher.Publish(asyncCtx, chatTopic, synchEnvelope); err != nil {
			bgLog.Error("failed to publish queue sync", slog.String("err", err.Error()))
		}
	}(bgCtx, roomID, syncID, topic)

	return nil
}

func (s *synchronizeService) ReportTimecode(ctx context.Context, userID, roomID, syncID uuid.UUID, timecode int) error {
	const op = "services.synchronizeService.ReportTimecode"

	err := s.timecodesRepo.SetTimecode(ctx, userID, roomID, syncID, timecode)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *synchronizeService) selectMaxTimecode(timecodes map[string]int) int {
	result := 0

	for _, timecode := range timecodes {
		if timecode > result {
			result = timecode
		}
	}

	return result
}
