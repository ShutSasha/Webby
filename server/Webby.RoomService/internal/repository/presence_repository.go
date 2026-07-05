package repository

import (
	"context"
	"fmt"
	"time"

	"webby/room-service/internal/models"
	"webby/room-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type redisPresenceRepository struct {
	client *redis.Client
}

func NewRedisPresenceRepository(client *redis.Client) *redisPresenceRepository {
	return &redisPresenceRepository{client: client}
}

func (r *redisPresenceRepository) GetActiveRooms(ctx context.Context, minMembers int, zombieTTL time.Duration) ([]models.ActiveRoom, error) {
	const op = "repository.redisPresenceRepository.GetActiveRooms"
	log := logger.FromContext(ctx).With("op", op)

	chats, err := r.client.SMembers(ctx, "active_chats").Result()
	if err != nil {
		return nil, fmt.Errorf("%s: get chats: %w", op, err)
	}

	var activeRooms []models.ActiveRoom

	deadlineStr := fmt.Sprintf("%d", time.Now().Add(-zombieTTL).Unix())

	for _, chatIDStr := range chats {
		presenceKey := "chat:" + chatIDStr + ":presence"

		err := r.client.ZRemRangeByScore(ctx, presenceKey, "-inf", deadlineStr).Err()
		if err != nil {
			log.Error("failed to ZRemRangeByScore", "err", err.Error())
			continue
		}

		count, err := r.client.ZCard(ctx, presenceKey).Result()
		if err != nil {
			log.Error("failed to ZCard", "err", err.Error())
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
			log.Error("failed to ZRange", "err", err.Error())
			continue
		}

		var userIDs []uuid.UUID
		for _, idStr := range userIDsStr {
			id, err := uuid.Parse(idStr)
			if err != nil {
				log.Error("failed to uuid.Parse", "err", err.Error())
				continue
			}
			userIDs = append(userIDs, id)
		}

		if len(userIDs) > 0 {
			activeRooms = append(activeRooms, models.ActiveRoom{
				ChatID:  chatIDStr,
				UserIDs: userIDs,
			})
		}
	}

	return activeRooms, nil
}
