package services

import (
	"context"
	"fmt"
	"webby/room-service/internal/apperrors"

	"github.com/google/uuid"
)

func (r *RoomService) RemoveMember(
	ctx context.Context,
	roomID, memberID, hostID uuid.UUID,
) error {
	const op = "services.RoomService.RemoveMember"

	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != hostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	if memberID == room.HostID {
		return fmt.Errorf("%s: %w: cannot remove host from room", op, apperrors.ErrInvalidInput)
	}

	return r.roomMemberRepo.Delete(ctx, roomID, memberID)
}
