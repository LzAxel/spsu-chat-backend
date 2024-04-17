package service

import (
	"errors"
	"spsu-chat/internal/apperror"
)

func handleNotFoundError(err error, newErr error) error {
	if errors.Is(err, apperror.ErrNotFound) {
		return newErr
	}
	return err
}
