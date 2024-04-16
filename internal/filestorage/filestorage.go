package filestorage

import (
	"errors"

	"github.com/google/uuid"
)

const (
	MangaBucket = "manga"
)

var (
	ErrNotFound = errors.New("file not found")
)

type FileInfo struct {
	ID        uuid.UUID
	Extension string
	URL       string
}

type FileStorage interface {
	SaveFile(bucket string, filename string, data []byte) (FileInfo, error)
	GetFile(bucket string, id uuid.UUID) ([]byte, error)
	GetFileInfo(bucket string, id uuid.UUID) (FileInfo, error)
	GetFilesFromBucket(bucket string) ([]FileInfo, error)
	DeleteFile(bucket string, id uuid.UUID) error
	DeleteBucket(bucket string) error
}
