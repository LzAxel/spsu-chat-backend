package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"spsu-chat/internal/apperror"
	"spsu-chat/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
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

func (p *ChatPosgresql) JoinUser(ctx context.Context, chatID int64, userID int64) error {
	query, args, _ := squirrel.
		Insert(ChatUsersTable).
		Columns(
			"chat_id",
			"user_id",
		).
		Values(
			chatID,
			userID,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		pgErr := GetPgError(err)

		switch {
		case pgErr != nil && pgErr.Code == pgerrcode.UniqueViolation:
			return models.ErrChatAlreadyJoined
		case errors.Is(err, sql.ErrNoRows):
			return apperror.ErrNotFound
		}

		return apperror.NewDBError(
			err,
			"Chat",
			"JoinUser",
			query,
			args,
		)
	}

	return nil
}
func (p *ChatPosgresql) LeaveUser(ctx context.Context, chatID int64, userID int64) error {
	query, args, _ := squirrel.
		Delete(ChatUsersTable).
		Where(squirrel.Eq{"chat_id": chatID, "user_id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.NewDBError(
			err,
			"Chat",
			"LeaveUser",
			query,
			args,
		)
	}

	return nil
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

func (p *ChatPosgresql) IsUserInChat(ctx context.Context, chatID, userID int64) (bool, error) {
	query, args, _ := squirrel.
		Select("user_id").
		From(ChatUsersTable).
		Where(squirrel.Eq{"user_id": userID, "chat_id": chatID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var returnedUserID int64
	if err := p.db.GetContext(ctx, &returnedUserID, query, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, nil
		default:
			return false, apperror.NewDBError(
				err,
				"Chat",
				"GetByID",
				query,
				args,
			)
		}
	}

	return true, nil
}
