package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"webby/internal/apperrors"
	"webby/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) Create(name string) (uuid.UUID, error) {
	const op = "repository.CategoryRepository.Create"

	id := uuid.New()

	query := `
		INSERT INTO categories (id, name)
		VALUES ($1, $2)
	`

	err := r.db.QueryRow(query, id, name).Err()
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return uuid.Nil, fmt.Errorf("%s: %w: category with name '%s' already exists", op, apperrors.ErrConflict, name)
			}
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *CategoryRepository) List(search string, offset int, limit int) ([]models.Category, int64, error) {
	const op = "repository.CategoryRepository.List"

	countQuery := `SELECT COUNT(*) FROM categories`
	dataQuery := `
		SELECT id, name
		FROM categories
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	queryArgs := []any{limit, offset}

	if search != "" {
		countQuery += ` WHERE name ILIKE $1`
		dataQuery = `
			SELECT id, name
			FROM categories
			WHERE name ILIKE $1
			ORDER BY name ASC
			LIMIT $2 OFFSET $3
		`
		queryArgs = []any{"%" + search + "%", limit, offset}
	}

	var total int64
	countQueryArgs := []any{}
	if search != "" {
		countQueryArgs = []any{"%" + search + "%"}
	}

	err := r.db.QueryRow(countQuery, countQueryArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count query failed: %w", op, err)
	}

	rows, err := r.db.Query(dataQuery, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: data query failed: %w", op, err)
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.Id,
			&category.Name,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return categories, total, nil
}

func (r *CategoryRepository) Update(id uuid.UUID, name string) error {
	const op = "repository.CategoryRepository.Update"

	if id == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid category id", op, apperrors.ErrInvalidInput)
	}

	if name == "" {
		return fmt.Errorf("%s: %w: category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	query := `
		UPDATE categories
		SET name = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, name, id)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: getting rows affected failed: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: category %s: %w", op, id.String(), apperrors.ErrNotFound)
	}

	return nil
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	const op = "repository.CategoryRepository.Delete"

	if id == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid category id", op, apperrors.ErrInvalidInput)
	}

	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: getting rows affected failed: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: category %s: %w", op, id.String(), apperrors.ErrNotFound)
	}

	return nil
}
