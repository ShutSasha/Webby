package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"webby/internal/apperrors"
	"webby/internal/models"

	"github.com/google/uuid"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *models.Room) (uuid.UUID, error) {
	const op = "repository.RoomRepository.Create"

	if room == nil {
		return uuid.Nil, fmt.Errorf("%s: %w: room cannot be nil", op, apperrors.ErrInvalidInput)
	}

	room.Id = uuid.New()

	query := `
		INSERT INTO rooms (id, host_id, category_id, name, thumbnail, token, is_private)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	err := r.db.QueryRow(
		query,
		room.Id,
		room.HostId,
		room.CategoryId,
		room.Name,
		room.Thumbnail,
		room.Token,
		room.IsPrivate,
	).Scan(&room.CreatedAt)

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return room.Id, nil
}

func (r *RoomRepository) Delete(id uuid.UUID) error {
	const op = "repository.RoomRepository.Delete"

	if id == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `DELETE FROM rooms WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: getting rows affected failed: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: room %s: %w", op, id.String(), apperrors.ErrNotFound)
	}

	return nil
}

func (r *RoomRepository) GetById(id uuid.UUID) (*models.Room, error) {
	const op = "repository.RoomRepository.GetById"

	if id == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT id, host_id, category_id, name, thumbnail, token, is_private, created_at
		FROM rooms
		WHERE id = $1
	`

	var room models.Room
	err := r.db.QueryRow(query, id).Scan(
		&room.Id,
		&room.HostId,
		&room.CategoryId,
		&room.Name,
		&room.Thumbnail,
		&room.Token,
		&room.IsPrivate,
		&room.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: room %s: %w", op, id.String(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &room, nil
}

func (r *RoomRepository) GetByToken(token string) (*models.Room, error) {
	const op = "repository.RoomRepository.GetByToken"

	if token == "" {
		return nil, fmt.Errorf("%s: %w: token cannot be empty", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT id, host_id, category_id, name, thumbnail, token, is_private, created_at
		FROM rooms
		WHERE token = $1
	`

	var room models.Room
	err := r.db.QueryRow(query, token).Scan(
		&room.Id,
		&room.HostId,
		&room.CategoryId,
		&room.Name,
		&room.Thumbnail,
		&room.Token,
		&room.IsPrivate,
		&room.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: room with token %s: %w", op, token, apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &room, nil
}

func (r *RoomRepository) ListMy(userId uuid.UUID, page int, limit int) ([]models.Room, int64, error) {
	const op = "repository.RoomRepository.ListMy"

	if userId == uuid.Nil {
		return nil, 0, fmt.Errorf("%s: %w: invalid user id", op, apperrors.ErrInvalidInput)
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

	countQuery := `SELECT COUNT(*) FROM rooms WHERE host_id = $1`
	var total int64
	err := r.db.QueryRow(countQuery, userId).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count query failed: %w", op, err)
	}

	query := `
		SELECT id, host_id, category_id, name, thumbnail, token, is_private, created_at
		FROM rooms
		WHERE host_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: data query failed: %w", op, err)
	}
	defer rows.Close()

	rooms := []models.Room{}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.Id,
			&room.HostId,
			&room.CategoryId,
			&room.Name,
			&room.Thumbnail,
			&room.Token,
			&room.IsPrivate,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return rooms, total, nil
}

func (r *RoomRepository) ListPublic(page int, limit int, search string, categoryId *uuid.UUID) ([]models.Room, int64, error) {
	const op = "repository.RoomRepository.ListPublic"

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

	whereClause := "WHERE is_private = false"
	countArgs := []any{}
	paramN := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND name ILIKE $%d", paramN)
		countArgs = append(countArgs, "%"+search+"%")
		paramN++
	}

	if categoryId != nil {
		whereClause += fmt.Sprintf(" AND category_id = $%d", paramN)
		countArgs = append(countArgs, *categoryId)
		paramN++
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM rooms %s`, whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count query failed: %w", op, err)
	}

	query := fmt.Sprintf(`
		SELECT id, host_id, category_id, name, thumbnail, token, is_private, created_at
		FROM rooms
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, paramN, paramN+1)

	dataArgs := append(countArgs, limit, offset)

	rows, err := r.db.Query(query, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: data query failed: %w", op, err)
	}
	defer rows.Close()

	rooms := []models.Room{}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.Id,
			&room.HostId,
			&room.CategoryId,
			&room.Name,
			&room.Thumbnail,
			&room.Token,
			&room.IsPrivate,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return rooms, total, nil
}

func (r *RoomRepository) Update(room *models.Room) (uuid.UUID, error) {
	const op = "repository.RoomRepository.Update"

	if room == nil || room.Id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `
		UPDATE rooms
		SET category_id = $1, name = $2, thumbnail = $3, is_private = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(
		query,
		room.CategoryId,
		room.Name,
		room.Thumbnail,
		room.IsPrivate,
		room.Id,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: getting rows affected failed: %w", op, err)
	}

	if rowsAffected == 0 {
		return uuid.Nil, fmt.Errorf("%s: room %s: %w", op, room.Id.String(), apperrors.ErrNotFound)
	}

	return room.Id, nil
}
