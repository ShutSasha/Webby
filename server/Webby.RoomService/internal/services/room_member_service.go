package services

import (
	"context"
	"fmt"
	"log/slog"

	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

const MaxRoomMembers = 20

type roomRetriever interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
}

type roomMemberRepository interface {
	ListByRoom(ctx context.Context, roomID uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error)
	EnsureMember(ctx context.Context, roomID, userID uuid.UUID) error
	GetMemberStatus(ctx context.Context, roomID, userID uuid.UUID) (*models.MemberStatus, error)
	GetMemberPoints(ctx context.Context, roomID, userID uuid.UUID) (int, error)
	BanRoomMember(ctx context.Context, roomID, userID uuid.UUID) error
	UnbanRoomMember(ctx context.Context, roomID, userID uuid.UUID) error
	CountMembers(ctx context.Context, roomID uuid.UUID) (int, error)
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

func (s *roomMemberService) AddMembers(ctx context.Context, roomID, hostID uuid.UUID, memberIDs []uuid.UUID) error {
	const op = "services.roomMemberService.AddMembers"

	room, err := s.roomRetriever.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != hostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNotHost)
	}

	actualMemberCount, err := s.roomMemberRepo.CountMembers(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if MaxRoomMembers-actualMemberCount < len(memberIDs) {
		return fmt.Errorf("%s: %w", op, apperrors.ErrMaxMembersReached)
	}

	chatID, err := s.chatManager.GetChatByRoomID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, memberID := range memberIDs {
		stats, err := s.roomMemberRepo.GetMemberStatus(ctx, roomID, memberID)
		if err != nil {
			slog.Error("failed to get member status",
				slog.String("memberID", memberID.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		if stats.IsBanned {
			err := s.roomMemberRepo.UnbanRoomMember(ctx, roomID, memberID)
			if err != nil {
				slog.Error("failed to unban member",
					slog.String("memberID", memberID.String()),
					slog.String("error", err.Error()),
				)
				continue
			}
		}

		err = s.roomMemberRepo.EnsureMember(ctx, roomID, memberID)
		if err != nil {
			slog.Error("failed to ensure member",
				slog.String("memberID", memberID.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		err = s.chatManager.AddChatMember(ctx, chatID, memberID)
		if err != nil {
			slog.Warn("failed to add room member as chat member",
				slog.String("memberID", memberID.String()),
				slog.String("error", err.Error()),
			)
		}

		err = s.notificationClient.SendNotificationToUser(ctx, memberID, roomID)
		if err != nil {
			slog.Warn("failed to notify member",
				slog.String("memberID", memberID.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	return nil
}

func (s *roomMemberService) ListMembers(ctx context.Context, roomID uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error) {
	const op = "service.RoomMemberService.ListMembers"

	members, total, err := s.roomMemberRepo.ListByRoom(ctx, roomID, page, limit, search)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return members, total, nil
}

func (s *roomMemberService) RemoveMember(ctx context.Context, roomID, memberID, hostID uuid.UUID) error {
	const op = "services.roomMemberService.RemoveMember"

	room, err := s.roomRetriever.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if room.HostID != hostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNotHost)
	}

	if memberID == room.HostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrRemoveHost)
	}

	err = s.roomMemberRepo.BanRoomMember(ctx, roomID, memberID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *roomMemberService) GetMemberPoints(ctx context.Context, roomID, userID uuid.UUID) (int, error) {
	const op = "services.roomMemberService.GetMemberPoints"

	points, err := s.roomMemberRepo.GetMemberPoints(ctx, roomID, userID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return points, nil
}
