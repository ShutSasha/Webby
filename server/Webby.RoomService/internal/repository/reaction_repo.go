package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/room-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type reactionRepository struct {
	db *pgxpool.Pool
}

func NewReactionRepository(db *pgxpool.Pool) *reactionRepository {
	return &reactionRepository{db: db}
}

func (r *reactionRepository) List(ctx context.Context) ([]models.Reaction, error) {
	const op = "repository.reactionRepository.List"

	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "name", "cost", "sticker_url").
		From("reactions").
		OrderBy("name").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Reaction{}, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	reactions := []models.Reaction{}
	for rows.Next() {
		var reaction models.Reaction
		if err := rows.Scan(&reaction.ID, &reaction.Name, &reaction.Cost, &reaction.StickerURL); err != nil {
			return nil, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		reactions = append(reactions, reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return reactions, nil
}

func (r *reactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Reaction, error) {
	const op = "reactionRepository.GetByID"

	sql, args, err := sq.Select("id", "name", "cost", "sticker_url").
		From("reactions").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var reaction models.Reaction
	err = r.db.QueryRow(ctx, sql, args...).Scan(&reaction.ID, &reaction.Name, &reaction.Cost, &reaction.StickerURL)
	if err != nil {
		return nil, fmt.Errorf("%s: row scan failed: %w", op, err)
	}

	return &reaction, nil
}
