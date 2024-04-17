package http

import (
	"context"
	"net"
	"strconv"

	middle "spsu-chat/internal/handlers/http/middleware"
	"spsu-chat/internal/logger"
	"spsu-chat/internal/models"
	"spsu-chat/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Host string `yaml:"host" env:"HOST"`
	Port uint   `yaml:"port" env:"PORT"`
}
type Handler struct {
	jwtValidator JWTValidator
	services     *service.Services
	server       *echo.Echo
	config       Config
	logger       logger.Logger
}

func New(config Config, services *service.Services, logger logger.Logger, jwtValidator JWTValidator) *Handler {
	echo := echo.New()
	echo.HideBanner = true
	echo.HidePort = true
	handler := Handler{
		server:       echo,
		config:       config,
		services:     services,
		logger:       logger,
		jwtValidator: jwtValidator,
	}
	handler.initMiddlewares()
	handler.initRoutes()

	return &handler
}

func (h *Handler) initMiddlewares() {
	h.server.Use(
		middleware.RequestID(),
		middleware.Recover(),
		middle.Logger(h.logger),
	)
}

func (h *Handler) initRoutes() {
	api := h.server.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	auth := v1.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refreshTokens)
	}

	user := v1.Group("/users", h.Authorized())
	{
		user.GET("", h.getAllUsers, h.WithPagination(), h.RequireUserType(models.UserTypeAdmin))
		user.GET("/:id", h.getUserByID)
		user.GET("/self", h.getSelfUser)
	}

	chat := v1.Group("/chats", h.Authorized())
	{
		chat.GET("", h.getAllChats, h.WithPagination())
		chat.POST("", h.createChat, h.RequireUserType(models.UserTypeAdmin))
		chat.GET("/:id", h.getChatByID)
	}
}

func (h *Handler) Stop(ctx context.Context) error {
	h.logger.Infof("shutting down server")
	return h.server.Shutdown(ctx)
}

func (h *Handler) Start() error {
	h.logger.Infof("starting server on %s:%d", h.config.Host, h.config.Port)
	return h.server.Start(net.JoinHostPort(h.config.Host, strconv.Itoa(int(h.config.Port))))
}
