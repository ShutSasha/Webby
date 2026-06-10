package services

import (
	"context"
	"fmt"
	"log/slog"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type roomRetriever interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
}

type roomMemberRepository interface {
	Delete(ctx context.Context, roomID, userID uuid.UUID) error
	ListByRoom(ctx context.Context, roomID uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error)
	EnsureMember(ctx context.Context, roomID, userID uuid.UUID) error
	GetMemberPoints(ctx context.Context, roomID, userID uuid.UUID) (int, error)
}

type roomMemberChatManager interface {
	GetChatByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
	AddChatMember(ctx context.Context, chatID, userID uuid.UUID) error
}

type notificationSender interface {
	SendNotificationToUser(ctx context.Context, userID, roomID uuid.UUID) error
}

type roomMemberService struct {
	roomRetriever      roomRetriever
	roomMemberRepo     roomMemberRepository
	chatManager        roomMemberChatManager
	notificationClient notificationSender
}

func NewRoomMemberService(
	roomRepo roomRetriever,
	roomMemberRepo roomMemberRepository,
	chatClient roomMemberChatManager,
	notificationClient notificationSender,
) *roomMemberService {
	return &roomMemberService{
		roomRetriever:      roomRepo,
		roomMemberRepo:     roomMemberRepo,
		chatManager:        chatClient,
		notificationClient: notificationClient,
	}
}

func (svc *roomMemberService) AddMembers(ctx context.Context, roomID, hostID uuid.UUID, memberIDs []uuid.UUID) error {
	const op = "services.RoomMemberService.AddMembers"

	room, err := svc.roomRetriever.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != hostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	for _, memberID := range memberIDs {
		err := svc.roomMemberRepo.EnsureMember(ctx, roomID, memberID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if svc.chatManager != nil {
		chatID, err := svc.chatManager.GetChatByRoomID(ctx, roomID)
		if err != nil {
			slog.Warn("failed to get chat for room when adding members",
				slog.String("roomID", roomID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			for _, memberID := range memberIDs {
				err := svc.chatManager.AddChatMember(ctx, chatID, memberID)
				if err != nil {
					slog.Warn("failed to add room member as chat member",
						slog.String("roomID", roomID.String()),
						slog.String("chatID", chatID.String()),
						slog.String("memberID", memberID.String()),
						slog.String("error", err.Error()),
					)
				}

				err = svc.notificationClient.SendNotificationToUser(ctx, memberID, roomID)
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

func (svc *roomMemberService) ListMembers(
	ctx context.Context,
	roomID uuid.UUID,
	page, limit int,
	search string,
) ([]models.RoomMemberInfo, int64, error) {
	const op = "service.RoomMemberService.ListMembers"

	members, total, err := svc.roomMemberRepo.ListByRoom(ctx, roomID, page, limit, search)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return members, total, nil
}

func (svc *roomMemberService) RemoveMember(ctx context.Context, roomID, memberID, hostID uuid.UUID) error {
	const op = "services.RoomMemberService.RemoveMember"

	room, err := svc.roomRetriever.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != hostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrForbidden)
	}

	if memberID == room.HostID {
		return fmt.Errorf("%s: %w: cannot remove host from room", op, apperrors.ErrInvalidInput)
	}

	err = svc.roomMemberRepo.Delete(ctx, roomID, memberID)
	if err != nil {
		return fmt.Errorf("%s: %w: cannot remove room member", op, apperrors.ErrInvalidInput)
	}

	return nil
}

func (svc *roomMemberService) GetMemberPoints(ctx context.Context, roomID, userID uuid.UUID) (int, error) {
	const op = "services.roomMemberService.GetMemberPoints"

	points, err := svc.roomMemberRepo.GetMemberPoints(ctx, roomID, userID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return points, nil
}
