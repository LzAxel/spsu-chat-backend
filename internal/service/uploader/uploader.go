package uploader

import (
	"spsu-chat/internal/filestorage"
)

const (
	AvatarFolder = "avatar"
)

type Uploader struct {
	fileStorage filestorage.FileStorage
}

func NewUploader(fileStorage filestorage.FileStorage) *Uploader {
	return &Uploader{
		fileStorage: fileStorage,
	}
}
