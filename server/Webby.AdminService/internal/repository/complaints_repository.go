package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/admin-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type complaintsRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *complaintsRepository {
	return &complaintsRepository{db: db}
}

func (r *complaintsRepository) ListComplaints(ctx context.Context, offset, limit int) ([]models.Complaint, int, error) {
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
		Where("NOT EXISTS (SELECT 1 FROM complaint_results cr WHERE cr.complaint_id = \"Complaints\".\"ComplaintId\")").
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
		if err := rows.Scan(&complaint.ID, &complaint.AuthorID, &complaint.TargetType, &complaint.TargetID, &complaint.ReasonType, &complaint.AdditionalInfo, &complaint.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		complaints = append(complaints, complaint)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return complaints, 0, nil
}

func (r *complaintsRepository) ResolveComplaint(ctx context.Context, complaintID, userID uuid.UUID, isAccepted bool) error {
	const op = "repository.AcceptComplaint"

	sql, args, err := sq.Insert("complaint_results").
		Columns("complaint_id", "admin_id", "is_accepted").
		Values(complaintID, userID, isAccepted).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *complaintsRepository) IsResolved(ctx context.Context, complaintID uuid.UUID) (bool, error) {
	const op = "repository.IsResolved"

	sql, args, err := sq.Select("1").
		From("complaint_results").
		Where(sq.Eq{"complaint_id": complaintID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: sql build failed: %w", op, err)
	}

	var isResolved int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&isResolved)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("%s: failed to execute: %w", op, err)
	}

	return true, nil
}

func (r *complaintsRepository) GetActiveReportsCount(ctx context.Context) (int, error) {
	const op = "repository.GetActiveReportsCount"

	sql, args, err := sq.Select("COUNT(\"ComplaintId\")").
		From("\"Complaints\"").
		Where("NOT EXISTS (SELECT 1 FROM complaint_results cr WHERE cr.complaint_id = \"Complaints\".\"ComplaintId\")").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: build failed %w", op, err)
	}

	var activeComplaints int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&activeComplaints)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("%s: failed to execute: %w", op, err)
	}

	return activeComplaints, nil
}
