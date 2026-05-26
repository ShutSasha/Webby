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

	time.Sleep(800 * time.Millisecond)

	timecodes, err := s.timecodesRepo.RetrieveTimecodes(ctx, roomID, syncID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	timecode := s.calculateTimecode(timecodes)

	synchEnvelope := EventEnvelope[SynchronizePayload]{
		Type: EventTypeSynchronize,
		Payload: SynchronizePayload{
			Timecode: timecode,
		},
	}
	if err := s.publisher.Publish(ctx, topic, synchEnvelope); err != nil {
		log.Error("failed to publish queue sync", slog.Any("err", err))
	}

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
