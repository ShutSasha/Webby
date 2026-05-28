package services

import (
	"context"
	"fmt"
	"log/slog"
	"webby/room-service/internal/apperrors"

	"github.com/google/uuid"
)

func (r *RoomService) AddMembers(ctx context.Context, roomID, hostID uuid.UUID, memberIDs []uuid.UUID) error {
	const op = "services.RoomService.AddMembers"

	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != hostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	for _, memberID := range memberIDs {
		if err := r.roomMemberRepo.EnsureMember(ctx, roomID, memberID); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if r.chatClient != nil {
		chatID, err := r.chatClient.GetChatByRoomID(ctx, roomID)
		if err != nil {
			slog.Warn("failed to get chat for room when adding members",
				slog.String("roomID", roomID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			for _, memberID := range memberIDs {
				if err := r.chatClient.AddChatMember(ctx, chatID, memberID); err != nil {
					slog.Warn("failed to add room member as chat member",
						slog.String("roomID", roomID.String()),
						slog.String("chatID", chatID.String()),
						slog.String("memberID", memberID.String()),
						slog.String("error", err.Error()),
					)
				}
			}
		}
	}

	return nil
}
