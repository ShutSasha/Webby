package repository

import (
	"context"
	"fmt"
	"time"

	"webby/room-service/internal/workers"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type redisPresenceRepository struct {
	client *redis.Client
}

func NewRedisPresenceRepository(client *redis.Client) *redisPresenceRepository {
	return &redisPresenceRepository{client: client}
}

func (r *redisPresenceRepository) GetActiveRooms(ctx context.Context, minMembers int, zombieTTL time.Duration) ([]workers.ActiveRoom, error) {
	const op = "repository.redisPresenceRepository.GetActiveRooms"

	chats, err := r.client.SMembers(ctx, "active_chats").Result()
	if err != nil {
		return nil, fmt.Errorf("%s: get chats: %w", op, err)
	}

	var activeRooms []workers.ActiveRoom

	deadlineStr := fmt.Sprintf("%d", time.Now().Add(-zombieTTL).Unix())

	for _, chatIDStr := range chats {
		presenceKey := "chat:" + chatIDStr + ":presence"

		err := r.client.ZRemRangeByScore(ctx, presenceKey, "-inf", deadlineStr).Err()
		if err != nil {
			continue
		}

		count, err := r.client.ZCard(ctx, presenceKey).Result()
		if err != nil {
			continue
		}

		if count < int64(minMembers) {
			if count == 0 {
				r.client.SRem(ctx, "active_chats", chatIDStr)
			}
			continue
		}

		userIDsStr, err := r.client.ZRange(ctx, presenceKey, 0, -1).Result()
		if err != nil {
			continue
		}

		var userIDs []uuid.UUID
		for _, idStr := range userIDsStr {
			if id, err := uuid.Parse(idStr); err == nil {
				userIDs = append(userIDs, id)
			}
		}

		if len(userIDs) > 0 {
			activeRooms = append(activeRooms, workers.ActiveRoom{
				ChatID:  chatIDStr,
				UserIDs: userIDs,
			})
		}
	}

	return activeRooms, nil
}
