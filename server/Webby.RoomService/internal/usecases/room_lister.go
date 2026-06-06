package usecases

import (
	"context"
	"fmt"

	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type RoomListerRepository interface {
	ListMy(
		ctx context.Context,
		userID uuid.UUID,
		page, limit int,
		search, category string,
	) ([]models.Room, int64, error)
	ListPublic(
		ctx context.Context,
		page, limit int,
		search, category string,
	) ([]models.PublicRoom, int64, error)
}

type RoomLister struct {
	repo RoomListerRepository
}

func NewRoomLister(repo RoomListerRepository) *RoomLister {
	return &RoomLister{repo: repo}
}

func (uc *RoomLister) ExecuteListMy(
	ctx context.Context,
	userID uuid.UUID,
	page, limit int,
	search, category string,
) ([]models.Room, int64, error) {
	const op = "usecases.RoomLister.ExecuteListMy"

	rooms, total, err := uc.repo.ListMy(ctx, userID, page, limit, search, category)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return rooms, total, nil
}

func (uc *RoomLister) ExecuteListPublic(
	ctx context.Context,
	page, limit int,
	search, category string,
) ([]models.PublicRoom, int64, error) {
	const op = "usecases.RoomLister.ExecuteListPublic"

	rooms, total, err := uc.repo.ListPublic(ctx, page, limit, search, category)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return rooms, total, nil
}
