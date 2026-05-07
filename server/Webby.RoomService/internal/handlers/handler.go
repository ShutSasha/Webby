package handlers

import (
	"context"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string) (*models.Room, error)
	Delete(ctx context.Context, id uuid.UUID, userId uuid.UUID) error
	GetById(ctx context.Context, roomId uuid.UUID, userId uuid.UUID) (*models.Room, error)
	ListMy(ctx context.Context, userId uuid.UUID, page int, limit int, search string, categoryName *string) ([]models.Room, int64, error)
	ListPublic(ctx context.Context, page int, limit int, search string, categoryName *string) ([]models.PublicRoom, int64, error)
	ListMembers(ctx context.Context, roomId uuid.UUID, page int, limit int, search string) ([]models.RoomMemberInfo, int64, error)
	AddMembers(ctx context.Context, roomId uuid.UUID, memberIds []uuid.UUID, userId uuid.UUID) error
	RemoveMember(ctx context.Context, roomId uuid.UUID, memberId uuid.UUID, userId uuid.UUID) error
	UpdateMemberPoints(ctx context.Context, roomId uuid.UUID, memberId uuid.UUID, delta int, userId uuid.UUID) (*models.RoomMemberInfo, error)
	Update(ctx context.Context, roomId uuid.UUID, name *string, categoryName *string, isPrivate *bool, thumbnailData *[]byte, thumbnailFilename *string, userId uuid.UUID) (*models.Room, error)
}

type handler struct {
	service Service
}

func New(service Service) handler {
	return handler{service: service}
}
