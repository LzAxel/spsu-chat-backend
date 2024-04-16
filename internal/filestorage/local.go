package filestorage

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"spsu-chat/internal/logger"

	"github.com/google/uuid"
)

type LocalFileStorageConfig struct {
	StoragePath string `yaml:"storagePath" env:"LOCAL_FILE_STORAGE_PATH"`
	Host        string `yaml:"host" env:"LOCAL_FILE_STORAGE_HOST"`
	Port        uint   `yaml:"port" env:"LOCAL_FILE_STORAGE_PORT"`
}

type LocalFileStorage struct {
	logger      logger.Logger
	storagePath string
	host        string
	port        uint
}

func NewLocalFileStorage(config LocalFileStorageConfig, logger logger.Logger) *LocalFileStorage {
	return &LocalFileStorage{
		storagePath: config.StoragePath,
		host:        config.Host,
		port:        config.Port,
		logger:      logger,
	}
}

func (s *LocalFileStorage) SaveFile(bucket string, filename string, data []byte) (FileInfo, error) {
	id := uuid.New()
	fileExtansion := filepath.Ext(filename)

	filePath := filepath.Join(s.storagePath, bucket, id.String()+fileExtansion)
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)

	var fileInfo = FileInfo{
		ID:        id,
		Extension: fileExtansion,
		URL:       s.formatFileURL(bucket, filePath),
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fileInfo, fmt.Errorf("LocalFileStorage.SaveFile: creating file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fileInfo, fmt.Errorf("LocalFileStorage.SaveFile: writing file: %w", err)
	}

	return fileInfo, nil
}

func (s *LocalFileStorage) GetFile(bucket string, id uuid.UUID) ([]byte, error) {
	filepath, err := s.findFilePathByUUID(bucket, id)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("LocalFileStorage.GetFile: reading file: %w", err)
	}

	return data, nil
}

func (s *LocalFileStorage) DeleteFile(bucket string, id uuid.UUID) error {
	filepath, err := s.findFilePathByUUID(bucket, id)
	if err != nil {
		return err
	}

	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("LocalFileStorage.DeleteFile: removing file: %w", err)
	}

	return nil
}

func (s *LocalFileStorage) GetFileInfo(bucket string, id uuid.UUID) (FileInfo, error) {
	var fileInfo = FileInfo{
		ID: id,
	}

	filePath, err := s.findFilePathByUUID(bucket, id)
	if err != nil {
		return fileInfo, fmt.Errorf("LocalFileStorage.GetFileURL: %w", err)
	}
	fileInfo.Extension = filepath.Ext(filePath)
	fileInfo.URL = s.formatFileURL(bucket, filePath)

	return fileInfo, nil
}

func (s *LocalFileStorage) GetFilesFromBucket(bucket string) ([]FileInfo, error) {
	files, err := os.ReadDir(filepath.Join(s.storagePath, bucket))
	if err != nil {
		return nil, fmt.Errorf("LocalFileStorage.GetFilesFromBucket: %w", err)
	}

	var fileInfos = make([]FileInfo, len(files))

	for i, file := range files {
		ext := filepath.Ext(file.Name())
		id, err := uuid.Parse(file.Name()[:len(file.Name())-len(ext)])
		if err != nil {
			return nil, fmt.Errorf("LocalFileStorage.GetFilesFromBucket: parsing id from filename:%w", err)
		}
		fileInfos[i] = FileInfo{
			ID:        id,
			Extension: filepath.Ext(file.Name()),
			URL:       s.formatFileURL(bucket, filepath.Join(bucket, file.Name())),
		}
	}

	return fileInfos, nil
}

func (s *LocalFileStorage) findFilePathByUUID(bucket string, id uuid.UUID) (string, error) {
	pattern := filepath.Join(s.storagePath, bucket, id.String()) + ".*"
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("LocalFileStorage.findFilePathByUUID: %w: paths=%s", err, paths)
	}
	if len(paths) != 1 {
		return "", fmt.Errorf("LocalFileStorage.findFilePathByUUID: %w: pattern=%s", ErrNotFound, pattern)
	}

	return paths[0], nil
}

func (s *LocalFileStorage) DeleteBucket(bucket string) error {
	err := os.RemoveAll(filepath.Join(s.storagePath, bucket))
	if err != nil {
		return fmt.Errorf("LocalFileStorage.DeleteBucket: %w", err)
	}
	return nil
}

func (s *LocalFileStorage) Serve() {
	s.logger.Infof("starting local file storage server on %s:%d (path: %s)", s.host, s.port, s.storagePath)
	http.Handle("/", http.FileServer(http.Dir(s.storagePath)))
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port),
		nil,
	)
}

func (s *LocalFileStorage) formatFileURL(bucket string, filepath string) string {
	return fmt.Sprintf("http://%s:%d/", s.host, s.port) + path.Join(bucket, path.Base(filepath))
}
