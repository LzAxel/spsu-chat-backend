package app

import (
	"context"

	"spsu-chat/internal/config"
	"spsu-chat/internal/filestorage"
	"spsu-chat/internal/handlers/http"
	"spsu-chat/internal/jwt"
	"spsu-chat/internal/logger"
	"spsu-chat/internal/repository"
	"spsu-chat/internal/repository/postgresql"
	"spsu-chat/internal/service"
	"spsu-chat/pkg/clock"
	"spsu-chat/pkg/uuid"
)

type App struct {
	services   *service.Services
	repository *repository.Repository
	hanlder    *http.Handler
	logger     logger.Logger
}

func New(config config.Config) *App {
	ctx := context.Background()
	logger := logger.NewLogrusLogger(config.App.LogLevel, config.App.IsDev)
	logger.Infof("config loaded")

	if config.App.IsTesting {
		clock.InitClock(true)
		uuid.InitUUID(true)
	}

	logger.Infof("connecting to postgresql on %s:%d", config.Postgresql.Host, config.Postgresql.Port)
	psql, err := postgresql.New(config.Postgresql)
	if err != nil {
		logger.Fatalf("connect to postgresql: %s", err)
	}

	err = postgresql.Migrate(ctx, config.Postgresql)
	if err != nil {
		logger.Warnf("migrate database: %s", err)
	}
	fileStorage := filestorage.NewLocalFileStorage(config.FileStorage, logger)
	go fileStorage.Serve()

	jwt := jwt.New(config.JWT)
	repository := repository.New(psql, logger)
	services := service.New(ctx, repository, jwt, fileStorage)
	handler := http.New(config.Server, services, logger, jwt)

	return &App{
		services:   services,
		repository: repository,
		hanlder:    handler,
		logger:     logger,
	}
}

func (app *App) Start() error {
	if err := app.hanlder.Start(); err != nil {
		app.logger.Errorf("failed to start app: %s", err)
		return err
	}

	return nil
}

func (app *App) Shutdown(context.Context) {
	if err := app.hanlder.Stop(context.Background()); err != nil {
		panic(err)
	}
}
