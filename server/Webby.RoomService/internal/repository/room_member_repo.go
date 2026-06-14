package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type roomMemberRepository struct {
	db *pgxpool.Pool
}

func NewRoomMemberRepository(db *pgxpool.Pool) *roomMemberRepository {
	return &roomMemberRepository{db: db}
}

func (r *roomMemberRepository) Exists(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
	const op = "repository.roomMemberRepository.Exists"

	query, args, err := sq.Select("1").
		From("room_members").
		Where(sq.Eq{
			"room_id": roomID,
			"user_id": userID,
		}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var dummy int
	err = r.db.QueryRow(ctx, query, args...).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return true, nil
}

func (r *roomMemberRepository) EnsureMember(ctx context.Context, roomID, userID uuid.UUID) error {
	const op = "repository.roomMemberRepository.EnsureMember"

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("room_members").
		Columns("room_id", "user_id", "room_points").
		Values(roomID, userID, 0).
		Suffix("ON CONFLICT (room_id, user_id) DO NOTHING").
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build failed: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return nil
}

func (r *roomMemberRepository) ListByRoom(ctx context.Context, roomId uuid.UUID, page, limit int, search string) ([]models.RoomMemberInfo, int64, error) {
	const op = "repository.roomMemberRepository.ListByRoom"

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
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	countBuilder := builder.
		Select("COUNT(*)").
		From("room_members rm").
		Join("\"Users\" u ON rm.user_id = u.\"UserId\"").
		Where(sq.Eq{"rm.room_id": roomId})

	dataBuilder := builder.
		Select("u.\"UserId\"", "u.\"Username\"", "u.\"AvatarUrl\"", "rm.room_points").
		From("room_members rm").
		Join("\"Users\" u ON rm.user_id = u.\"UserId\"").
		Where(sq.Eq{"rm.room_id": roomId}).
		OrderBy("rm.room_points DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	if search != "" {
		pattern := "%" + search + "%"
		countBuilder = countBuilder.Where(sq.Expr("u.\"Username\" ILIKE ?", pattern))
		dataBuilder = dataBuilder.Where(sq.Expr("u.\"Username\" ILIKE ?", pattern))
	}

	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var total int64
	err = r.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count query failed: %w", op, err)
	}

	dataSQL, dataArgs, err := dataBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: build failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: data query failed: %w", op, err)
	}
	defer rows.Close()

	members := []models.RoomMemberInfo{}
	for rows.Next() {
		var member models.RoomMemberInfo
		err := rows.Scan(
			&member.UserID,
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

func (r *roomMemberRepository) Delete(ctx context.Context, roomID, userID uuid.UUID) error {
	const op = "repository.roomMemberRepository.Delete"

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("room_members").
		Where(sq.Eq{"room_id": roomID, "user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build failed: %w", op, err)
	}

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, apperrors.ErrRoomMemberNotFound)
	}

	return nil
}

func (r *roomMemberRepository) AddPointsBulk(ctx context.Context, userIDs []uuid.UUID, pointsToAdd int) (map[uuid.UUID]int, error) {
	const op = "repository.roomMemberRepository.AddPointsBulk"

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("room_members").
		Set("room_points", sq.Expr("room_points + ?", pointsToAdd)).
		Where(sq.Eq{"user_id": userIDs}).
		Suffix("RETURNING user_id, room_points").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}
	defer rows.Close()

	updatedTotals := make(map[uuid.UUID]int, len(userIDs))

	for rows.Next() {
		var uID uuid.UUID
		var totalPoints int

		if err := rows.Scan(&uID, &totalPoints); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		updatedTotals[uID] = totalPoints
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return updatedTotals, nil
}

func (r *roomMemberRepository) GetMemberPoints(ctx context.Context, roomID, userID uuid.UUID) (int, error) {
	const op = "repository.roomMemberRepository.GetMemberPoints"

	sql, args, err := sq.Select("room_points").
		From("room_members").
		Where(sq.And{
			sq.Eq{"room_id": roomID},
			sq.Eq{"user_id": userID},
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var points int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&points)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, apperrors.ErrRoomMemberNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return points, nil
}
