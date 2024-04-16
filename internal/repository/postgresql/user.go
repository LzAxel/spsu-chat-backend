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

type UserPosgresql struct {
	db DB
}

func NewUser(db DB) *UserPosgresql {
	return &UserPosgresql{
		db: db,
	}
}

func (p *UserPosgresql) Create(ctx context.Context, user models.CreateUserRecord) error {
	query, args, _ := squirrel.
		Insert(UsersTable).
		Columns(
			"username",
			"display_name",
			"password_hash",
			"type",
			"created_at",
		).
		Values(
			user.Username,
			user.DisplayName,
			user.PasswordHash,
			user.Type,
			user.CreatedAt,
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := p.db.ExecContext(ctx, query, args...); err != nil {
		pgErr := GetPgError(err)

		switch {
		case pgErr != nil && pgErr.Code == pgerrcode.UniqueViolation:
			return models.ErrUsernameExists
		case errors.Is(err, sql.ErrNoRows):
			return apperror.ErrNotFound
		}

		return apperror.NewDBError(
			err,
			"User",
			"Create",
			query,
			args,
		)
	}

	return nil
}

func (p *UserPosgresql) GetByID(ctx context.Context, id int64) (models.User, error) {
	query, args, _ := squirrel.
		Select("*").
		From(UsersTable).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var user models.User
	if err := p.db.GetContext(ctx, &user, query, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return user, apperror.ErrNotFound
		default:
			return user, apperror.NewDBError(
				err,
				"User",
				"GetByID",
				query,
				args,
			)
		}
	}

	return user, nil
}
func (p *UserPosgresql) GetByUsername(ctx context.Context, username string) (models.User, error) {
	query, args, _ := squirrel.
		Select("*").
		From(UsersTable).
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var user models.User
	if err := p.db.GetContext(ctx, &user, query, args...); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return user, apperror.ErrNotFound
		}
		return user, apperror.NewDBError(
			err,
			"User",
			"GetByUsername",
			query,
			args,
		)
	}

	return user, nil
}

func (p *UserPosgresql) GetAll(ctx context.Context, pagination models.DBPagination) ([]models.User, uint64, error) {
	// getting users
	query := squirrel.
		Select("*").
		From(UsersTable)

	queryString, args, _ := query.Limit(pagination.Limit).
		Offset(pagination.Offset).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var count uint64
	var users = make([]models.User, 0)
	if err := p.db.SelectContext(ctx, &users, queryString, args...); err != nil {
		return users, count, apperror.NewDBError(
			err,
			"User",
			"GetAll",
			queryString,
			args,
		)
	}
	// counting users
	query = squirrel.
		Select("COUNT(*)").
		From(UsersTable)

	queryString, args, _ = query.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err := p.db.GetContext(ctx, &count, queryString, args...); err != nil {
		return users, count, apperror.NewDBError(
			err,
			"User",
			"GetAll",
			queryString,
			args,
		)
	}

	return users, count, nil
}
