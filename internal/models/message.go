package models

import (
	"errors"
	"time"
)

var (
	ErrNotYourMessage  = errors.New("message is not your")
	ErrMessageNotFound = errors.New("message not found")
)

// BASE MODEL
type Message struct {
	ID        int64     `db:"id" json:"id"`
	ChatID    int64     `db:"chat_id" json:"chat_id"`
	SenderID  int64     `db:"user_id" json:"sender_id"`
	Text      string    `db:"text" json:"text"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// CREATE MODELS
type CreateMessageInput struct {
	ChatID   int64
	SenderID int64
	Text     string
}

func NewCreateMessageInput(chatID, senderID int64, text string) CreateMessageInput {
	return CreateMessageInput{
		ChatID:   chatID,
		SenderID: senderID,
		Text:     text,
	}
}

type CreateMessageRecord struct {
	ChatID    int64
	SenderID  int64
	Text      string
	CreatedAt time.Time
}

// FILTER MODELS
type GetMessagesFilters struct {
	ChatID int64 `query:"chat_id"`
}
