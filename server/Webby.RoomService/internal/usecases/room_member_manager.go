package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type MemberManagerRoomRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
}

type MemberManagerRoomMemberRepository interface {
	Create(ctx context.Context, member *models.RoomMember) error
	EnsureMember(ctx context.Context, roomID, userID uuid.UUID) error
	Delete(ctx context.Context, roomID, userID uuid.UUID) error
	ListByRoom(
		ctx context.Context,
		roomID uuid.UUID,
		page, limit int,
		search string,
	) ([]models.RoomMemberInfo, int64, error)
}

type MemberManagerChatClient interface {
	GetChatByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
	AddChatMember(ctx context.Context, chatID, userID uuid.UUID) error
}

type MemberManagerNotificationClient interface {
	SendNotificationToUser(ctx context.Context, userID, roomID uuid.UUID) error
}

type RoomMemberManager struct {
	roomRepo           MemberManagerRoomRepository
	roomMemberRepo     MemberManagerRoomMemberRepository
	chatClient         MemberManagerChatClient
	notificationClient MemberManagerNotificationClient
}

func NewRoomMemberManager(
	roomRepo MemberManagerRoomRepository,
	roomMemberRepo MemberManagerRoomMemberRepository,
	chatClient MemberManagerChatClient,
	notificationClient MemberManagerNotificationClient,
) *RoomMemberManager {
	return &RoomMemberManager{
		roomRepo:           roomRepo,
		roomMemberRepo:     roomMemberRepo,
		chatClient:         chatClient,
		notificationClient: notificationClient,
	}
}

func (uc *RoomMemberManager) ExecuteAddMembers(ctx context.Context, roomID, hostID uuid.UUID, memberIDs []uuid.UUID) error {
	const op = "usecases.RoomMemberManager.ExecuteAddMembers"

	room, err := uc.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != hostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	for _, memberID := range memberIDs {
		if err := uc.roomMemberRepo.EnsureMember(ctx, roomID, memberID); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if uc.chatClient != nil {
		chatID, err := uc.chatClient.GetChatByRoomID(ctx, roomID)
		if err != nil {
			slog.Warn("failed to get chat for room when adding members",
				slog.String("roomID", roomID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			for _, memberID := range memberIDs {
				err := uc.chatClient.AddChatMember(ctx, chatID, memberID)
				if err != nil {
					slog.Warn("failed to add room member as chat member",
						slog.String("roomID", roomID.String()),
						slog.String("chatID", chatID.String()),
						slog.String("memberID", memberID.String()),
						slog.String("error", err.Error()),
					)
				}

				err = uc.notificationClient.SendNotificationToUser(ctx, memberID, roomID)
				if err != nil {
					slog.Warn("failed to notify member",
						slog.String("roomID", roomID.String()),
						slog.String("chatID", chatID.String()),
						slog.String("memberID", memberID.String()),
						slog.String("error", err.Error()),
					)
				}
			}
		}
	}

	return nil
}

func (uc *RoomMemberManager) ExecuteRemoveMember(ctx context.Context, roomID, memberID, hostID uuid.UUID) error {
	const op = "usecases.RoomMemberManager.ExecuteRemoveMember"

	room, err := uc.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != hostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	if memberID == room.HostID {
		return fmt.Errorf("%s: %w: cannot remove host from room", op, apperrors.ErrInvalidInput)
	}

	return uc.roomMemberRepo.Delete(ctx, roomID, memberID)
}

func (uc *RoomMemberManager) ExecuteListMembers(
	ctx context.Context,
	roomID uuid.UUID,
	page, limit int,
	search string,
) ([]models.RoomMemberInfo, int64, error) {
	const op = "usecases.RoomMemberManager.ExecuteListMembers"

	members, total, err := uc.roomMemberRepo.ListByRoom(ctx, roomID, page, limit, search)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return members, total, nil
}
