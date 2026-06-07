package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type CreatorRoomRepository interface {
	Create(ctx context.Context, room *models.Room) (uuid.UUID, error)
	Update(ctx context.Context, room *models.Room) (uuid.UUID, error)
}

type CreatorRoomMemberRepository interface {
	EnsureMember(ctx context.Context, roomID, userID uuid.UUID) error
}

type CreatorFileRepository interface {
	Save(ctx context.Context, key string, data []byte) (string, error)
}

type CreatorCategoryClient interface {
	Exists(ctx context.Context, name string) (bool, error)
}

type CreatorChatClient interface {
	CreateChat(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
	AddChatMember(ctx context.Context, chatID, userID uuid.UUID) error
}

type RoomCreator struct {
	roomRepo       CreatorRoomRepository
	roomMemberRepo CreatorRoomMemberRepository
	fileRepo       CreatorFileRepository
	categoryClient CreatorCategoryClient
	chatClient     CreatorChatClient
}

func NewRoomCreator(
	roomRepo CreatorRoomRepository,
	roomMemberRepo CreatorRoomMemberRepository,
	fileRepo CreatorFileRepository,
	categoryClient CreatorCategoryClient,
	chatClient CreatorChatClient,
) *RoomCreator {
	return &RoomCreator{
		roomRepo:       roomRepo,
		roomMemberRepo: roomMemberRepo,
		fileRepo:       fileRepo,
		categoryClient: categoryClient,
		chatClient:     chatClient,
	}
}

func (uc *RoomCreator) Execute(
	ctx context.Context,
	room *models.Room,
	thumbnailData []byte,
	thumbnailFilename string,
) (*models.Room, error) {
	const op = "usecases.RoomCreator.Execute"

	exists, err := uc.categoryClient.Exists(ctx, room.Category)
	if err != nil {
		return nil, fmt.Errorf("%s: category check failed: %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s: category '%s': %w", op, room.Category, apperrors.ErrInvalidInput)
	}

	id, err := uc.roomRepo.Create(ctx, room)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var thumbnailURL string
	if len(thumbnailData) > 0 && thumbnailFilename != "" {
		key := uc.generateThumbnailKey(thumbnailFilename)
		thumbnailURL, err = uc.fileRepo.Save(ctx, key, thumbnailData)
		if err != nil {
			return nil, fmt.Errorf("%s: thumbnail upload failed: %w", op, err)
		}
		room.Thumbnail = thumbnailURL
	} else {
		room.Thumbnail = defaultThumbnail
	}

	room.ID = id
	if _, err := uc.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("%s: failed to update room thumbnail: %w", op, err)
	}

	err = uc.roomMemberRepo.EnsureMember(ctx, room.ID, room.HostID)
	if err != nil {
		slog.Warn("failed to create room member for host",
			slog.String("roomID", room.ID.String()),
			slog.String("hostID", room.HostID.String()),
			slog.String("error", err.Error()),
		)
	}

	if uc.chatClient != nil {
		chatId, err := uc.chatClient.CreateChat(ctx, room.ID)
		if err != nil {
			slog.Warn("failed to create chat for room",
				slog.String("roomId", room.ID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			room.ChatID = &chatId

			if err := uc.chatClient.AddChatMember(ctx, chatId, room.HostID); err != nil {
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

func (uc *RoomCreator) generateThumbnailKey(filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("rooms/%d/thumbnail%s", time.Now().UnixMilli(), ext)
}
