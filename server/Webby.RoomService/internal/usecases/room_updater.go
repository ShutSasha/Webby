package usecases

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type UpdaterRoomRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
	Update(ctx context.Context, room *models.Room) (uuid.UUID, error)
}

type UpdaterFileRepository interface {
	Save(ctx context.Context, key string, data []byte) (string, error)
	Remove(ctx context.Context, key string) error
}

type RoomUpdater struct {
	roomRepo UpdaterRoomRepository
	fileRepo UpdaterFileRepository
}

func NewRoomUpdater(
	roomRepo UpdaterRoomRepository,
	fileRepo UpdaterFileRepository,
) *RoomUpdater {
	return &RoomUpdater{
		roomRepo: roomRepo,
		fileRepo: fileRepo,
	}
}

func (uc *RoomUpdater) Execute(
	ctx context.Context,
	roomID, userID uuid.UUID,
	name, category, thumbnailFilename *string,
	thumbnailData *[]byte,
	isPrivate *bool,
) (*models.Room, error) {
	const op = "usecases.RoomUpdater.Execute"

	existingRoom, err := uc.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if existingRoom.HostID != userID {
		return nil, fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	if name != nil {
		existingRoom.Name = *name
	}
	if category != nil {
		existingRoom.Category = *category
	}
	if isPrivate != nil {
		existingRoom.IsPrivate = *isPrivate
	}

	if thumbnailData != nil && thumbnailFilename != nil {
		if existingRoom.Thumbnail != "" && existingRoom.Thumbnail != defaultThumbnail {
			oldKey := fmt.Sprintf("rooms/%s/thumbnail", roomID.String())
			_ = uc.fileRepo.Remove(ctx, oldKey)
		}

		key := uc.generateThumbnailKey(*thumbnailFilename)
		thumbnailURL, err := uc.fileRepo.Save(ctx, key, *thumbnailData)
		if err != nil {
			return nil, fmt.Errorf("%s: thumbnail upload failed: %w", op, err)
		}
		existingRoom.Thumbnail = thumbnailURL
	}

	if _, err := uc.roomRepo.Update(ctx, existingRoom); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return existingRoom, nil
}

func (uc *RoomUpdater) generateThumbnailKey(filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("rooms/%d/thumbnail%s", time.Now().UnixMilli(), ext)
}
