package models

import (
	"errors"
	"time"
)

const (
	ChatTypePublic = iota
	ChatTypePrivate

	MinChatNameLength = 5
)

var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrChatNameTooShort  = errors.New("chat name is too short")
	ErrChatNotPrivate    = errors.New("chat is not private")
	ErrChatWrongPassword = errors.New("wrong chat password")
	ErrChatAlreadyJoined = errors.New("you are already joined this chat")
)

type ChatType int8

type Chat struct {
	ID           int64     `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	CreatorID    int64     `db:"creator_id" json:"creator_id"`
	Type         ChatType  `db:"type" json:"type"`
	PasswordHash []byte    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type CreateChatInput struct {
	Name      string
	CreatorID int64
	Password  *string
}

func NewCreateChatInput(name string, creatorID int64, password *string) (CreateChatInput, error) {
	if len(name) < MinChatNameLength {
		return CreateChatInput{}, ErrChatNameTooShort
	}
	return CreateChatInput{
		Name:      name,
		CreatorID: creatorID,
		Password:  password,
	}, nil
}

type CreateChatRecord struct {
	Name         string
	CreatorID    int64
	Type         int8
	PasswordHash []byte
	CreatedAt    time.Time
}
