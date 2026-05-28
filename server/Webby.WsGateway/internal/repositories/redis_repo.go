package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrTokenNotFound = errors.New("token not found or expired")

type TokenRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) *TokenRepository {
	return &TokenRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *TokenRepository) SaveToken(ctx context.Context, token string, userID uuid.UUID) error {
	const op = "repositories.TokenRepository.SaveToken"

	key := buildKey(token)
	err := r.client.Set(ctx, key, userID.String(), r.ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *TokenRepository) GetUserID(ctx context.Context, token string) (uuid.UUID, error) {
	const op = "repositories.TokenRepository.GetUserID"

	key := buildKey(token)
	val, err := r.client.GetDel(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, ErrTokenNotFound)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	parsedUUID, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: invalid uuid format in redis: %w", op, err)
	}

	return parsedUUID, nil
}

func buildKey(token string) string {
	return "session:" + token
}
