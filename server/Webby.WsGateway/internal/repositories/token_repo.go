package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"
	"webby/wsgateway/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type tokenRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewTokenRepository(client *redis.Client, ttl time.Duration) *tokenRepository {
	return &tokenRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *tokenRepository) SaveToken(ctx context.Context, token string, userID uuid.UUID) error {
	const op = "repositories.tokenRepository.SaveToken"

	key := buildKey(token)
	err := r.client.Set(ctx, key, userID.String(), r.ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *tokenRepository) GetUserID(ctx context.Context, token string) (uuid.UUID, error) {
	const op = "repositories.tokenRepository.GetUserID"

	key := buildKey(token)
	val, err := r.client.GetDel(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, domain.ErrTokenNotFound)
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
