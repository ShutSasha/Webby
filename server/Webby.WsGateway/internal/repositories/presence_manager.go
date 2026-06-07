package repositories

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type presenceRepository struct {
	client *redis.Client
}

func NewPresenceRepository(client *redis.Client) *presenceRepository {
	return &presenceRepository{client}
}

func (p *presenceRepository) AddUser(ctx context.Context, chatID, userID string) error {
	const op = "repositories.presenceRepository.AddUser"

	pipe := p.client.Pipeline()
	pipe.SAdd(ctx, "chat:"+chatID+":presence", userID)
	pipe.SAdd(ctx, "active_chats", chatID)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *presenceRepository) RemoveUser(ctx context.Context, chatID, userID string) error {
	const op = "repositories.presenceRepository.RemoveUser"

	pipe := p.client.TxPipeline()

	pipe.SRem(ctx, "chat:"+chatID+":presence", userID)
	cardCmd := pipe.SCard(ctx, "chat:"+chatID+":presence")

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: tx exec %w", op, err)
	}

	if cardCmd.Val() == 0 {
		if err := p.client.SRem(ctx, "active_chats", chatID).Err(); err != nil {
			return fmt.Errorf("%s: removin g active chat %w", op, err)
		}
	}

	return nil
}
