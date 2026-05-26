package services

import (
	"context"
	"fmt"
	"log/slog"
	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"
)

func (r *RoomService) Create(
	ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string,
) (*models.Room, error) {
	const op = "services.RoomService.Create"

	exists, err := r.categoryClient.Exists(ctx, room.Category)
	if err != nil {
		return nil, fmt.Errorf("%s: category check failed: %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s: category '%s': %w", op, room.Category, apperrors.ErrInvalidInput)
	}

	id, err := r.roomRepo.Create(ctx, room)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var thumbnailURL string
	if len(thumbnailData) > 0 && thumbnailFilename != "" {
		key := generateFileKey(thumbnailFilename)
		thumbnailURL, err = r.fileRepo.Save(ctx, key, thumbnailData)
		if err != nil {
			return nil, fmt.Errorf("%s: thumbnail upload failed: %w", op, err)
		}
		room.Thumbnail = thumbnailURL
	} else {
		room.Thumbnail = defaultThumbnail
	}

	room.ID = id
	if _, err := r.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("%s: failed to update room thumbnail: %w", op, err)
	}

	if err := r.roomMemberRepo.EnsureMember(ctx, room.ID, room.HostID); err != nil {
		slog.Warn("failed to create room member for host",
			slog.String("roomId", room.ID.String()),
			slog.String("hostId", room.HostID.String()),
			slog.String("error", err.Error()),
		)
	}

	if r.chatClient != nil {
		chatId, err := r.chatClient.CreateChat(ctx, room.ID)
		if err != nil {
			slog.Warn("failed to create chat for room",
				slog.String("roomId", room.ID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			room.ChatID = &chatId

			if err := r.chatClient.AddChatMember(ctx, chatId, room.HostID); err != nil {
				slog.Warn("failed to add host as chat member",
					slog.String("roomId", room.ID.String()),
					slog.String("chatId", chatId.String()),
					slog.String("hostId", room.HostID.String()),
					slog.String("error", err.Error()),
				)
			}
		}
	}

	return room, nil
}
