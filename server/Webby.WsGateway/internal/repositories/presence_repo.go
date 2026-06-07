package repositories

import (
	"context"
	"fmt"
	"time"

	"webby/wsgateway/internal/ws"

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
	pipe.ZAdd(ctx, "chat:"+chatID+":presence", redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: userID,
	})
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

	pipe.ZRem(ctx, "chat:"+chatID+":presence", userID)
	cardCmd := pipe.ZCard(ctx, "chat:"+chatID+":presence")

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: tx exec %w", op, err)
	}

	if cardCmd.Val() == 0 {
		if err := p.client.SRem(ctx, "active_chats", chatID).Err(); err != nil {
			return fmt.Errorf("%s: removing active chat %w", op, err)
		}
	}

	return nil
}

func (p *presenceRepository) UpdateHeartbeats(ctx context.Context, sessions []ws.SessionInfo) error {
	const op = "repositories.presenceRepository.UpdateHeartbeats"

	if len(sessions) == 0 {
		return nil
	}

	pipe := p.client.Pipeline()
	now := float64(time.Now().Unix())

	for _, sess := range sessions {
		pipe.ZAdd(ctx, "chat:"+sess.ChatID+":presence", redis.Z{
			Score:  now,
			Member: sess.UserID,
		})
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
