package repository

import (
	"context"
	"spsu-chat/internal/logger"
	"spsu-chat/internal/models"
	"spsu-chat/internal/repository/postgresql"
)

type User interface {
	Create(ctx context.Context, user models.CreateUserRecord) error
	GetByID(ctx context.Context, id int64) (models.User, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
	GetAll(ctx context.Context, pagination models.DBPagination) ([]models.User, uint64, error)
}

type Chat interface {
	Create(ctx context.Context, chat models.CreateChatRecord) error
	GetByID(ctx context.Context, id int64) (models.Chat, error)
	GetAll(ctx context.Context, pagination models.DBPagination) ([]models.Chat, uint64, error)
}

type Repository struct {
	User
	Chat
}

func New(psql postgresql.PostgresqlRepository, logger logger.Logger) *Repository {
	return &Repository{
		User: postgresql.NewUser(psql.DB),
		Chat: postgresql.NewChat(psql.DB),
	}
}
