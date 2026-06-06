package usecases

import (
	"context"
	"fmt"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type DeleterRoomRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DeleterFileRepository interface {
	Remove(ctx context.Context, key string) error
}

type RoomDeleter struct {
	roomRepo DeleterRoomRepository
	fileRepo DeleterFileRepository
}

func NewRoomDeleter(roomRepo DeleterRoomRepository, fileRepo DeleterFileRepository) *RoomDeleter {
	return &RoomDeleter{
		roomRepo: roomRepo,
		fileRepo: fileRepo,
	}
}

func (uc *RoomDeleter) Execute(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) error {
	const op = "usecases.RoomDeleter.Execute"

	room, err := uc.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != userID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	if err := uc.roomRepo.Delete(ctx, roomID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.Thumbnail != "" && room.Thumbnail != defaultThumbnail {
		key := fmt.Sprintf("rooms/%s/thumbnail", roomID.String())
		if err := uc.fileRepo.Remove(ctx, key); err != nil {
			return nil
		}
	}

	return nil
}
