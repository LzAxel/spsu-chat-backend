package models

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type LoginUserInput struct {
	Username string
	Password string
}

func NewLoginUserInput(username string, password string) (LoginUserInput, error) {
	input := LoginUserInput{
		Username: username,
		Password: password,
	}

	return input, nil
}
