package services

import (
	"context"
	"fmt"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

func (r *RoomService) ListMembers(
	ctx context.Context,
	roomID uuid.UUID,
	page, limit int,
	search string,
) ([]models.RoomMemberInfo, int64, error) {
	const op = "services.RoomService.ListMembers"

	members, total, err := r.roomMemberRepo.ListByRoom(ctx, roomID, page, limit, search)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return members, total, nil
}
