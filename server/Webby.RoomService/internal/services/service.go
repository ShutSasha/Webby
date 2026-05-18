package services

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

const defaultThumbnail = "https://webby-watch-platform-bucket.s3.eu-north-1.amazonaws.com/rooms/default-room-preview.jpg"

type RoomRepository interface {
	Create(ctx context.Context, room *models.Room) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
	ListMy(
		ctx context.Context,
		userID uuid.UUID,
		page, limit int,
		search string, category *string,
	) ([]models.Room, int64, error)
	ListPublic(
		ctx context.Context,
		page, limit int,
		search string, category *string,
	) ([]models.PublicRoom, int64, error)
	Update(ctx context.Context, room *models.Room) (uuid.UUID, error)
}

type RoomMemberRepository interface {
	Create(ctx context.Context, member *models.RoomMember) error
	Exists(ctx context.Context, roomId, userID uuid.UUID) (bool, error)
	EnsureMember(ctx context.Context, roomID, userID uuid.UUID) error
	Delete(ctx context.Context, roomID, userID uuid.UUID) error
	ListByRoom(
		ctx context.Context,
		roomID uuid.UUID, page, limit int,
		search string,
	) ([]models.RoomMemberInfo, int64, error)
}

type FileRepository interface {
	Save(ctx context.Context, key string, data []byte) (string, error)
	Remove(ctx context.Context, key string) error
}

type ChatClientInterface interface {
	CreateChat(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
	GetChatByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
	AddChatMember(ctx context.Context, chatID, userID uuid.UUID) error
}

type CategoryClientInterface interface {
	Exists(ctx context.Context, name string) (bool, error)
}

type RoomService struct {
	roomRepo       RoomRepository
	roomMemberRepo RoomMemberRepository
	fileRepo       FileRepository
	chatClient     ChatClientInterface
	categoryClient CategoryClientInterface
}

func NewRoomService(
	repo RoomRepository,
	roomMemberRepo RoomMemberRepository,
	fileRepo FileRepository,
	chatClient ChatClientInterface,
	categoryClient CategoryClientInterface,
) *RoomService {
	return &RoomService{
		roomRepo:       repo,
		roomMemberRepo: roomMemberRepo,
		fileRepo:       fileRepo,
		chatClient:     chatClient,
		categoryClient: categoryClient,
	}
}

func generateFileKey(filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("rooms/%d/thumbnail%s", time.Now().UnixMilli(), ext)
}

func (r *RoomService) Create(ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string) (*models.Room, error) {
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

func (r *RoomService) Delete(ctx context.Context, id uuid.UUID, userId uuid.UUID) error {
	const op = "services.RoomService.Delete"

	room, err := r.roomRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != userId {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	if err := r.roomRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.Thumbnail != "" && room.Thumbnail != defaultThumbnail {
		key := fmt.Sprintf("rooms/%s/thumbnail", id.String())
		if err := r.fileRepo.Remove(ctx, key); err != nil {
			return nil
		}
	}

	return nil
}

func (r *RoomService) GetByID(ctx context.Context, roomId uuid.UUID, userId uuid.UUID) (*models.Room, error) {
	const op = "services.RoomService.GetByID"

	room, err := r.roomRepo.GetByID(ctx, roomId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if room.IsPrivate {
		isMember, err := r.roomMemberRepo.Exists(ctx, room.ID, userId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if !isMember {
			return nil, fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
		}
	} else {
		if err := r.roomMemberRepo.EnsureMember(ctx, room.ID, userId); err != nil {
			slog.Warn("failed to ensure room member on GetById",
				slog.String("roomId", room.ID.String()),
				slog.String("userId", userId.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	if r.chatClient != nil {
		chatId, err := r.chatClient.GetChatByRoomID(ctx, room.ID)
		if err != nil {
			slog.Warn("failed to get chat for room",
				slog.String("roomId", room.ID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			room.ChatID = &chatId

			if err := r.chatClient.AddChatMember(ctx, chatId, userId); err != nil {
				slog.Warn("failed to add user as chat member",
					slog.String("roomId", room.ID.String()),
					slog.String("chatId", chatId.String()),
					slog.String("userId", userId.String()),
					slog.String("error", err.Error()),
				)
			}
		}
	}

	return room, nil
}

func (r *RoomService) ListMy(
	ctx context.Context,
	userID uuid.UUID,
	page, limit int,
	search string, category *string,
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
	page int, limit int,
	search string, category *string,
) ([]models.PublicRoom, int64, error) {
	const op = "services.RoomService.ListPublic"

	rooms, total, err := r.roomRepo.ListPublic(ctx, page, limit, search, category)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return rooms, total, nil
}

func (r *RoomService) ListMembers(
	ctx context.Context,
	roomID uuid.UUID,
	page, limit int,
	search string,
) ([]models.RoomMemberInfo, int64, error) {
	const op = "services.RoomService.ListMembers"

	members, total, err := r.roomMemberRepo.ListByRoom(ctx, roomID, page, limit, search)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return members, total, nil
}

func (r *RoomService) AddMembers(
	ctx context.Context,
	roomID uuid.UUID, memberIDs []uuid.UUID, userID uuid.UUID,
) error {
	const op = "services.RoomService.AddMembers"

	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != userID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	for _, memberId := range memberIDs {
		if err := r.roomMemberRepo.EnsureMember(ctx, roomID, memberId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if r.chatClient != nil {
		chatId, err := r.chatClient.GetChatByRoomID(ctx, roomID)
		if err != nil {
			slog.Warn("failed to get chat for room when adding members",
				slog.String("roomId", roomID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			for _, memberId := range memberIDs {
				if err := r.chatClient.AddChatMember(ctx, chatId, memberId); err != nil {
					slog.Warn("failed to add room member as chat member",
						slog.String("roomId", roomID.String()),
						slog.String("chatId", chatId.String()),
						slog.String("memberId", memberId.String()),
						slog.String("error", err.Error()),
					)
				}
			}
		}
	}

	return nil
}

func (r *RoomService) RemoveMember(
	ctx context.Context,
	roomID, memberID, userID uuid.UUID,
) error {
	const op = "services.RoomService.RemoveMember"

	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != userID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	if memberID == room.HostID {
		return fmt.Errorf("%s: %w: cannot remove host from room", op, apperrors.ErrInvalidInput)
	}

	return r.roomMemberRepo.Delete(ctx, roomID, memberID)
}

func (r *RoomService) Update(
	ctx context.Context,
	roomID uuid.UUID, name *string, category *string, isPrivate *bool,
	thumbnailData *[]byte, thumbnailFilename *string, userID uuid.UUID,
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
