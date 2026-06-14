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
	const op = "repository.chatRepository.Create"

	query := sq.Insert("chats").
		Columns("room_id").
		Values(chat.RoomID).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&chat.ID, &chat.CreatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return chat.ID, nil
}

func (r *chatRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Chat, error) {
	const op = "repository.chatRepository.GetByID"

	query := sq.Select("id", "room_id", "created_at").
		From("chats").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var chat models.Chat
	err = r.db.QueryRow(ctx, sql, args...).Scan(&chat.ID, &chat.RoomID, &chat.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: chat %s: %w", op, id.String(), apperrors.ErrChatNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &chat, nil
}

func (r *chatRepository) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*models.Chat, error) {
	const op = "repository.chatRepository.GetByRoomID"

	query := sq.Select("id", "room_id", "created_at").
		From("chats").
		Where(sq.Eq{"room_id": roomID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var chat models.Chat
	err = r.db.QueryRow(ctx, sql, args...).Scan(&chat.ID, &chat.RoomID, &chat.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, apperrors.ErrChatNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &chat, nil
}

func (r *chatRepository) GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	const op = "repository.chatRepository.GetChatIDByRoomID"

	query := sq.Select("id").
		From("chats").
		Where(sq.Eq{"room_id": roomID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var id uuid.UUID
	err = r.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, apperrors.ErrChatNotFound)
		}
		return uuid.Nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return id, nil
}

func (r *chatRepository) History(ctx context.Context, userID uuid.UUID, userIDs []uuid.UUID) ([]models.ChatHistoryItem, int, error) {
	const op = "repository.chatRepository.History"

	subQuery := sq.Select("chat_id").
		From("chat_members").
		Where(sq.Eq{"user_id": userID})

	subQuerySql, subQueryArgs, err := subQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: building subquery: %w", op, err)
	}

	queryBuilder := sq.Select(
		"cm.chat_id",
		"cm.user_id",
		"m.content AS last_message",
		"m.created_at AS last_message_sent_at",
	).
		From("chat_members cm").
		Join("chats c ON cm.chat_id = c.id").
		LeftJoin("messages m ON c.last_message_id = m.id").
		Where(sq.Expr("cm.chat_id IN ("+subQuerySql+")", subQueryArgs...)).
		Where(sq.NotEq{"cm.user_id": userID}).
		Where(sq.Eq{"c.room_id": nil}).
		Where(sq.Eq{"cm.user_id": userIDs}).
		OrderBy("last_message_sent_at DESC NULLS LAST").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: building query: %w", op, err)
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: executing query: %w", op, err)
	}
	defer rows.Close()

	items := make([]models.ChatHistoryItem, 0)
	for rows.Next() {
		var item models.ChatHistoryItem
		var content *string
		var createdAt *time.Time

		err := rows.Scan(
			&item.ChatID,
			&item.User.ID,
			&content,
			&createdAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: scanning row: %w", op, err)
		}

		if content != nil {
			item.LastMessage.Content = *content
		}
		if createdAt != nil {
			item.LastMessage.CreatedAt = *createdAt
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: iterating rows: %w", op, err)
	}

	return items, len(items), nil
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
			return fmt.Errorf("%s: chat not found: %w", op, apperrors.ErrChatNotFound)
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

func (r *chatRepository) Delete(ctx context.Context, chatID uuid.UUID) error {
	const op = "repository.chatRepository.Delete"

	sql, args, err := sq.Delete("chats").
		Where(sq.Eq{"id": chatID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build sql query: %w", op, err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, apperrors.ErrChatNotFound)
	}

	return nil
}

func (r *chatRepository) GetByMembers(ctx context.Context, firstUserID, secondUserID uuid.UUID) (*models.Chat, error) {
	const op = "repository.chatRepository.GetByMembers"

	sql, args, err := sq.Select("c.id", "c.created_at").
		From("chats c").
		Join("chat_members cm1 ON c.id = cm1.chat_id").
		Join("chat_members cm2 ON c.id = cm2.chat_id").
		Where(sq.Eq{"cm1.user_id": firstUserID, "cm2.user_id": secondUserID}).
		Where(sq.Eq{"c.room_id": nil}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build sql query: %w", op, err)
	}

	var chat models.Chat
	err = r.db.QueryRow(ctx, sql, args...).Scan(&chat.ID, &chat.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: exec query: %w", op, err)
	}

	return &chat, nil
}
