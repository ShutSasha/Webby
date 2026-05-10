package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomMemberRepository struct {
	db *pgxpool.Pool
}

func NewRoomMemberRepository(db *pgxpool.Pool) *RoomMemberRepository {
	return &RoomMemberRepository{db: db}
}

func (r *RoomMemberRepository) Create(ctx context.Context, member *models.RoomMember) error {
	const op = "repository.RoomMemberRepository.Create"

	if member == nil {
		return fmt.Errorf("%s: %w: member cannot be nil", op, apperrors.ErrInvalidInput)
	}

	if member.RoomId == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	if member.UserId == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid user id", op, apperrors.ErrInvalidInput)
	}

	query := `
		INSERT INTO room_members (room_id, user_id, room_points)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, member.RoomId, member.UserId, member.RoomPoints)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return nil
}

func (r *RoomMemberRepository) Exists(ctx context.Context, roomId, userId uuid.UUID) (bool, error) {
	const op = "repository.RoomMemberRepository.Exists"

	if roomId == uuid.Nil {
		return false, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	if userId == uuid.Nil {
		return false, fmt.Errorf("%s: %w: invalid user id", op, apperrors.ErrInvalidInput)
	}

	query := `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, roomId, userId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return exists, nil
}

func (r *RoomMemberRepository) EnsureMember(ctx context.Context, roomID, userID uuid.UUID) error {
	const op = "repository.RoomMemberRepository.EnsureMember"

	query := `
		INSERT INTO room_members (room_id, user_id, room_points)
		VALUES ($1, $2, 0)
		ON CONFLICT (room_id, user_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, roomID, userID)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return nil
}

func (r *RoomMemberRepository) ListByRoom(ctx context.Context, roomId uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error) {
	const op = "repository.RoomMemberRepository.ListByRoom"

	if roomId == uuid.Nil {
		return nil, 0, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
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

	whereClause := `WHERE rm.room_id = $1`
	countArgs := []any{roomId}
	paramN := 2

	if search != "" {
		whereClause += fmt.Sprintf(` AND u."Username" ILIKE $%d`, paramN)
		countArgs = append(countArgs, "%"+search+"%")
		paramN++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM room_members rm
		JOIN "Users" u ON rm.user_id = u."UserId"
		%s
	`, whereClause)

	var total int64
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count query failed: %w", op, err)
	}

	dataQuery := fmt.Sprintf(`
		SELECT u."UserId", u."Username", u."AvatarUrl", rm.room_points
		FROM room_members rm
		JOIN "Users" u ON rm.user_id = u."UserId"
		%s
		ORDER BY rm.room_points DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, paramN, paramN+1)

	dataArgs := append(countArgs, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: data query failed: %w", op, err)
	}
	defer rows.Close()

	members := []models.RoomMemberInfo{}
	for rows.Next() {
		var member models.RoomMemberInfo
		err := rows.Scan(
			&member.UserId,
			&member.Username,
			&member.AvatarUrl,
			&member.RoomPoints,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return members, total, nil
}

func (r *RoomMemberRepository) Delete(ctx context.Context, roomId, userId uuid.UUID) error {
	const op = "repository.RoomMemberRepository.Delete"

	if roomId == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	if userId == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid user id", op, apperrors.ErrInvalidInput)
	}

	query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`

	tag, err := r.db.Exec(ctx, query, roomId, userId)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: member not found in room: %w", op, apperrors.ErrNotFound)
	}

	return nil
}

func (r *RoomMemberRepository) UpdatePoints(ctx context.Context, roomId, userId uuid.UUID, delta int) (*models.RoomMemberInfo, error) {
	const op = "repository.RoomMemberRepository.UpdatePoints"

	if roomId == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	if userId == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid user id", op, apperrors.ErrInvalidInput)
	}

	query := `
		UPDATE room_members
		SET room_points = room_points + $3
		WHERE room_id = $1 AND user_id = $2 AND room_points + $3 >= 0
		RETURNING (SELECT u."UserId" FROM "Users" u WHERE u."UserId" = room_members.user_id),
			(SELECT u."Username" FROM "Users" u WHERE u."UserId" = room_members.user_id),
			(SELECT u."AvatarUrl" FROM "Users" u WHERE u."UserId" = room_members.user_id),
			room_points
	`

	var member models.RoomMemberInfo
	err := r.db.QueryRow(ctx, query, roomId, userId, delta).Scan(
		&member.UserId,
		&member.Username,
		&member.AvatarUrl,
		&member.RoomPoints,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			exists, exErr := r.Exists(ctx, roomId, userId)
			if exErr != nil {
				return nil, fmt.Errorf("%s: %w", op, exErr)
			}
			if exists {
				return nil, fmt.Errorf("%s: %w: insufficient points", op, apperrors.ErrInvalidInput)
			}
			return nil, fmt.Errorf("%s: member not found in room: %w", op, apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return &member, nil
}
