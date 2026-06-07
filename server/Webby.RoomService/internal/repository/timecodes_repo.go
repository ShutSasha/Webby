package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type timecodesRepo struct {
	client *redis.Client
}

func NewTimecodesRepo(client *redis.Client) *timecodesRepo {
	return &timecodesRepo{client: client}
}

func (r *timecodesRepo) RetrieveTimecodes(ctx context.Context, roomID, syncID uuid.UUID) (map[string]int, error) {
	const op = "repository.RedisRepo.RetrieveTimecodes"

	key := fmt.Sprintf("room:%s:sync:%s", roomID, syncID)
	timecodesMap, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	timecodes := make(map[string]int, len(timecodesMap))
	for userID, timeStr := range timecodesMap {
		time, err := strconv.Atoi(timeStr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		timecodes[userID] = time
	}
	return timecodes, nil
}

func (r *timecodesRepo) SetTimecode(ctx context.Context, userID, roomID, syncID uuid.UUID, timecode int) error {
	const op = "repository.RedisRepo.SetTimecode"

	topic := fmt.Sprintf("room:%s:sync:%s", roomID.String(), syncID.String())
	err := r.client.HSet(ctx, topic, userID, timecode).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	r.client.Expire(ctx, topic, 10*time.Second)

	return nil
}
