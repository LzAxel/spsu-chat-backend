package models

import "time"

const (
	ChatTypePublic = iota
	ChatTypePrivate
)

type ChatType int8

type Chat struct {
	ID           int64
	Name         string
	CreatorID    int64
	Type         ChatType
	PasswordHash []byte
	CreatedAt    time.Time
}
