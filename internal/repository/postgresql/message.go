package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"spsu-chat/internal/apperror"
	"spsu-chat/internal/models"

	"github.com/Masterminds/squirrel"
)

type MessagesPosgresql struct {
	db DB
}

func NewMessages(db DB) *MessagesPosgresql {
	return &MessagesPosgresql{
		db: db,
	}
}

func (m *MessagesPosgresql) Create(ctx context.Context, message models.CreateMessageRecord) error {
	query, args, _ := squirrel.
		Insert(MessagesTable).
		Columns(
			"chat_id",
			"user_id",
			"text",
			"created_at",
		).
		Values(
			message.ChatID,
			message.SenderID,
			message.Text,
			message.CreatedAt,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := m.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.NewDBError(
			err,
			"Message",
			"Create",
			query,
			args,
		)
	}

	return nil
}
func (m *MessagesPosgresql) GetByID(ctx context.Context, id int64) (models.Message, error) {
	query, args, _ := squirrel.
		Select("*").
		From(MessagesTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var message models.Message
	if err := m.db.GetContext(ctx, &message, query, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return message, apperror.ErrNotFound
		default:
			return message, apperror.NewDBError(
				err,
				"Message",
				"GetByID",
				query,
				args,
			)
		}
	}

	return message, nil
}
func (m *MessagesPosgresql) GetAll(ctx context.Context, pagination models.DBPagination, filters models.GetMessagesFilters) ([]models.Message, uint64, error) {
	// getting messages
	query := squirrel.
		Select("*").
		From(MessagesTable).
		Where(squirrel.Eq{"chat_id": filters.ChatID})

	queryString, args, _ := query.
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var count uint64
	var messages = make([]models.Message, 0)
	if err := m.db.SelectContext(ctx, &messages, queryString, args...); err != nil {
		return messages, count, apperror.NewDBError(
			err,
			"Message",
			"GetAll",
			queryString,
			args,
		)
	}

	// counting messages
	query = squirrel.
		Select("COUNT(*)").
		From(MessagesTable).
		Where(squirrel.Eq{"chat_id": filters.ChatID})

	queryString, args, _ = query.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err := m.db.GetContext(ctx, &count, queryString, args...); err != nil {
		return messages, count, apperror.NewDBError(
			err,
			"Message",
			"GetAll",
			queryString,
			args,
		)
	}

	return messages, count, nil
}
func (m *MessagesPosgresql) Delete(ctx context.Context, id int64) error {
	query, args, _ := squirrel.
		Delete(MessagesTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := m.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.NewDBError(
			err,
			"Message",
			"GetByID",
			query,
			args,
		)
	}

	return nil
}
