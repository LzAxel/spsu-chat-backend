package models

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/samber/lo"
)

const (
	MaxAvatarFileSize = (1 << 20) * 4 // 4 MB
)

var (
	ErrInvalidFileExt = errors.New("invalid file extension")
)

type UploadReader struct {
	Reader   io.Reader
	Filename string
}

type UploadFile struct {
	Data     []byte
	Filename string
}

// ValidateExtension checks if file extension is valid.
// Extensions format example: .png .jpg .gif
func ValidateExtension(filename string, extensions ...string) error {
	if !lo.Contains[string](extensions, filepath.Ext(filename)) {
		return fmt.Errorf("Allowed file extensions: %s", extensions)
	}

	return nil
}
