package services

import (
	"context"
	"fmt"
	"time"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
)

type roomIDRetriever interface {
	RoomIDByChatIDBatch(ctx context.Context, chatIDs []string) (map[string]uuid.UUID, error)
}

type presenceRetriever interface {
	GetActiveRooms(ctx context.Context, minMembers int, zombieTTL time.Duration) ([]models.ActiveRoom, error)
}

type presenceService struct {
	presenceRetriever presenceRetriever
	roomIDRetriever   roomIDRetriever
}

func NewPresenceService(presenceRetriever presenceRetriever, roomIDRetriever roomIDRetriever) *presenceService {
	return &presenceService{
		presenceRetriever: presenceRetriever,
		roomIDRetriever:   roomIDRetriever,
	}
}

func (s *presenceService) GetActiveRooms(ctx context.Context, minMembers int, zombieTTL time.Duration) ([]models.ActiveRoom, error) {
	const op = "presenceService.GetActiveRooms"

	rooms, err := s.presenceRetriever.GetActiveRooms(ctx, minMembers, zombieTTL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chatIDs := make([]string, 0, len(rooms))
	for _, room := range rooms {
		chatIDs = append(chatIDs, room.ChatID)
	}

	chatIDroomIDMap, err := s.roomIDRetriever.RoomIDByChatIDBatch(ctx, chatIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i := range rooms {
		rooms[i].ID = chatIDroomIDMap[rooms[i].ChatID]
	}

	return rooms, nil
}
