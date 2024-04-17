package service

import (
	"context"

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

	return user, handleNotFoundError(err, models.ErrUserNotFound)
}
func (u *UserService) GetByUsername(ctx context.Context, username string) (models.User, error) {
	user, err := u.repo.GetByUsername(ctx, username)

	return user, handleNotFoundError(err, models.ErrUserNotFound)
}
func (u *UserService) GetAll(ctx context.Context, pagination models.Pagination) ([]models.User, models.FullPagination, error) {
	users, total, err := u.repo.GetAll(ctx, models.DBPagination{
		Offset: pagination.Offset(),
		Limit:  pagination.Limit(),
	})

	return users, pagination.GetFull(total), err
}
