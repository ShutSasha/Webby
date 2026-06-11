package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/room-category-service/internal/apperrors"
	"webby/room-category-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, name string) error {
	const op = "repository.Create"

	query := sq.Insert("categories").
		Columns("name").
		Values(name).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w: category with name '%s' already exists", op, apperrors.ErrCategoryAlreadyExists, name)
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, search string, offset int, limit int) ([]models.Category, int, error) {
	const op = "repository.List"

	query := sq.Select("id", "name", "COUNT(*) OVER() AS total").
		From("categories").
		OrderBy("name ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	if search != "" {
		query = query.Where(sq.ILike{"name": "%" + search + "%"})
	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query building failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var total int
	categories := make([]models.Category, 0, limit)

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &total); err != nil {
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

	query := sq.Update("categories").
		Set("name", newName).
		Where(sq.Eq{"name": oldName}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("%s: query building failed: %w", op, err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w: category with name '%s' already exists", op, apperrors.ErrCategoryAlreadyExists, newName)
			}
		}
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: category '%s': %w", op, oldName, apperrors.ErrCategoryNotFound)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, name string) error {
	const op = "repository.CategoryRepository.Delete"

	query := sq.Delete("categories").
		Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("%s: query building failed: %w", op, err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: category '%s': %w", op, name, apperrors.ErrCategoryNotFound)
	}

	return nil
}

func (r *Repository) Exists(ctx context.Context, name string) (bool, error) {
	const op = "repository.CategoryRepository.Exists"

	subQuery := sq.Select("1").
		From("categories").
		Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar)

	subSQL, subArgs, err := subQuery.ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: subquery building failed: %w", op, err)
	}
	query := fmt.Sprintf("SELECT EXISTS(%s)", subSQL)

	var exists bool
	err = r.db.QueryRow(ctx, query, subArgs...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}
