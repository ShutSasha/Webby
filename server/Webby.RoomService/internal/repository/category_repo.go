package repository

import (
	"database/sql"
	"webby/internal/models"

	"github.com/google/uuid"
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
	if r.db == nil {
		return uuid.Nil, ErrDatabaseConnection("database connection is nil")
	}

	if name == "" {
		return uuid.Nil, ErrCategoryCreationFailed("category name cannot be empty")
	}

	id := uuid.New()

	query := `
		INSERT INTO categories (id, name)
		VALUES ($1, $2)
	`

	err := r.db.QueryRow(query, id, name).Err()
	if err != nil {
		return uuid.Nil, ErrCategoryCreationFailed(err.Error())
	}

	return id, nil
}

func (r *CategoryRepository) List(search string, page int, limit int) ([]models.Category, int64, error) {
	if r.db == nil {
		return nil, 0, ErrDatabaseConnection("database connection is nil")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

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
		return nil, 0, ErrCategoriesFetchFailed(err.Error())
	}

	rows, err := r.db.Query(dataQuery, queryArgs...)
	if err != nil {
		return nil, 0, ErrCategoriesFetchFailed(err.Error())
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
			return nil, 0, ErrCategoriesFetchFailed(err.Error())
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, ErrCategoriesFetchFailed(err.Error())
	}

	return categories, total, nil
}

func (r *CategoryRepository) Update(id uuid.UUID, name string) error {
	if r.db == nil {
		return ErrDatabaseConnection("database connection is nil")
	}

	if id == uuid.Nil {
		return ErrCategoryUpdateFailed("invalid category id")
	}

	if name == "" {
		return ErrCategoryUpdateFailed("category name cannot be empty")
	}

	query := `
		UPDATE categories
		SET name = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, name, id)
	if err != nil {
		return ErrCategoryUpdateFailed(err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrCategoryUpdateFailed(err.Error())
	}

	if rowsAffected == 0 {
		return ErrCategoryNotFound(id.String())
	}

	return nil
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return ErrDatabaseConnection("database connection is nil")
	}

	if id == uuid.Nil {
		return ErrCategoryDeletionFailed("invalid category id")
	}

	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return ErrCategoryDeletionFailed(err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrCategoryDeletionFailed(err.Error())
	}

	if rowsAffected == 0 {
		return ErrCategoryNotFound(id.String())
	}

	return nil
}
