package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"webby/room-service/pkg/logger"

	"github.com/google/uuid"
)

type chatIDRetriever interface {
	GetChatIDByRoomID(ctx context.Context, roomID, userID uuid.UUID) (uuid.UUID, error)
}

type eventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type timecodesRepository interface {
	RetrieveTimecodes(ctx context.Context, roomID, syncID uuid.UUID) (map[string]int, error)
	SetTimecode(ctx context.Context, userID, roomID, syncID uuid.UUID, timecode int) error
}

type SynchronizeService struct {
	chatIDRetriever chatIDRetriever
	publisher       eventPublisher
	timecodesRepo   timecodesRepository
}

func NewSynchronizeService(
	chatClient chatIDRetriever,
	publisher eventPublisher,
	timecodesRepo timecodesRepository,
) *SynchronizeService {
	return &SynchronizeService{
		chatIDRetriever: chatClient,
		publisher:       publisher,
		timecodesRepo:   timecodesRepo,
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

func (svc *SynchronizeService) Synchronize(ctx context.Context, userID, roomID uuid.UUID) error {
	const op = "service.SynchronizeService.Synchronize"
	log := logger.FromContext(ctx).With(slog.String("op", op))

	chatID, err := svc.chatIDRetriever.GetChatIDByRoomID(ctx, roomID, userID)
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
	// TODO: Move to Publish
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := svc.publisher.Publish(ctx, topic, reportEnvelope); err != nil {
		log.Error("failed to publish queue report", slog.String("err", err.Error()))
	}

	bgCtx := context.WithoutCancel(ctx)

	go func(asyncCtx context.Context, rID uuid.UUID, sID uuid.UUID, chatTopic string) {
		bgLog := logger.FromContext(asyncCtx).With(slog.String("op", op+"_async"))

		time.Sleep(syncSkipDuration)

		timecodes, err := svc.timecodesRepo.RetrieveTimecodes(asyncCtx, rID, sID)
		if err != nil {
			bgLog.Error("failed to retrieve timecodes", slog.String("err", err.Error()))
			return
		}

		timecode := svc.selectMaxTimecode(timecodes)

		synchEnvelope := eventEnvelope[synchronizePayload]{
			Type: eventTypeSynchronize,
			Payload: synchronizePayload{
				Timecode: timecode + int(syncSkipDuration.Seconds()),
			},
		}
		if err := svc.publisher.Publish(asyncCtx, chatTopic, synchEnvelope); err != nil {
			bgLog.Error("failed to publish queue sync", slog.String("err", err.Error()))
		}
	}(bgCtx, roomID, syncID, topic)

	return nil
}

func (svc *SynchronizeService) ReportTimecode(ctx context.Context, userID, roomID, syncID uuid.UUID, timecode int) error {
	const op = "services.SynchronizeService.ReportTimecode"

	err := svc.timecodesRepo.SetTimecode(ctx, userID, roomID, syncID, timecode)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (svc *SynchronizeService) selectMaxTimecode(timecodes map[string]int) int {
	result := 0

	for _, timecode := range timecodes {
		if timecode > result {
			result = timecode
		}
	}

	return result
}
