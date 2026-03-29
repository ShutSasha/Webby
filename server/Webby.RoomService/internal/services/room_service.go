package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"path/filepath"
	"webby/internal/apperrors"
	"webby/internal/models"

	"github.com/google/uuid"
)

const defaultThumbnail = "https://webby-watch-platform-bucket.s3.eu-north-1.amazonaws.com/rooms/default-room-preview.jpg"

type RoomRepository interface {
	Create(room *models.Room) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	GetById(id uuid.UUID) (*models.Room, error)
	GetByToken(token string) (*models.Room, error)
	ListMy(userId uuid.UUID, page int, limit int) ([]models.Room, int64, error)
	ListPublic(page int, limit int, search string, categoryId *uuid.UUID) ([]models.Room, int64, error)
	Update(room *models.Room) (uuid.UUID, error)
}

type RoomMemberRepository interface {
	Create(member *models.RoomMember) error
	Exists(roomId, userId uuid.UUID) (bool, error)
	EnsureMember(roomId, userId uuid.UUID) error
	Delete(roomId, userId uuid.UUID) error
	ListByRoom(roomId uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error)
}

type FileRepository interface {
	Save(ctx context.Context, key string, data []byte) (string, error)
	Remove(ctx context.Context, key string) error
}

type RoomService struct {
	roomRepo       RoomRepository
	roomMemberRepo RoomMemberRepository
	fileRepo       FileRepository
}

func NewRoomService(repo RoomRepository, roomMemberRepo RoomMemberRepository, fileRepo FileRepository) *RoomService {
	return &RoomService{
		roomRepo:       repo,
		roomMemberRepo: roomMemberRepo,
		fileRepo:       fileRepo,
	}
}

func generateFileKey(roomID uuid.UUID, filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("rooms/%s/thumbnail%s", roomID.String(), ext)
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (r *RoomService) Create(ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string) (*models.Room, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	room.Token = token

	id, err := r.roomRepo.Create(room)
	if err != nil {
		return nil, err
	}

	var thumbnailURL string
	if len(thumbnailData) > 0 && thumbnailFilename != "" {
		key := generateFileKey(id, thumbnailFilename)
		thumbnailURL, err = r.fileRepo.Save(ctx, key, thumbnailData)
		if err != nil {
			return nil, fmt.Errorf("thumbnail upload failed: %w", err)
		}
		room.Thumbnail = thumbnailURL
	} else {
		room.Thumbnail = defaultThumbnail
	}

	room.Id = id
	if _, err := r.roomRepo.Update(room); err != nil {
		return nil, fmt.Errorf("failed to update room thumbnail: %w", err)
	}

	if err := r.roomMemberRepo.EnsureMember(room.Id, room.HostId); err != nil {
		slog.Warn("failed to create room member for host",
			slog.String("roomId", room.Id.String()),
			slog.String("hostId", room.HostId.String()),
			slog.String("error", err.Error()),
		)
	}

	return room, nil
}

func (r *RoomService) Delete(ctx context.Context, id uuid.UUID, userId uuid.UUID) error {
	room, err := r.roomRepo.GetById(id)
	if err != nil {
		return err
	}

	if room.HostId != userId {
		return apperrors.ErrForbidden
	}

	if err := r.roomRepo.Delete(id); err != nil {
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
	room, err := r.roomRepo.GetById(roomId)
	if err != nil {
		return nil, err
	}

	if room.HostId != userId {
		return nil, apperrors.ErrForbidden
	}

	if err := r.roomMemberRepo.EnsureMember(room.Id, userId); err != nil {
		slog.Warn("failed to ensure room member on GetById",
			slog.String("roomId", room.Id.String()),
			slog.String("userId", userId.String()),
			slog.String("error", err.Error()),
		)
	}

	return room, nil
}

func (r *RoomService) GetByToken(ctx context.Context, token string, userId uuid.UUID) (*models.Room, error) {
	room, err := r.roomRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}

	if err := r.roomMemberRepo.EnsureMember(room.Id, userId); err != nil {
		slog.Warn("failed to ensure room member on GetByToken",
			slog.String("roomId", room.Id.String()),
			slog.String("userId", userId.String()),
			slog.String("error", err.Error()),
		)
	}

	return room, nil
}

func (r *RoomService) ListMy(userId uuid.UUID, page int, limit int) ([]models.Room, int64, error) {
	rooms, total, err := r.roomRepo.ListMy(userId, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *RoomService) ListPublic(page int, limit int, search string, categoryId *uuid.UUID) ([]models.Room, int64, error) {
	rooms, total, err := r.roomRepo.ListPublic(page, limit, search, categoryId)
	if err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *RoomService) ListMembers(roomId uuid.UUID, page int, limit int, search string) ([]models.RoomMemberInfo, int64, error) {
	members, total, err := r.roomMemberRepo.ListByRoom(roomId, page, limit, search)
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

func (r *RoomService) AddMember(ctx context.Context, roomId uuid.UUID, memberId uuid.UUID, userId uuid.UUID) error {
	room, err := r.roomRepo.GetById(roomId)
	if err != nil {
		return err
	}

	if room.HostId != userId {
		return apperrors.ErrForbidden
	}

	return r.roomMemberRepo.EnsureMember(roomId, memberId)
}

func (r *RoomService) RemoveMember(ctx context.Context, roomId uuid.UUID, memberId uuid.UUID, userId uuid.UUID) error {
	room, err := r.roomRepo.GetById(roomId)
	if err != nil {
		return err
	}

	if room.HostId != userId {
		return apperrors.ErrForbidden
	}

	if memberId == room.HostId {
		return fmt.Errorf("%w: cannot remove the host from the room", apperrors.ErrInvalidInput)
	}

	return r.roomMemberRepo.Delete(roomId, memberId)
}

func (r *RoomService) Update(ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string, userId uuid.UUID) (uuid.UUID, error) {
	existingRoom, err := r.roomRepo.GetById(room.Id)
	if err != nil {
		return uuid.Nil, err
	}

	if existingRoom.HostId != userId {
		return uuid.Nil, apperrors.ErrForbidden
	}

	if len(thumbnailData) > 0 && thumbnailFilename != "" {
		if existingRoom.Thumbnail != "" && existingRoom.Thumbnail != defaultThumbnail {
			oldKey := fmt.Sprintf("rooms/%s/thumbnail", room.Id.String())
			_ = r.fileRepo.Remove(ctx, oldKey)
		}

		key := generateFileKey(room.Id, thumbnailFilename)
		thumbnailURL, err := r.fileRepo.Save(ctx, key, thumbnailData)
		if err != nil {
			return uuid.Nil, fmt.Errorf("thumbnail upload failed: %w", err)
		}
		room.Thumbnail = thumbnailURL
	} else if room.Thumbnail == "" && existingRoom.Thumbnail != defaultThumbnail {
		oldKey := fmt.Sprintf("rooms/%s/thumbnail", room.Id.String())
		_ = r.fileRepo.Remove(ctx, oldKey)
		room.Thumbnail = defaultThumbnail
	}

	id, err := r.roomRepo.Update(room)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
