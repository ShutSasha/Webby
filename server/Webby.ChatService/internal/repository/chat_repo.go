package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type chatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *chatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) Create(ctx context.Context, chat *models.Chat) (uuid.UUID, error) {
	const op = "repository.ChatRepository.Create"

	chat.ID = uuid.New()

	query := `
		INSERT INTO chats (id, room_id)
		VALUES ($1, $2)
		RETURNING created_at
	`

	err := r.db.QueryRow(ctx, query, chat.ID, chat.RoomID).Scan(&chat.CreatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return chat.ID, nil
}

func (r *chatRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Chat, error) {
	const op = "repository.ChatRepository.GetByID"

	if id == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid chat id", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT id, room_id, created_at
		FROM chats
		WHERE id = $1
	`

	var chat models.Chat
	err := r.db.QueryRow(ctx, query, id).Scan(&chat.ID, &chat.RoomID, &chat.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: chat %s: %w", op, id.String(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &chat, nil
}

func (r *chatRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*models.Chat, error) {
	const op = "repository.ChatRepository.GetByRoomID"

	if roomID == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT id, room_id, created_at
		FROM chats
		WHERE room_id = $1
	`

	var chat models.Chat
	err := r.db.QueryRow(ctx, query, roomID).Scan(&chat.ID, &chat.RoomID, &chat.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: room %s: %w", op, roomID.String(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &chat, nil
}

func (r *chatRepository) GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	const op = "repository.ChatRepository.GetChatIDByRoomID"

	query := `SELECT id FROM chats WHERE room_id = $1`
	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, roomID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: room %s: %w", op, roomID.String(), apperrors.ErrNotFound)
		}
		return uuid.Nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return id, nil
}

func (r *chatRepository) History(ctx context.Context, userID uuid.UUID, userIDs []uuid.UUID, offset, limit int) ([]models.ChatHistoryItem, int, error) {
	const op = "repository.chatRepository.History"

	query := `
	SELECT cm.chat_id, cm.user_id, m.content, m.created_at
	FROM chat_members cm
	JOIN chats c ON cm.chat_id = c.id
	JOIN messages m ON c.last_message_id = m.id
	WHERE cm.chat_id IN (
		SELECT cm2chat_id
		FROM chat_members cm2
		JOIN (
			SELECT cm3.chat_id
			FROM chat_members cm3
			WHERE cm3.user_id = $1
		) ON cm2.chat_id = cm3.chat_id
		WHERE cm2.id IN ($2)
	) AND c.room_id IS NULL;
	`
	_ = query
	return nil, 0, nil
}

func (r *chatRepository) UpdateLastMessage(ctx context.Context, chatID, lastMessageID uuid.UUID, createdAt time.Time) error {
	const op = "repositories.chatRepository.UpdateLastMessage"

	currentMessageTimeSubquery := sq.Select("COALESCE(m.created_at, '0001-01-01 00:00:00'::timestamp)").
		From("chats c").
		LeftJoin("messages m ON c.last_message_id = m.id").
		Where(sq.Eq{"c.id": chatID})

	subQuerySql, subQueryArgs, err := currentMessageTimeSubquery.ToSql()
	if err != nil {
		return fmt.Errorf("%s: building subquery: %w", op, err)
	}

	exprArgs := append([]interface{}{createdAt}, subQueryArgs...)

	query, args, err := sq.Update("chats").
		Set("last_message_id", lastMessageID).
		Where(sq.Eq{"id": chatID}).
		Where(sq.Expr("? >= ("+subQuerySql+")", exprArgs...)).	
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%s: building query: %w", op, err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		var exists bool
		checkErr := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM chats WHERE id = $1)", chatID).Scan(&exists)
		if checkErr == nil && !exists {
			return fmt.Errorf("%s: chat not found: %w", op, apperrors.ErrNotFound)
		}

		return nil
	}
	return nil
}

func (r *chatRepository) Exists(ctx context.Context, firstUserID, secondUserID uuid.UUID) (bool, error) {
	const op = "repositories.chatRepository.Exists"

	sql, args, err := sq.Select("1").
		From("chat_members cm1").
		Join("chat_members cm2 ON cm1.chat_id = cm2.chat_id").
		Join("chats c on cm1.chat_id = c.id").
		Where(sq.And{
			sq.Eq{"cm1.user_id": firstUserID.String()},
			sq.Eq{"cm2.user_id": secondUserID.String()},
			sq.Eq{"c.room_id": nil},
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: building query: %w", op, err)
	}

	var exists bool
	err = r.db.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return true, nil
}
