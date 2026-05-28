package services

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

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
		search, category string,
	) ([]models.Room, int64, error)
	ListPublic(
		ctx context.Context,
		page, limit int,
		search, category string,
	) ([]models.PublicRoom, int64, error)
	Update(ctx context.Context, room *models.Room) (uuid.UUID, error)
}

type RoomMemberRepository interface {
	Create(ctx context.Context, member *models.RoomMember) error
	Exists(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
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
	GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (uuid.UUID, error)
}

type CategoryClientInterface interface {
	Exists(ctx context.Context, name string) (bool, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type TimecodesRepo interface {
	RetrieveTimecodes(ctx context.Context, roomID, syncID uuid.UUID) (map[string]int, error)
	SetTimecode(ctx context.Context, userID, roomID, syncID uuid.UUID, timecode int) error
}

type RoomService struct {
	roomRepo       RoomRepository
	roomMemberRepo RoomMemberRepository
	fileRepo       FileRepository
	chatClient     ChatClientInterface
	categoryClient CategoryClientInterface
	publisher      EventPublisher
	timecodesRepo  TimecodesRepo
}

func NewRoomService(
	repo RoomRepository,
	roomMemberRepo RoomMemberRepository,
	fileRepo FileRepository,
	chatClient ChatClientInterface,
	categoryClient CategoryClientInterface,
	publisher EventPublisher,
	timecodesRepo TimecodesRepo,
) *RoomService {
	return &RoomService{
		roomRepo:       repo,
		roomMemberRepo: roomMemberRepo,
		fileRepo:       fileRepo,
		chatClient:     chatClient,
		categoryClient: categoryClient,
		publisher:      publisher,
		timecodesRepo:  timecodesRepo,
	}
}

func generateFileKey(filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("rooms/%d/thumbnail%s", time.Now().UnixMilli(), ext)
}
