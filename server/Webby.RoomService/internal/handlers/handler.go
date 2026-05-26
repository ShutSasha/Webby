package handlers

import (
	"context"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string) (*models.Room, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetByID(ctx context.Context, roomID, userID uuid.UUID) (*models.Room, error)
	ListMy(ctx context.Context, userID uuid.UUID, page, limit int, search, category string) ([]models.Room, int64, error)
	ListPublic(ctx context.Context, page, limit int, search, category string) ([]models.PublicRoom, int64, error)
	ListMembers(ctx context.Context, roomID uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error)
	AddMembers(ctx context.Context, roomID, userID uuid.UUID, memberIDs []uuid.UUID) error
	RemoveMember(ctx context.Context, roomID, memberID, hostID uuid.UUID) error
	Update(ctx context.Context, roomID, userID uuid.UUID, name, category, thumbnailFilename *string, thumbnailData *[]byte, isPrivate *bool) (*models.Room, error)
	Synchronize(ctx context.Context, userID, roomID uuid.UUID) error
	ReportTimecode(ctx context.Context, userID, roomID, syncID uuid.UUID, timecode int) error
}

type handler struct {
	service Service
}

func New(service Service) handler {
	return handler{service: service}
}
