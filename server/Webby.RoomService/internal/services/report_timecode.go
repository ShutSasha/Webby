package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *RoomService) ReportTimecode(ctx context.Context, userID, roomID, syncID uuid.UUID, timecode int) error {
	const op = "services.RoomService.ReportTimecode"

	err := s.timecodesRepo.SetTimecode(ctx, userID, roomID, syncID, timecode)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
