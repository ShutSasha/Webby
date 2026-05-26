package services

import (
	"context"
	"fmt"
	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

func (r *RoomService) Update(
	ctx context.Context,
	roomID, userID uuid.UUID,
	name, category, thumbnailFilename *string,
	thumbnailData *[]byte, isPrivate *bool,
) (*models.Room, error) {
	const op = "services.RoomService.Update"

	existingRoom, err := r.roomRepo.GetByID(ctx, roomID)
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
			_ = r.fileRepo.Remove(ctx, oldKey)
		}

		key := generateFileKey(*thumbnailFilename)
		thumbnailURL, err := r.fileRepo.Save(ctx, key, *thumbnailData)
		if err != nil {
			return nil, fmt.Errorf("thumbnail upload failed: %w", err)
		}
		existingRoom.Thumbnail = thumbnailURL
	}

	id, err := r.roomRepo.Update(ctx, existingRoom)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	updatedRoom, err := r.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updatedRoom, nil
}
