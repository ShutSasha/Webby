package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"webby/vote-service/internal/apperrors"
	"webby/vote-service/internal/models"
	"webby/vote-service/pkg/logger"

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
	roomIdxKey := fmt.Sprintf("room:%s:votings", vote.RoomID)

	pipe := r.client.TxPipeline()

	fields := map[string]any{
		"id":         vote.ID.String(),
		"room_id":    vote.RoomID.String(),
		"vote_text":  vote.VoteText,
		"duration":   vote.Duration,
		"created_at": vote.CreatedAt.Format(time.RFC3339),
		"status":     "active",
	}

	pipe.HSet(ctx, key, fields)
	pipe.Expire(ctx, key, 24*time.Hour)

	pipe.SAdd(ctx, roomIdxKey, vote.ID.String())
	pipe.Expire(ctx, roomIdxKey, 24*time.Hour)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("%s: pipeline failed: %w", op, err)
	}

	return nil
}

func (r *repository) SaveChoices(ctx context.Context, voteID uuid.UUID, choices []string) error {
	const op = "repository.SaveChoices"
	key := fmt.Sprintf("votings:%s:options", voteID)

	interfaceChoices := make([]any, len(choices))
	for i, v := range choices {
		interfaceChoices[i] = v
	}

	pipe := r.client.TxPipeline()
	pipe.SAdd(ctx, key, interfaceChoices...)
	pipe.Expire(ctx, key, 24*time.Hour)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("%s: pipeline failed: %w", op, err)
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

func (r *repository) SetVotingRightOption(ctx context.Context, roomID, voteID uuid.UUID, rightChoice string) error {
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

			optionsKey := fmt.Sprintf("votings:%s:options", voteID)
			votesKey := fmt.Sprintf("votings:%s:user_choices", voteID)
			roomIdxKey := fmt.Sprintf("room:%s:votings", roomID)

			pipe.Expire(ctx, key, 5*time.Minute)
			pipe.Expire(ctx, optionsKey, 5*time.Minute)
			pipe.Expire(ctx, votesKey, 5*time.Minute)
			pipe.SRem(ctx, roomIdxKey, voteID.String())

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

func (r *repository) GetVoteByID(ctx context.Context, id uuid.UUID) (*models.Vote, error) {
	const op = "repository.GetVoteByID"
	log := logger.FromContext(ctx).With("op", op)

	key := fmt.Sprintf("votings:%s", id)

	res, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("%s: vote not found", op)
	}

	parseUUID := func(val string) uuid.UUID {
		u, parseErr := uuid.Parse(val)
		if parseErr != nil {
			u, _ = uuid.FromBytes([]byte(val))
		}
		return u
	}

	duration, _ := strconv.Atoi(res["duration"])
	createdAt, err := time.Parse(time.RFC3339, res["created_at"])
	if err != nil {
		createdAt.UnmarshalText([]byte(res["created_at"]))
	}

	vote := &models.Vote{
		ID:        parseUUID(res["id"]),
		RoomID:    parseUUID(res["room_id"]),
		VoteText:  res["vote_text"],
		Duration:  duration,
		CreatedAt: createdAt,
		Status:    res["status"],
	}

	log.Debug("Retrieve vote from redis", "vote", vote)
	if vote.ID == uuid.Nil {
		return nil, fmt.Errorf("%s: vote not found", op)
	}

	return vote, nil
}

func (r *repository) ListByRoom(ctx context.Context, roomID uuid.UUID) ([]models.Vote, error) {
	const op = "repository.ListByRoom"
	log := logger.FromContext(ctx).With("op", op)

	roomIdxKey := fmt.Sprintf("room:%s:votings", roomID)

	voteIDs, err := r.client.SMembers(ctx, roomIdxKey).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: get room votings index: %w", op, err)
	}

	var votes []models.Vote
	for _, vIDStr := range voteIDs {
		vID, parseErr := uuid.Parse(vIDStr)
		if parseErr != nil {
			log.Warn("Failed to parse vote ID, removing from set", "vIDStr", vIDStr)
			r.client.SRem(ctx, roomIdxKey, vIDStr)
			continue
		}

		vote, err := r.GetVoteByID(ctx, vID)
		if err != nil {
			if errors.Is(err, redis.Nil) || strings.Contains(err.Error(), "not found") {
				log.Info("Vote body expired naturally, cleaning up orphaned ID", "voteID", vID)
				r.client.SRem(ctx, roomIdxKey, vIDStr)
			} else {
				log.Error("Failed to fetch vote", "voteID", vID, "error", err.Error())
			}
			continue
		}

		if vote != nil {
			votes = append(votes, *vote)
		}
	}

	return votes, nil
}

func (r *repository) DeleteVote(ctx context.Context, id uuid.UUID) error {
	const op = "repository.DeleteVote"

	vote, err := r.GetVoteByID(ctx, id)
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

func (r *repository) GetChoicesForVoting(ctx context.Context, voteID uuid.UUID) ([]string, error) {
	const op = "repository.GetChoicesForVoting"
	optionsKey := fmt.Sprintf("votings:%s:options", voteID)

	choices, err := r.client.SMembers(ctx, optionsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return choices, nil
}

func (r *repository) CastVote(ctx context.Context, voteID, userID uuid.UUID, choice string) error {
	const op = "repository.CastVote"
	key := fmt.Sprintf("votings:%s", voteID)
	votesKey := fmt.Sprintf("votings:%s:user_choices", voteID)

	err := r.client.Watch(ctx, func(tx *redis.Tx) error {
		status, err := tx.HGet(ctx, key, "status").Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return fmt.Errorf("voting not found")
			}
			return err
		}

		if status != "active" {
			return apperrors.ErrVotingLocked
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, votesKey, userID.String(), choice)
			return nil
		})
		return err
	}, key)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) CreateVotingForNextVideo(ctx context.Context, roomID uuid.UUID, duration int) error {
	const op = "repository.CreateVotingForNextVideo"
	nextVotingsKey := "rooms:nextvideo"
	votingKey := fmt.Sprintf("rooms:nextvideo:%s", roomID.String())

	err := r.client.SAdd(ctx, nextVotingsKey, roomID.String()).Err()
	if err != nil {
		return fmt.Errorf("%s: add voting: %w", op, err)
	}

	err = r.client.HSet(ctx, votingKey, "duration", duration, "created_at", time.Now()).Err()
	if err != nil {
		return fmt.Errorf("%s: create voting: %w", op, err)
	}

	return nil
}

func (r *repository) GetNextVideoResults(ctx context.Context, roomID uuid.UUID) (map[uuid.UUID]int, error) {
	const op = "repository.GetNextVideoResults"
	log := logger.FromContext(ctx).With("op", op)

	votesKey := fmt.Sprintf("votings:%s:user_choices", roomID)
	activeVotingsKey := "rooms:nextvideo"

	allVotes, err := r.client.HGetAll(ctx, votesKey).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results := make(map[uuid.UUID]int)
	for userID, choice := range allVotes {
		choiceID, parseErr := uuid.Parse(choice)
		if parseErr != nil {
			log.Warn("Failed to parse choice ID, removing from set", "choice", choice)
			r.client.HDel(ctx, votesKey, userID)
			continue
		}
		results[choiceID]++
	}

	r.client.Expire(ctx, votesKey, 30*time.Second)
	r.client.SRem(ctx, activeVotingsKey, roomID.String())

	return results, nil
}

func (r *repository) VoteForNextVideo(ctx context.Context, roomID, userID, queueItemID uuid.UUID) error {
	const op = "repository.VoteForNextVideo"

	votesKey := fmt.Sprintf("votings:%s:user_choices", roomID)
	err := r.client.HSet(ctx, votesKey, userID.String(), queueItemID.String()).Err()
	if err != nil {
		return fmt.Errorf("%s: cast vote: %w", op, err)
	}

	return nil
}

func (r *repository) GetNextVideoVoting(ctx context.Context, roomID uuid.UUID) (*models.NextVideoInfo, error) {
	const op = "repository.GetNextVideoVoting"
	nextVotingsKey := "rooms:nextvideo"
	votingKey := fmt.Sprintf("rooms:nextvideo:%s", roomID.String())

	hasVoting, err := r.client.SIsMember(ctx, nextVotingsKey, roomID.String()).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !hasVoting {
		return &models.NextVideoInfo{
			Exists: false,
		}, nil
	}

	durationStr, err := r.client.HGet(ctx, votingKey, "duration").Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	duration, _ := strconv.Atoi(durationStr)

	createdAtStr, err := r.client.HGet(ctx, votingKey, "created_at").Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	createdAt, err := time.Parse(time.RFC3339Nano, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("%s: parse created_at: %w", op, err)
	}

	expiresAt := createdAt.Add(time.Duration(duration) * time.Second)
	return &models.NextVideoInfo{
		Exists:    true,
		Duration:  &duration,
		ExpiresAt: &expiresAt,
	}, nil
}

func (r *repository) GetUserVote(ctx context.Context, voteID, userID uuid.UUID) (*string, error) {
	const op = "repository.GetUserVote"
	votesKey := fmt.Sprintf("votings:%s:user_choices", voteID)

	choice, err := r.client.HGet(ctx, votesKey, userID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &choice, nil
}
