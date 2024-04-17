package service

import (
	"context"
	"spsu-chat/internal/filestorage"
	"spsu-chat/internal/jwt"
	"spsu-chat/internal/models"
	"spsu-chat/internal/repository"
	"spsu-chat/internal/service/uploader"
)

type Authorization interface {
	Login(ctx context.Context, input models.LoginUserInput) (jwt.TokenPair, error)
	Register(ctx context.Context, input models.CreateUserInput) error
	RefreshTokens(ctx context.Context, refreshToken string) (jwt.TokenPair, error)
}

type User interface {
	GetByID(ctx context.Context, id int64) (models.User, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
	GetAll(ctx context.Context, pagination models.Pagination) ([]models.User, models.FullPagination, error)
}

type Chat interface {
	GetAll(ctx context.Context, pagination models.Pagination) ([]models.Chat, models.FullPagination, error)
	GetByID(ctx context.Context, id int64) (models.Chat, error)
	Create(ctx context.Context, input models.CreateChatInput) error
}

type Services struct {
	User
	Authorization
	Chat
}

func New(
	ctx context.Context,
	repository *repository.Repository,
	jwt *jwt.JWT,
	fileStorage filestorage.FileStorage,
) *Services {
	_ = uploader.NewUploader(fileStorage)
	return &Services{
		User:          NewUserService(repository.User),
		Authorization: NewAuthorizationSerive(jwt, repository.User),
		Chat:          NewChatService(repository.Chat),
	}
}
