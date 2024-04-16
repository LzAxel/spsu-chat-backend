package config

import (
	"sync"

	"spsu-chat/internal/filestorage"
	"spsu-chat/internal/handlers/http"
	"spsu-chat/internal/jwt"
	"spsu-chat/internal/repository/postgresql"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "configs/dev.yaml"
)

type AppConfig struct {
	IsDev     bool   `yaml:"isDev" env:"IS_DEV"`
	IsTesting bool   `yaml:"isTesting" env:"IS_TESTING"`
	LogLevel  string `yaml:"logLevel" env:"LOG_LEVEL"`
}

type Config struct {
	Postgresql  postgresql.Config                  `yaml:"postgres"`
	Server      http.Config                        `yaml:"server"`
	App         AppConfig                          `yaml:"app"`
	JWT         jwt.Config                         `yaml:"jwt"`
	FileStorage filestorage.LocalFileStorageConfig `yaml:"fileStorage"`
}

var (
	config Config
	once   sync.Once
)

func ReadConfig() Config {
	once.Do(func() {
		if err := cleanenv.ReadConfig(configPath, &config); err != nil {
			panic(err)
		}
	})

	return config
}
