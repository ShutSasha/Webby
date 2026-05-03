package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveToken(ctx context.Context, token string, userID uuid.UUID) error {
	const op = "repository.SaveToken"

	query := `
		INSERT INTO user_wstokens (user_id, token)
		VALUES($1, $2)
		ON CONFLICT (user_id) DO UPDATE 
        SET token = EXCLUDED.token, created_at = NOW();
	`

	if _, err := r.db.Exec(ctx, query, userID, token); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
