package handlers

import (
	"context"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string) (*models.Room, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetByID(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (*models.Room, error)
	ListMy(ctx context.Context, userID uuid.UUID, page int, limit int, search string, category *string) ([]models.Room, int64, error)
	ListPublic(ctx context.Context, page int, limit int, search string, category *string) ([]models.PublicRoom, int64, error)
	ListMembers(ctx context.Context, roomID uuid.UUID, page int, limit int, search string) ([]models.RoomMemberInfo, int64, error)
	AddMembers(ctx context.Context, roomID uuid.UUID, memberIDs []uuid.UUID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, roomID uuid.UUID, memberID uuid.UUID, userID uuid.UUID) error
	Update(ctx context.Context, roomID uuid.UUID, name *string, category *string, isPrivate *bool, thumbnailData *[]byte, thumbnailFilename *string, userID uuid.UUID) (*models.Room, error)
}

type handler struct {
	service Service
}

func New(service Service) handler {
	return handler{service: service}
}
