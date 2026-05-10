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
	GetById(ctx context.Context, id uuid.UUID) (*models.Room, error)
	ListMy(ctx context.Context, userId uuid.UUID, page int, limit int, search string, categoryName *string) ([]models.Room, int64, error)
	ListPublic(ctx context.Context, page int, limit int, search string, categoryName *string) ([]models.PublicRoom, int64, error)
	Update(ctx context.Context, room *models.Room) (uuid.UUID, error)
}

type RoomMemberRepository interface {
	Create(ctx context.Context, member *models.RoomMember) error
	Exists(ctx context.Context, roomId, userId uuid.UUID) (bool, error)
	EnsureMember(ctx context.Context, roomId, userId uuid.UUID) error
	Delete(ctx context.Context, roomId, userId uuid.UUID) error
	ListByRoom(ctx context.Context, roomId uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error)
	UpdatePoints(ctx context.Context, roomId, userId uuid.UUID, delta int) (*models.RoomMemberInfo, error)
}

type FileRepository interface {
	Save(ctx context.Context, key string, data []byte) (string, error)
	Remove(ctx context.Context, key string) error
}

type ChatClientInterface interface {
	CreateChat(ctx context.Context, roomId uuid.UUID) (uuid.UUID, error)
	GetChatByRoomId(ctx context.Context, roomId uuid.UUID) (uuid.UUID, error)
	AddChatMember(ctx context.Context, chatId, userId uuid.UUID) error
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

func NewRoomService(repo RoomRepository, roomMemberRepo RoomMemberRepository, fileRepo FileRepository, chatClient ChatClientInterface, categoryClient CategoryClientInterface) *RoomService {
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
	exists, err := r.categoryClient.Exists(ctx, room.CategoryName)
	if err != nil {
		return nil, fmt.Errorf("category check failed: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("category '%s': %w", room.CategoryName, apperrors.ErrInvalidInput)
	}

	id, err := r.roomRepo.Create(ctx, room)
	if err != nil {
		return nil, err
	}

	var thumbnailURL string
	if len(thumbnailData) > 0 && thumbnailFilename != "" {
		key := generateFileKey(thumbnailFilename)
		thumbnailURL, err = r.fileRepo.Save(ctx, key, thumbnailData)
		if err != nil {
			return nil, fmt.Errorf("thumbnail upload failed: %w", err)
		}
		room.Thumbnail = thumbnailURL
	} else {
		room.Thumbnail = defaultThumbnail
	}

	room.Id = id
	if _, err := r.roomRepo.Update(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to update room thumbnail: %w", err)
	}

	if err := r.roomMemberRepo.EnsureMember(ctx, room.Id, room.HostId); err != nil {
		slog.Warn("failed to create room member for host",
			slog.String("roomId", room.Id.String()),
			slog.String("hostId", room.HostId.String()),
			slog.String("error", err.Error()),
		)
	}

	if r.chatClient != nil {
		chatId, err := r.chatClient.CreateChat(ctx, room.Id)
		if err != nil {
			slog.Warn("failed to create chat for room",
				slog.String("roomId", room.Id.String()),
				slog.String("error", err.Error()),
			)
		} else {
			room.ChatId = &chatId

			if err := r.chatClient.AddChatMember(ctx, chatId, room.HostId); err != nil {
				slog.Warn("failed to add host as chat member",
					slog.String("roomId", room.Id.String()),
					slog.String("chatId", chatId.String()),
					slog.String("hostId", room.HostId.String()),
					slog.String("error", err.Error()),
				)
			}
		}
	}

	return room, nil
}

func (r *RoomService) Delete(ctx context.Context, id uuid.UUID, userId uuid.UUID) error {
	room, err := r.roomRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if room.HostId != userId {
		return apperrors.ErrForbidden
	}

	if err := r.roomRepo.Delete(ctx, id); err != nil {
		return err
	}

	if room.Thumbnail != "" && room.Thumbnail != defaultThumbnail {
		key := fmt.Sprintf("rooms/%s/thumbnail", id.String())
		if err := r.fileRepo.Remove(ctx, key); err != nil {
			return nil
		}
	}

	return nil
}

func (r *RoomService) GetById(ctx context.Context, roomId uuid.UUID, userId uuid.UUID) (*models.Room, error) {
	room, err := r.roomRepo.GetById(ctx, roomId)
	if err != nil {
		return nil, err
	}

	if room.IsPrivate {
		isMember, err := r.roomMemberRepo.Exists(ctx, room.Id, userId)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, apperrors.ErrForbidden
		}
	} else {
		if err := r.roomMemberRepo.EnsureMember(ctx, room.Id, userId); err != nil {
			slog.Warn("failed to ensure room member on GetById",
				slog.String("roomId", room.Id.String()),
				slog.String("userId", userId.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	if r.chatClient != nil {
		chatId, err := r.chatClient.GetChatByRoomId(ctx, room.Id)
		if err != nil {
			slog.Warn("failed to get chat for room",
				slog.String("roomId", room.Id.String()),
				slog.String("error", err.Error()),
			)
		} else {
			room.ChatId = &chatId

			if err := r.chatClient.AddChatMember(ctx, chatId, userId); err != nil {
				slog.Warn("failed to add user as chat member",
					slog.String("roomId", room.Id.String()),
					slog.String("chatId", chatId.String()),
					slog.String("userId", userId.String()),
					slog.String("error", err.Error()),
				)
			}
		}
	}

	return room, nil
}

func (r *RoomService) ListMy(ctx context.Context, userId uuid.UUID, page int, limit int, search string, categoryName *string) ([]models.Room, int64, error) {
	rooms, total, err := r.roomRepo.ListMy(ctx, userId, page, limit, search, categoryName)
	if err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *RoomService) ListPublic(ctx context.Context, page int, limit int, search string, categoryName *string) ([]models.PublicRoom, int64, error) {
	rooms, total, err := r.roomRepo.ListPublic(ctx, page, limit, search, categoryName)
	if err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *RoomService) ListMembers(ctx context.Context, roomId uuid.UUID, page int, limit int, search string) ([]models.RoomMemberInfo, int64, error) {
	members, total, err := r.roomMemberRepo.ListByRoom(ctx, roomId, page, limit, search)
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

func (r *RoomService) AddMembers(ctx context.Context, roomId uuid.UUID, memberIds []uuid.UUID, userId uuid.UUID) error {
	room, err := r.roomRepo.GetById(ctx, roomId)
	if err != nil {
		return err
	}

	if room.HostId != userId {
		return apperrors.ErrForbidden
	}

	for _, memberId := range memberIds {
		if err := r.roomMemberRepo.EnsureMember(ctx, roomId, memberId); err != nil {
			return err
		}
	}

	if r.chatClient != nil {
		chatId, err := r.chatClient.GetChatByRoomId(ctx, roomId)
		if err != nil {
			slog.Warn("failed to get chat for room when adding members",
				slog.String("roomId", roomId.String()),
				slog.String("error", err.Error()),
			)
		} else {
			for _, memberId := range memberIds {
				if err := r.chatClient.AddChatMember(ctx, chatId, memberId); err != nil {
					slog.Warn("failed to add room member as chat member",
						slog.String("roomId", roomId.String()),
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

func (r *RoomService) RemoveMember(ctx context.Context, roomId uuid.UUID, memberId uuid.UUID, userId uuid.UUID) error {
	room, err := r.roomRepo.GetById(ctx, roomId)
	if err != nil {
		return err
	}

	if room.HostId != userId {
		return apperrors.ErrForbidden
	}

	if memberId == room.HostId {
		return fmt.Errorf("%w: cannot remove the host from the room", apperrors.ErrInvalidInput)
	}

	return r.roomMemberRepo.Delete(ctx, roomId, memberId)
}

func (r *RoomService) UpdateMemberPoints(ctx context.Context, roomId uuid.UUID, memberId uuid.UUID, delta int, userId uuid.UUID) (*models.RoomMemberInfo, error) {
	room, err := r.roomRepo.GetById(ctx, roomId)
	if err != nil {
		return nil, err
	}

	if room.HostId != userId {
		return nil, apperrors.ErrForbidden
	}

	return r.roomMemberRepo.UpdatePoints(ctx, roomId, memberId, delta)
}

func (r *RoomService) Update(ctx context.Context, roomId uuid.UUID, name *string, categoryName *string, isPrivate *bool, thumbnailData *[]byte, thumbnailFilename *string, userId uuid.UUID) (*models.Room, error) {
	existingRoom, err := r.roomRepo.GetById(ctx, roomId)
	if err != nil {
		return nil, err
	}

	if existingRoom.HostId != userId {
		return nil, apperrors.ErrForbidden
	}

	if name != nil {
		existingRoom.Name = *name
	}
	if categoryName != nil {
		existingRoom.CategoryName = *categoryName
	}
	if isPrivate != nil {
		existingRoom.IsPrivate = *isPrivate
	}

	if thumbnailData != nil && thumbnailFilename != nil {
		if existingRoom.Thumbnail != "" && existingRoom.Thumbnail != defaultThumbnail {
			oldKey := fmt.Sprintf("rooms/%s/thumbnail", roomId.String())
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
		return nil, err
	}

	updatedRoom, err := r.roomRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedRoom, nil
}
