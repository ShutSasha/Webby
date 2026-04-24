package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/internal/apperrors"
	"webby/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, name string) error {
	const op = "repository.CategoryRepository.Create"

	query := `
		INSERT INTO categories (name)
		VALUES ($1)
	`

	_, err := r.db.Exec(ctx, query, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w: category with name '%s' already exists", op, apperrors.ErrConflict, name)
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, search string, offset int, limit int) ([]models.Category, int64, error) {
	const op = "repository.CategoryRepository.List"

	query := `
        WITH filtered_cats AS (
            SELECT id, name 
            FROM categories 
            WHERE ($1 = '' OR name ILIKE $1)
        ),
        total_count AS (
            SELECT count(*) AS total FROM filtered_cats
        )
        SELECT f.id, f.name, t.total
        FROM filtered_cats f, total_count t
        ORDER BY f.name ASC
        LIMIT $2 OFFSET $3
    `

	searchParam := ""
	if search != "" {
		searchParam = "%" + search + "%"
	}

	rows, err := r.db.Query(ctx, query, searchParam, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var total int64
	categories := make([]models.Category, 0, limit)

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.Id, &category.Name, &total); err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return categories, total, nil
}

func (r *Repository) Update(ctx context.Context, oldName string, newName string) error {
	const op = "repository.CategoryRepository.Update"

	query := `
		UPDATE categories
		SET name = $1
		WHERE name = $2
	`

	tag, err := r.db.Exec(ctx, query, newName, oldName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w: category with name '%s' already exists", op, apperrors.ErrConflict, newName)
			}
		}
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: category '%s': %w", op, oldName, apperrors.ErrNotFound)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, name string) error {
	const op = "repository.CategoryRepository.Delete"

	query := `DELETE FROM categories WHERE name = $1`

	tag, err := r.db.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: category '%s': %w", op, name, apperrors.ErrNotFound)
	}

	return nil
}

func (r *Repository) Exists(ctx context.Context, name string) (bool, error) {
	const op = "repository.CategoryRepository.Exists"

	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}
