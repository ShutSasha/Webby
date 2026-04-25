package repository

import (
	"context"
	"errors"
	"fmt"
	"webby-vote-service/internal/apperrors"
	"webby-vote-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository struct {
	db *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) CreateVote(
	ctx context.Context, vote *models.Vote,
) (uuid.UUID, error) {
	const op = "repository.VoteRepository.CreateVote"

	if vote == nil {
		return uuid.Nil, fmt.Errorf(
			"%s: %w: vote cannot be nil", op, apperrors.ErrInvalidInput,
		)
	}

	vote.Id = uuid.New()

	query := `
		INSERT INTO votes (id, room_id, type, vote_text, duration_seconds)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	err := r.db.QueryRow(
		ctx, query,
		vote.Id, vote.RoomId, vote.Type, vote.VoteText, vote.DurationSeconds,
	).Scan(&vote.CreatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return vote.Id, nil
}

func (r *VoteRepository) CreateChoice(
	ctx context.Context, choice *models.VoteChoice,
) (uuid.UUID, error) {
	const op = "repository.VoteRepository.CreateChoice"

	if choice == nil {
		return uuid.Nil, fmt.Errorf(
			"%s: %w: choice cannot be nil", op, apperrors.ErrInvalidInput,
		)
	}

	choice.Id = uuid.New()

	query := `
		INSERT INTO vote_choices (id, vote_id, name, is_correct, queue_item_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx, query,
		choice.Id, choice.VoteId, choice.Name, choice.IsCorrect, choice.QueueItemId,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return choice.Id, nil
}

func (r *VoteRepository) DeleteVote(
	ctx context.Context, id uuid.UUID,
) error {
	const op = "repository.VoteRepository.DeleteVote"

	if id == uuid.Nil {
		return fmt.Errorf(
			"%s: %w: invalid vote id", op, apperrors.ErrInvalidInput,
		)
	}

	query := `DELETE FROM votes WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf(
			"%s: vote %s: %w", op, id.String(), apperrors.ErrNotFound,
		)
	}

	return nil
}

func (r *VoteRepository) GetVoteById(
	ctx context.Context, id uuid.UUID,
) (*models.Vote, error) {
	const op = "repository.VoteRepository.GetVoteById"

	if id == uuid.Nil {
		return nil, fmt.Errorf(
			"%s: %w: invalid vote id", op, apperrors.ErrInvalidInput,
		)
	}

	query := `
		SELECT id, room_id, type, vote_text, created_at, duration_seconds
		FROM votes
		WHERE id = $1
	`

	var vote models.Vote
	err := r.db.QueryRow(ctx, query, id).Scan(
		&vote.Id,
		&vote.RoomId,
		&vote.Type,
		&vote.VoteText,
		&vote.CreatedAt,
		&vote.DurationSeconds,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(
				"%s: vote %s: %w", op, id.String(), apperrors.ErrNotFound,
			)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &vote, nil
}

func (r *VoteRepository) ListByRoom(
	ctx context.Context, roomId uuid.UUID,
) ([]models.Vote, error) {
	const op = "repository.VoteRepository.ListByRoom"

	if roomId == uuid.Nil {
		return nil, fmt.Errorf(
			"%s: %w: invalid room id", op, apperrors.ErrInvalidInput,
		)
	}

	query := `
		SELECT id, room_id, type, vote_text, created_at, duration_seconds
		FROM votes
		WHERE room_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, roomId)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	votes := []models.Vote{}
	for rows.Next() {
		var vote models.Vote
		err := rows.Scan(
			&vote.Id,
			&vote.RoomId,
			&vote.Type,
			&vote.VoteText,
			&vote.CreatedAt,
			&vote.DurationSeconds,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		votes = append(votes, vote)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return votes, nil
}

func (r *VoteRepository) GetChoicesByVoteId(
	ctx context.Context, voteId uuid.UUID,
) ([]models.VoteChoice, error) {
	const op = "repository.VoteRepository.GetChoicesByVoteId"

	if voteId == uuid.Nil {
		return nil, fmt.Errorf(
			"%s: %w: invalid vote id", op, apperrors.ErrInvalidInput,
		)
	}

	query := `
		SELECT id, vote_id, name, votes, is_correct, queue_item_id
		FROM vote_choices
		WHERE vote_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query, voteId)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	choices := []models.VoteChoice{}
	for rows.Next() {
		var choice models.VoteChoice
		err := rows.Scan(
			&choice.Id,
			&choice.VoteId,
			&choice.Name,
			&choice.Votes,
			&choice.IsCorrect,
			&choice.QueueItemId,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		choices = append(choices, choice)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return choices, nil
}

func (r *VoteRepository) GetChoiceById(
	ctx context.Context, choiceId uuid.UUID,
) (*models.VoteChoice, error) {
	const op = "repository.VoteRepository.GetChoiceById"

	if choiceId == uuid.Nil {
		return nil, fmt.Errorf(
			"%s: %w: invalid choice id", op, apperrors.ErrInvalidInput,
		)
	}

	query := `
		SELECT id, vote_id, name, votes, is_correct, queue_item_id
		FROM vote_choices
		WHERE id = $1
	`

	var choice models.VoteChoice
	err := r.db.QueryRow(ctx, query, choiceId).Scan(
		&choice.Id,
		&choice.VoteId,
		&choice.Name,
		&choice.Votes,
		&choice.IsCorrect,
		&choice.QueueItemId,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(
				"%s: choice %s: %w",
				op, choiceId.String(), apperrors.ErrNotFound,
			)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &choice, nil
}

func (r *VoteRepository) CastVote(
	ctx context.Context, choiceId, userId uuid.UUID,
) error {
	const op = "repository.VoteRepository.CastVote"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO user_votes (vote_choice_id, user_id) VALUES ($1, $2)`,
		choiceId, userId,
	)
	if err != nil {
		return fmt.Errorf("%s: insert user vote: %w", op, err)
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE vote_choices SET votes = votes + 1 WHERE id = $1`,
		choiceId,
	)
	if err != nil {
		return fmt.Errorf("%s: increment votes: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}

func (r *VoteRepository) RemoveUserVote(
	ctx context.Context, voteId, userId uuid.UUID,
) error {
	const op = "repository.VoteRepository.RemoveUserVote"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var choiceId uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT uv.vote_choice_id FROM user_votes uv
		JOIN vote_choices vc ON uv.vote_choice_id = vc.id
		WHERE vc.vote_id = $1 AND uv.user_id = $2
	`, voteId, userId).Scan(&choiceId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"%s: %w: user has not voted", op, apperrors.ErrNotFound,
			)
		}
		return fmt.Errorf("%s: find user vote: %w", op, err)
	}

	_, err = tx.Exec(
		ctx,
		`DELETE FROM user_votes WHERE vote_choice_id = $1 AND user_id = $2`,
		choiceId, userId,
	)
	if err != nil {
		return fmt.Errorf("%s: delete user vote: %w", op, err)
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE vote_choices SET votes = votes - 1 WHERE id = $1`,
		choiceId,
	)
	if err != nil {
		return fmt.Errorf("%s: decrement votes: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}

func (r *VoteRepository) GetUserVoteForVote(
	ctx context.Context, voteId, userId uuid.UUID,
) (*uuid.UUID, error) {
	const op = "repository.VoteRepository.GetUserVoteForVote"

	query := `
		SELECT uv.vote_choice_id FROM user_votes uv
		JOIN vote_choices vc ON uv.vote_choice_id = vc.id
		WHERE vc.vote_id = $1 AND uv.user_id = $2
	`

	var choiceId uuid.UUID
	err := r.db.QueryRow(ctx, query, voteId, userId).Scan(&choiceId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &choiceId, nil
}
