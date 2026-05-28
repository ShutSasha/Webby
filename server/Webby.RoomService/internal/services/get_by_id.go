package services

import (
	"context"
	"fmt"
	"log/slog"
	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

func (r *RoomService) GetByID(ctx context.Context, roomID, userID uuid.UUID) (*models.Room, error) {
	const op = "services.RoomService.GetByID"

	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if room.IsPrivate {
		isMember, err := r.roomMemberRepo.Exists(ctx, room.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if !isMember {
			return nil, fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
		}
	} else {
		if err := r.roomMemberRepo.EnsureMember(ctx, room.ID, userID); err != nil {
			slog.Warn("failed to ensure room member on GetById",
				slog.String("roomId", room.ID.String()),
				slog.String("userId", userID.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	if r.chatClient != nil {
		chatID, err := r.chatClient.GetChatByRoomID(ctx, room.ID)
		if err != nil {
			slog.Warn("failed to get chat for room",
				slog.String("roomId", room.ID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			room.ChatID = &chatID

			if err := r.chatClient.AddChatMember(ctx, chatID, userID); err != nil {
				slog.Warn("failed to add user as chat member",
					slog.String("roomID", room.ID.String()),
					slog.String("chatID", chatID.String()),
					slog.String("userID", userID.String()),
					slog.String("error", err.Error()),
				)
			}
		}
	}

	return room, nil
}
