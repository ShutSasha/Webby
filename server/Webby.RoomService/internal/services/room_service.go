package services

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type roomRepository interface {
	Create(ctx context.Context, room *models.Room) (uuid.UUID, error)
	Update(ctx context.Context, room *models.Room) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListMy(ctx context.Context, userID uuid.UUID, page, limit int, search, category string) ([]models.Room, int64, error)
	ListPublic(ctx context.Context, page, limit int, search, category string) ([]models.PublicRoom, int64, error)
}

type roomMemberChecker interface {
	GetMemberStatus(ctx context.Context, roomID, userID uuid.UUID) (*models.MemberStatus, error)
	EnsureMember(ctx context.Context, roomID, userID uuid.UUID) error
}

type fileRepository interface {
	Save(ctx context.Context, key string, data []byte) (string, error)
	Remove(ctx context.Context, key string) error
}

type categoryChecker interface {
	Exists(ctx context.Context, name string) (bool, error)
}

type roomChatManager interface {
	CreateChat(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
	GetChatByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
	AddChatMember(ctx context.Context, chatID, userID uuid.UUID) error
}

type roomService struct {
	roomRepo          roomRepository
	roomMemberChecker roomMemberChecker
	fileRepo          fileRepository
	categoryChecker   categoryChecker
	roomChatManager   roomChatManager
}

func NewRoomService(
	roomRepo roomRepository,
	roomMemberRepo roomMemberChecker,
	fileRepo fileRepository,
	categoryChecker categoryChecker,
	roomChatManager roomChatManager,
) *roomService {
	return &roomService{
		roomRepo:          roomRepo,
		roomMemberChecker: roomMemberRepo,
		fileRepo:          fileRepo,
		categoryChecker:   categoryChecker,
		roomChatManager:   roomChatManager,
	}
}

func (s *roomService) Create(
	ctx context.Context,
	room *models.Room,
	thumbnailData []byte,
	thumbnailFilename string,
) (*models.Room, error) {
	const op = "services.roomService.Create"

	exists, err := s.categoryChecker.Exists(ctx, room.Category)
	if err != nil {
		return nil, fmt.Errorf("%s: category check failed: %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s: %w", op, apperrors.ErrCategoryNotFound)
	}

	id, err := s.roomRepo.Create(ctx, room)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var thumbnailURL string
	if len(thumbnailData) > 0 && thumbnailFilename != "" {
		key := s.generateThumbnailKey(thumbnailFilename)
		thumbnailURL, err = s.fileRepo.Save(ctx, key, thumbnailData)
		if err != nil {
			return nil, fmt.Errorf("%s: thumbnail upload failed: %w", op, err)
		}
		room.Thumbnail = thumbnailURL
	} else {
		room.Thumbnail = defaultThumbnail
	}

	room.ID = id
	if _, err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("%s: failed to update room thumbnail: %w", op, err)
	}

	err = s.roomMemberChecker.EnsureMember(ctx, room.ID, room.HostID)
	if err != nil {
		slog.Warn("failed to create room member for host",
			slog.String("roomID", room.ID.String()),
			slog.String("hostID", room.HostID.String()),
			slog.String("error", err.Error()),
		)
	}

	chatID, err := s.roomChatManager.CreateChat(ctx, room.ID)
	if err != nil {
		slog.Warn("failed to create chat for room",
			slog.String("roomId", room.ID.String()),
			slog.String("error", err.Error()),
		)
	} else {
		room.ChatID = &chatID

		if err := s.roomChatManager.AddChatMember(ctx, chatID, room.HostID); err != nil {
			slog.Warn("failed to add host as chat member",
				slog.String("roomId", room.ID.String()),
				slog.String("chatId", chatID.String()),
				slog.String("hostId", room.HostID.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	return room, nil
}

func (s *roomService) AccessRoom(ctx context.Context, roomID, userID uuid.UUID) (*models.Room, error) {
	const op = "services.roomService.AccessRoom"

	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	memberStatus, err := s.roomMemberChecker.GetMemberStatus(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: get member status: %w", op, err)
	}

	if memberStatus.IsMember && memberStatus.IsBanned {
		return nil, fmt.Errorf("%s: %w", op, apperrors.ErrBanned)
	}

	if room.IsPrivate {
		if !memberStatus.IsMember {
			return nil, fmt.Errorf("%s: %w", op, apperrors.ErrNotMember)
		}
	} else {
		if !memberStatus.IsMember {
			err := s.roomMemberChecker.EnsureMember(ctx, room.ID, userID)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	chatID, err := s.roomChatManager.GetChatByRoomID(ctx, room.ID)
	if err != nil {
		slog.Warn("failed to get chat for room",
			slog.String("roomID", room.ID.String()),
			slog.String("error", err.Error()),
		)
	} else {
		room.ChatID = &chatID

		err := s.roomChatManager.AddChatMember(ctx, chatID, userID)
		if err != nil {
			slog.Warn("failed to add user as chat member",
				slog.String("roomID", room.ID.String()),
				slog.String("chatID", chatID.String()),
				slog.String("userID", userID.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	return room, nil
}

func (s *roomService) ListMyRooms(
	ctx context.Context,
	userID uuid.UUID,
	page, limit int,
	search, category string,
) ([]models.Room, int64, error) {
	const op = "services.RoomService.ListMyRooms"

	if strings.ToLower(category) == "all" {
		category = ""
	}

	rooms, total, err := s.roomRepo.ListMy(ctx, userID, page, limit, search, category)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return rooms, total, nil
}

func (s *roomService) ListPublicRooms(ctx context.Context, page, limit int, search, category string) ([]models.PublicRoom, int64, error) {
	const op = "services.roomService.ListPublicRooms"

	if strings.ToLower(category) == "all" {
		category = ""
	}

	rooms, total, err := s.roomRepo.ListPublic(ctx, page, limit, search, category)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return rooms, total, nil
}

func (s *roomService) Update(
	ctx context.Context,
	roomID, userID uuid.UUID,
	name, category, thumbnailFilename *string,
	thumbnailData *[]byte,
	isPrivate *bool,
) (*models.Room, error) {
	const op = "services.RoomService.Update"

	existingRoom, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if existingRoom.HostID != userID {
		return nil, fmt.Errorf("%s: %w", op, apperrors.ErrNotHost)
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
			_ = s.fileRepo.Remove(ctx, oldKey)
		}

		key := s.generateThumbnailKey(*thumbnailFilename)
		thumbnailURL, err := s.fileRepo.Save(ctx, key, *thumbnailData)
		if err != nil {
			return nil, fmt.Errorf("%s: thumbnail upload failed: %w", op, err)
		}
		existingRoom.Thumbnail = thumbnailURL
	}

	if _, err := s.roomRepo.Update(ctx, existingRoom); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return existingRoom, nil
}

func (s *roomService) Delete(ctx context.Context, roomID, userID uuid.UUID) error {
	const op = "services.RoomService.Delete"

	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != userID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNotHost)
	}

	if err := s.roomRepo.Delete(ctx, roomID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.Thumbnail != "" && room.Thumbnail != defaultThumbnail {
		key := fmt.Sprintf("rooms/%s/thumbnail", roomID.String())
		if err := s.fileRepo.Remove(ctx, key); err != nil {
			return nil
		}
	}

	return nil
}

func (s *roomService) generateThumbnailKey(filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("rooms/%d/thumbnail%s", time.Now().UnixMilli(), ext)
}
