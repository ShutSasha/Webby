package services

import (
	"context"
	"fmt"
	"webby/room-service/internal/apperrors"

	"github.com/google/uuid"
)

func (r *RoomService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	const op = "services.RoomService.Delete"

	room, err := r.roomRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != userID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	if err := r.roomRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.Thumbnail != "" && room.Thumbnail != defaultThumbnail {
		key := fmt.Sprintf("rooms/%s/thumbnail", id.String())
		if err := r.fileRepo.Remove(ctx, key); err != nil {
			return nil
		}
	}

	return nil
}
