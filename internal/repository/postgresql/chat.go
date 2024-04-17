package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"spsu-chat/internal/apperror"
	"spsu-chat/internal/models"

	"github.com/Masterminds/squirrel"
)

type ChatPosgresql struct {
	db DB
}

func NewChat(db DB) *ChatPosgresql {
	return &ChatPosgresql{
		db: db,
	}
}

func (p *ChatPosgresql) Create(ctx context.Context, chat models.CreateChatRecord) error {
	query, args, _ := squirrel.
		Insert(ChatsTable).
		Columns(
			"name",
			"creator_id",
			"type",
			"password_hash",
			"created_at",
		).
		Values(
			chat.Name,
			chat.CreatorID,
			chat.Type,
			chat.PasswordHash,
			chat.CreatedAt,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.ErrNotFound
		}

		return apperror.NewDBError(
			err,
			"Chat",
			"Create",
			query,
			args,
		)
	}

	return nil
}

func (p *ChatPosgresql) GetByID(ctx context.Context, id int64) (models.Chat, error) {
	query, args, _ := squirrel.
		Select("*").
		From(ChatsTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var chat models.Chat
	if err := p.db.GetContext(ctx, &chat, query, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return chat, apperror.ErrNotFound
		default:
			return chat, apperror.NewDBError(
				err,
				"Chat",
				"GetByID",
				query,
				args,
			)
		}
	}

	return chat, nil
}

func (p *ChatPosgresql) GetAll(ctx context.Context, pagination models.DBPagination) ([]models.Chat, uint64, error) {
	// getting chats
	query := squirrel.
		Select("*").
		From(ChatsTable)

	queryString, args, _ := query.Limit(pagination.Limit).
		Offset(pagination.Offset).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var count uint64
	var chats = make([]models.Chat, 0)
	if err := p.db.SelectContext(ctx, &chats, queryString, args...); err != nil {
		return chats, count, apperror.NewDBError(
			err,
			"Chat",
			"GetAll",
			queryString,
			args,
		)
	}
	// counting chats
	query = squirrel.
		Select("COUNT(*)").
		From(ChatsTable)

	queryString, args, _ = query.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err := p.db.GetContext(ctx, &count, queryString, args...); err != nil {
		return chats, count, apperror.NewDBError(
			err,
			"Chat",
			"GetAll",
			queryString,
			args,
		)
	}

	return chats, count, nil
}
