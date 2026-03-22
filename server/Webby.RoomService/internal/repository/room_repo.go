package repository

import (
	"database/sql"
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
	if r.db == nil {
		return uuid.Nil, ErrDatabaseConnection("database connection is nil")
	}

	room.Id = uuid.New()

	query := `
		INSERT INTO rooms (id, host_id, category_id, name, thumbnail, is_private, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	err := r.db.QueryRow(
		query,
		room.Id,
		room.HostId,
		room.CategoryId,
		room.Name,
		room.Thumbnail,
		room.IsPrivate,
		room.CreatedAt,
	).Err()

	if err != nil {
		return uuid.Nil, ErrRoomCreationFailed(err.Error())
	}

	return room.Id, nil
}

func (r *RoomRepository) Delete(id uuid.UUID) error {
	if r.db == nil {
		return ErrDatabaseConnection("database connection is nil")
	}

	if id == uuid.Nil {
		return ErrRoomDeletionFailed("invalid room id")
	}

	query := `DELETE FROM rooms WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return ErrRoomDeletionFailed(err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrRoomDeletionFailed(err.Error())
	}

	if rowsAffected == 0 {
		return ErrRoomNotFound(id.String())
	}

	return nil
}

func (r *RoomRepository) GetById(id uuid.UUID) (*models.Room, error) {
	if r.db == nil {
		return nil, ErrDatabaseConnection("database connection is nil")
	}

	if id == uuid.Nil {
		return nil, ErrRoomNotFound("invalid room id")
	}

	query := `
		SELECT id, host_id, category_id, name, thumbnail, is_private, created_at
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
		&room.IsPrivate,
		&room.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRoomNotFound(id.String())
		}
		return nil, ErrRoomsFetchFailed(err.Error())
	}

	return &room, nil
}

func (r *RoomRepository) ListMy(userId uuid.UUID, page int, limit int) ([]models.Room, int64, error) {
	if r.db == nil {
		return nil, 0, ErrDatabaseConnection("database connection is nil")
	}

	if userId == uuid.Nil {
		return nil, 0, ErrRoomsFetchFailed("invalid user id")
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
		return nil, 0, ErrRoomsFetchFailed(err.Error())
	}

	query := `
		SELECT id, host_id, category_id, name, thumbnail, is_private, created_at
		FROM rooms
		WHERE host_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, 0, ErrRoomsFetchFailed(err.Error())
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
			&room.IsPrivate,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, 0, ErrRoomsFetchFailed(err.Error())
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, ErrRoomsFetchFailed(err.Error())
	}

	return rooms, total, nil
}

func (r *RoomRepository) ListPublic(page int, limit int) ([]models.Room, int64, error) {
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

	countQuery := `SELECT COUNT(*) FROM rooms WHERE is_private = false`
	var total int64
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, ErrRoomsFetchFailed(err.Error())
	}

	query := `
		SELECT id, host_id, category_id, name, thumbnail, is_private, created_at
		FROM rooms
		WHERE is_private = false
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, ErrRoomsFetchFailed(err.Error())
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
			&room.IsPrivate,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, 0, ErrRoomsFetchFailed(err.Error())
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, ErrRoomsFetchFailed(err.Error())
	}

	return rooms, total, nil
}

func (r *RoomRepository) Update(room *models.Room) (uuid.UUID, error) {
	if r.db == nil {
		return uuid.Nil, ErrDatabaseConnection("database connection is nil")
	}

	if room.Id == uuid.Nil {
		return uuid.Nil, ErrRoomUpdateFailed("invalid room id")
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
		return uuid.Nil, ErrRoomUpdateFailed(err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return uuid.Nil, ErrRoomUpdateFailed(err.Error())
	}

	if rowsAffected == 0 {
		return uuid.Nil, ErrRoomNotFound(room.Id.String())
	}

	return room.Id, nil
}
