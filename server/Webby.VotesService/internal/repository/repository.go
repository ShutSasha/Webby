package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/vote-service/internal/apperrors"
	"webby/vote-service/internal/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *repository {
	return &repository{client}
}

func (r *repository) CreateVoteWithRightChoice(ctx context.Context, vote *models.Vote) error {
	const op = "repository.CreateVoteWithRightChoice"

	key := fmt.Sprintf("votings:%s", vote.ID)

	pipe := r.client.TxPipeline()
	pipe.HSet(ctx, key, vote)
	pipe.HSet(ctx, key, "status", "active")

	roomIdxKey := fmt.Sprintf("room:%s:votings", vote.RoomID)
	pipe.SAdd(ctx, roomIdxKey, vote.ID.String())

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("%s: pipeline failed: %w", op, err)
	}

	return nil
}

func (r *repository) SaveChoices(ctx context.Context, voteID uuid.UUID, choices []string) error {
	const op = "repository.SaveChoices"

	interfaceChoices := make([]any, len(choices))
	for i, v := range choices {
		interfaceChoices[i] = v
	}

	key := fmt.Sprintf("votings:%s:options", voteID)
	cmd := r.client.SAdd(ctx, key, interfaceChoices...)
	if cmd.Err() != nil {
		return fmt.Errorf("%s: %w", op, cmd.Err())
	}

	return nil
}

func (r *repository) IsChoiceValid(ctx context.Context, voteID uuid.UUID, choice string) (bool, error) {
	const op = "repository.IsChoiceValid"
	key := fmt.Sprintf("votings:%s:options", voteID)

	isValid, err := r.client.SIsMember(ctx, key, choice).Result()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isValid, nil
}

func (r *repository) SetVotingRightOption(ctx context.Context, voteID uuid.UUID, rightChoice string) error {
	const op = "repository.SetVotingRightOption"
	key := fmt.Sprintf("votings:%s", voteID)

	err := r.client.Watch(ctx, func(tx *redis.Tx) error {
		status, err := tx.HGet(ctx, key, "status").Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}

		if status == "resolved" {
			return apperrors.ErrAlreadyClosed
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, key, "status", "resolved")
			pipe.HSet(ctx, key, "right_choice", rightChoice)
			return nil
		})
		return err
	}, key)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *repository) MarkVoteAsLocked(ctx context.Context, voteID uuid.UUID) error {
	const op = "repository.MarkVoteAsLocked"
	key := fmt.Sprintf("votings:%s", voteID)

	err := r.client.Watch(ctx, func(tx *redis.Tx) error {
		status, err := tx.HGet(ctx, key, "status").Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}

		if status != "active" {
			return nil
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, key, "status", "locked")
			return nil
		})
		return err
	}, key)

	return err
}

func (r *repository) GetVotingUserWinners(ctx context.Context, voteID uuid.UUID, rightChoice string) ([]uuid.UUID, error) {
	const op = "repository.GetVotingUserWinners"
	votesKey := fmt.Sprintf("votings:%s:user_choices", voteID)

	allVotes, err := r.client.HGetAll(ctx, votesKey).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var winners []uuid.UUID
	for userIDStr, choice := range allVotes {
		if choice == rightChoice {
			if parsedID, parseErr := uuid.Parse(userIDStr); parseErr == nil {
				winners = append(winners, parsedID)
			}
		}
	}

	return winners, nil
}

func (r *repository) GetVotingRightOption(ctx context.Context, voteID uuid.UUID) (string, error) {
	const op = "repository.GetVotingRightOption"
	key := fmt.Sprintf("votings:%s", voteID)

	rightChoice, err := r.client.HGet(ctx, key, "right_choice").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("%s: right choice not set yet", op)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return rightChoice, nil
}

func (r *repository) GetVoteById(ctx context.Context, id uuid.UUID) (*models.Vote, error) {
	const op = "repository.GetVoteById"
	key := fmt.Sprintf("votings:%s", id)

	var vote models.Vote
	err := r.client.HGetAll(ctx, key).Scan(&vote)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if vote.ID == uuid.Nil {
		return nil, fmt.Errorf("%s: vote not found", op)
	}

	return &vote, nil
}

func (r *repository) ListByRoom(ctx context.Context, roomID uuid.UUID) ([]models.Vote, error) {
	const op = "repository.ListByRoom"
	roomIdxKey := fmt.Sprintf("room:%s:votings", roomID)

	voteIDs, err := r.client.SMembers(ctx, roomIdxKey).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: get room votings index: %w", op, err)
	}

	var votes []models.Vote
	for _, vIDStr := range voteIDs {
		vID, parseErr := uuid.Parse(vIDStr)
		if parseErr != nil {
			continue
		}

		vote, err := r.GetVoteById(ctx, vID)
		if err == nil && vote != nil {
			votes = append(votes, *vote)
		}
	}

	return votes, nil
}

func (r *repository) DeleteVote(ctx context.Context, id uuid.UUID) error {
	const op = "repository.DeleteVote"

	vote, err := r.GetVoteById(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	key := fmt.Sprintf("votings:%s", id)
	optionsKey := fmt.Sprintf("votings:%s:options", id)
	userChoicesKey := fmt.Sprintf("votings:%s:user_choices", id)
	roomIdxKey := fmt.Sprintf("room:%s:votings", vote.RoomID)

	pipe := r.client.TxPipeline()
	pipe.Del(ctx, key, optionsKey, userChoicesKey)
	pipe.SRem(ctx, roomIdxKey, id.String())

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
