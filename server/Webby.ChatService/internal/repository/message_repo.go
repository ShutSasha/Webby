package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *messageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, msg *models.Message) (*models.Message, error) {
	const op = "repository.MessageRepository.Create"

	query := sq.Insert("messages").
		Columns("sender_id", "chat_id", "content").
		Values(msg.SenderID, msg.ChatID, msg.Content).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return msg, nil
}

func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	const op = "repository.messageRepository.GetByID"

	query := sq.Select(
		"id", "sender_id", "chat_id", "content",
		"is_edited", "edited_at", "created_at",
	).
		From("messages").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var msg models.Message
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&msg.ID, &msg.SenderID, &msg.ChatID,
		&msg.Content, &msg.IsEdited, &msg.EditedAt, &msg.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: message %s: %w", op, id.String(), apperrors.ErrMessageNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &msg, nil
}

func (r *messageRepository) Update(ctx context.Context, id uuid.UUID, content string) (*models.Message, error) {
	const op = "repository.messageRepository.Update"

	query := sq.Update("messages").
		Set("content", content).
		Set("is_edited", true).
		Set("edited_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, sender_id, chat_id, content, is_edited, edited_at, created_at").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var msg models.Message
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&msg.ID, &msg.SenderID, &msg.ChatID,
		&msg.Content, &msg.IsEdited, &msg.EditedAt, &msg.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: message %s: %w", op, id.String(), apperrors.ErrMessageNotFound)
		}
		return nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return &msg, nil
}

func (r *messageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.messageRepository.Delete"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx failed: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var chatID uuid.UUID
	getChatQuery := sq.Select("chat_id").
		From("messages").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := getChatQuery.ToSql()
	if err != nil {
		return fmt.Errorf("%s: build select query failed: %w", op, err)
	}

	err = tx.QueryRow(ctx, sqlStr, args...).Scan(&chatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, apperrors.ErrMessageNotFound)
		}
		return fmt.Errorf("%s: get chat_id failed: %w", op, err)
	}

	updateSubQuery, updateSubArgs, err := sq.Select("id").
		From("messages").
		Where(sq.And{
			sq.Eq{"chat_id": chatID},
			sq.NotEq{"id": id},
		}).
		OrderBy("created_at DESC").
		Limit(1).ToSql()
	if err != nil {
		return fmt.Errorf("%s: build update subquery failed: %w", op, err)
	}

	updateChatSQL, updateArgs, err := sq.Update("chats").
		Set("last_message_id", sq.Expr("("+updateSubQuery+")", updateSubArgs...)).
		Where(sq.Eq{
			"id":              chatID,
			"last_message_id": id,
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%s: build update query failed: %w", op, err)
	}

	_, err = tx.Exec(ctx, updateChatSQL, updateArgs...)
	if err != nil {
		return fmt.Errorf("%s: update last_message_id failed: %w", op, err)
	}

	deleteQuery := sq.Delete("messages").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	delSqlStr, delArgs, err := deleteQuery.ToSql()
	if err != nil {
		return fmt.Errorf("%s: build delete failed: %w", op, err)
	}

	tag, err := tx.Exec(ctx, delSqlStr, delArgs...)
	if err != nil {
		return fmt.Errorf("%s: delete execution failed: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, apperrors.ErrMessageNotFound)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit tx failed: %w", op, err)
	}

	return nil
}

func (r *messageRepository) List(ctx context.Context, chatID uuid.UUID, offset, limit int) ([]models.Message, int, error) {
	const op = "repository.MessageRepository.ListByChat"

	countQuery := sq.Select("COUNT(*)").
		From("messages").
		Where(sq.Eq{"chat_id": chatID}).
		PlaceholderFormat(sq.Dollar)
	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count build failed: %w", op, err)
	}

	var total int
	if err := r.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%s: count failed: %w", op, err)
	}

	query := sq.Select(
		"id", "sender_id", "chat_id", "content",
		"is_edited", "edited_at", "created_at",
	).
		From("messages").
		Where(sq.Eq{"chat_id": chatID}).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: build failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.ChatID,
			&msg.Content, &msg.IsEdited, &msg.EditedAt, &msg.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%s: scan failed: %w", op, err)
		}
		messages = append(messages, msg)
	}

	return messages, total, nil
}
