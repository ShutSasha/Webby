package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisPublisher struct {
	client *redis.Client
}

func NewRedisPublisher(client *redis.Client) *redisPublisher {
	return &redisPublisher{client: client}
}

func (p *redisPublisher) Publish(ctx context.Context, channel string, payload any) error {
	const op = "repositories.redisPublisher.Publish"

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = p.client.Publish(ctx, channel, string(data)).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
