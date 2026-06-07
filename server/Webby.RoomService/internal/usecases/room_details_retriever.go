package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type DetailsRetrieverRoomRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
}

type DetailsRetrieverRoomMemberRepository interface {
	Exists(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
	EnsureMember(ctx context.Context, roomID, userID uuid.UUID) error
}

type DetailsRetrieverChatClient interface {
	GetChatByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
	AddChatMember(ctx context.Context, chatID, userID uuid.UUID) error
}

type RoomDetailsRetriever struct {
	roomRepo       DetailsRetrieverRoomRepository
	roomMemberRepo DetailsRetrieverRoomMemberRepository
	chatClient     DetailsRetrieverChatClient
}

func NewRoomDetailsRetriever(
	roomRepo DetailsRetrieverRoomRepository,
	roomMemberRepo DetailsRetrieverRoomMemberRepository,
	chatClient DetailsRetrieverChatClient,
) *RoomDetailsRetriever {
	return &RoomDetailsRetriever{
		roomRepo:       roomRepo,
		roomMemberRepo: roomMemberRepo,
		chatClient:     chatClient,
	}
}

func (uc *RoomDetailsRetriever) Execute(ctx context.Context, roomID, userID uuid.UUID) (*models.Room, error) {
	const op = "usecases.RoomDetailsRetriever.Execute"

	room, err := uc.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if room.IsPrivate {
		isMember, err := uc.roomMemberRepo.Exists(ctx, room.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if !isMember {
			return nil, fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
		}
	} else {
		err := uc.roomMemberRepo.EnsureMember(ctx, room.ID, userID)
		if err != nil {
			slog.Warn("failed to ensure room member on GetById",
				slog.String("roomID", room.ID.String()),
				slog.String("userID", userID.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	if uc.chatClient != nil {
		chatID, err := uc.chatClient.GetChatByRoomID(ctx, room.ID)
		if err != nil {
			slog.Warn("failed to get chat for room",
				slog.String("roomID", room.ID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			room.ChatID = &chatID

			err := uc.chatClient.AddChatMember(ctx, chatID, userID)
			if err != nil {
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
