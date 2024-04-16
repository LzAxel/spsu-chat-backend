package service

import (
	"context"
	"errors"

	"spsu-chat/internal/apperror"
	"spsu-chat/internal/models"
	"spsu-chat/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repository repository.User) *UserService {
	return &UserService{
		repo: repository,
	}
}

func (u *UserService) GetByID(ctx context.Context, id int64) (models.User, error) {
	user, err := u.repo.GetByID(ctx, id)

	return user, handleNotFoundError(err)
}
func (u *UserService) GetByUsername(ctx context.Context, username string) (models.User, error) {
	user, err := u.repo.GetByUsername(ctx, username)

	return user, handleNotFoundError(err)
}
func (u *UserService) GetAll(ctx context.Context, pagination models.Pagination) ([]models.User, models.FullPagination, error) {
	users, total, err := u.repo.GetAll(ctx, models.DBPagination{
		Offset: pagination.Offset(),
		Limit:  pagination.Limit(),
	})

	return users, pagination.GetFull(total), err
}

func handleNotFoundError(err error) error {
	if errors.Is(err, apperror.ErrNotFound) {
		return models.ErrUserNotFound
	}
	return err
}
