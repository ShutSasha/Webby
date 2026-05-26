package services

import (
	"context"
	"fmt"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

func (r *RoomService) ListMy(
	ctx context.Context,
	userID uuid.UUID,
	page, limit int,
	search, category string,
) ([]models.Room, int64, error) {
	const op = "services.RoomService.ListMy"

	rooms, total, err := r.roomRepo.ListMy(ctx, userID, page, limit, search, category)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return rooms, total, nil
}

func (r *RoomService) ListPublic(
	ctx context.Context,
	page, limit int,
	search, category string,
) ([]models.PublicRoom, int64, error) {
	const op = "services.RoomService.ListPublic"

	rooms, total, err := r.roomRepo.ListPublic(ctx, page, limit, search, category)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return rooms, total, nil
}
