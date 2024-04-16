package uploader

import (
	"context"
	"fmt"
	"path"

	"spsu-chat/internal/filestorage"
	"spsu-chat/internal/models"
)

func (u *Uploader) UploadAvatar(
	ctx context.Context,
	userID int64,
	file models.UploadFile) (filestorage.FileInfo, error) {

	pageFileInfo, err := u.fileStorage.SaveFile(u.formatAvatarFolder(userID), file.Filename, file.Data)
	if err != nil {
		return filestorage.FileInfo{}, fmt.Errorf("Uploader.UploadAvatar: %w", err)
	}

	return pageFileInfo, nil
}

func (u *Uploader) formatAvatarFolder(userID int64) string {
	return path.Join(AvatarFolder, fmt.Sprintf("%d", userID))
}
