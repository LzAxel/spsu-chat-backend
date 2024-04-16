package models

import (
	"errors"
	"time"
)

const (
	UserTypeUser = iota
	UserTypeAdmin
)

const (
	minPasswordLength = 6
	minUsernameLength = 4
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUsername = errors.New("invalid username")
	ErrUsernameExists  = errors.New("username already exists")
	ErrUserNotFound    = errors.New("user not found")
)

type UserType int8

type User struct {
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash []byte    `db:"password_hash" json:"-"`
	DisplayName  string    `db:"display_name" json:"display_name"`
	Type         UserType  `db:"type" json:"type"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type CreateUserInput struct {
	Username    string
	Password    string
	DisplayName *string
}

func NewCreateUserInput(username string, displayName *string, password string) (CreateUserInput, error) {
	if len(password) < minPasswordLength {
		return CreateUserInput{}, ErrInvalidPassword
	}

	if len(username) < minUsernameLength {
		return CreateUserInput{}, ErrInvalidUsername
	}

	return CreateUserInput{
		Username:    username,
		DisplayName: displayName,
		Password:    password,
	}, nil
}

type CreateUserRecord struct {
	Username     string
	DisplayName  string
	PasswordHash []byte
	Type         int8
	CreatedAt    time.Time
}
