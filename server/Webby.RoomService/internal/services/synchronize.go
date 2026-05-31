package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"webby/room-service/pkg/logger"

	"github.com/google/uuid"
)

type ReportPayload struct {
	SyncID string `json:"synchronizeId"`
}

type SynchronizePayload struct {
	Timecode int `json:"timecode"`
}

type EventEnvelope[T any] struct {
	Type    string `json:"type"`
	Payload T      `json:"payload"`
}

const (
	EventTypeReportTimecode = "REPORT_TIMECODE"
	EventTypeSynchronize    = "SYNCHRONIZE"
)

func (s *RoomService) Synchronize(ctx context.Context, userID, roomID uuid.UUID) error {
	const op = "services.RoomService.Synchronize"
	log := logger.FromContext(ctx).With(slog.String("op", op))

	chatID, err := s.chatClient.GetChatIDByRoomID(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	syncID := uuid.New()
	reportEnvelope := EventEnvelope[ReportPayload]{
		Type: EventTypeReportTimecode,
		Payload: ReportPayload{
			SyncID: syncID.String(),
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.publisher.Publish(ctx, topic, reportEnvelope); err != nil {
		log.Error("failed to publish queue report", slog.Any("err", err))
	}

	bgCtx := context.WithoutCancel(ctx)

	go func(asyncCtx context.Context, rID uuid.UUID, sID uuid.UUID, chatTopic string) {
		bgLog := logger.FromContext(asyncCtx).With(slog.String("op", op+"_async"))

		skipTime := 1 * time.Second
		time.Sleep(skipTime)

		timecodes, err := s.timecodesRepo.RetrieveTimecodes(asyncCtx, rID, sID)
		if err != nil {
			bgLog.Error("failed to retrieve timecodes", slog.Any("err", err))
			return
		}

		timecode := s.calculateTimecode(timecodes)

		synchEnvelope := EventEnvelope[SynchronizePayload]{
			Type: EventTypeSynchronize,
			Payload: SynchronizePayload{
				Timecode: timecode + int(skipTime.Seconds()),
			},
		}
		if err := s.publisher.Publish(asyncCtx, chatTopic, synchEnvelope); err != nil {
			bgLog.Error("failed to publish queue sync", slog.Any("err", err))
		}
	}(bgCtx, roomID, syncID, topic)

	return nil
}

func (s *RoomService) calculateTimecode(timecodes map[string]int) int {
	result := 0

	for _, timecode := range timecodes {
		if timecode > result {
			result = timecode
		}
	}

	return result
}
