package repository

import (
	"context"
	"fmt"
	"webby/admin-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListComplaints(ctx context.Context, offset, limit int) ([]models.Complaint, int, error) {
	const op = "repository.ListComplaints"

	query := sq.Select(
		"\"ComplaintId\"",
		"\"AuthorId\"",
		"\"TargetType\"",
		"\"TargetId\"",
		"\"ReasonType\"",
		"\"AdditionalInfo\"",
		"\"CreatedAt\"",
	).From("\"Complaints\"").
		OrderBy("\"CreatedAt\" DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	complaints := make([]models.Complaint, 0, limit)
	for rows.Next() {
		var complaint models.Complaint
		if err := rows.Scan(&complaint); err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		complaints = append(complaints, complaint)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return complaints, 0, nil
}
