package handlers

import (
	"context"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type roomService interface {
	Create(ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string) (*models.Room, error)
	GetDetails(ctx context.Context, roomID, userID uuid.UUID) (*models.Room, error)
	ListMyRooms(ctx context.Context, userID uuid.UUID, page, limit int, search, category string) ([]models.Room, int64, error)
	ListPublicRooms(ctx context.Context, page, limit int, search, category string) ([]models.PublicRoom, int64, error)
	Update(ctx context.Context, roomID, userID uuid.UUID, name, category, thumbnailFilename *string, thumbnailData *[]byte, isPrivate *bool) (*models.Room, error)
	Delete(ctx context.Context, roomID, userID uuid.UUID) error
}

type roomMemberService interface {
	AddMembers(ctx context.Context, roomID, hostID uuid.UUID, memberIDs []uuid.UUID) error
	ListMembers(ctx context.Context, roomID uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error)
	RemoveMember(ctx context.Context, roomID, memberID, hostID uuid.UUID) error
	GetMemberPoints(ctx context.Context, roomID, userID uuid.UUID) (int, error)
}

type synchronizeService interface {
	Synchronize(ctx context.Context, userID, roomID uuid.UUID) error
	ReportTimecode(ctx context.Context, userID, roomID, syncID uuid.UUID, timecode int) error
}

type handler struct {
	roomService        roomService
	roomMemberService  roomMemberService
	synchronizeService synchronizeService
}

func New(
	roomService roomService,
	roomMemberService roomMemberService,
	synchronizeService synchronizeService,
) handler {
	return handler{
		roomService:        roomService,
		roomMemberService:  roomMemberService,
		synchronizeService: synchronizeService,
	}
}
