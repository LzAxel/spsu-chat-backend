package service

import (
	"context"
	"errors"

	"spsu-chat/internal/apperror"
	"spsu-chat/internal/models"
	"spsu-chat/internal/repository"
	"spsu-chat/pkg/clock"
	"spsu-chat/pkg/hash"
)

type ChatService struct {
	repo repository.Chat
}

func NewChatService(repository repository.Chat) *ChatService {
	return &ChatService{
		repo: repository,
	}
}

func (c *ChatService) Create(ctx context.Context, input models.CreateChatInput) error {
	var (
		passwordHash []byte
		chatType     models.ChatType = models.ChatTypePublic
		err          error
	)

	if input.Password != nil {
		passwordHash, err = hash.Hash(*input.Password)
		if err != nil {
			return err
		}
		chatType = models.ChatTypePrivate
	}

	chat := models.CreateChatRecord{
		Name:         input.Name,
		Type:         int8(chatType),
		CreatorID:    input.CreatorID,
		PasswordHash: passwordHash,
		CreatedAt:    clock.Now(),
	}

	return c.repo.Create(ctx, chat)
}

func (c *ChatService) GetByID(ctx context.Context, id int64) (models.Chat, error) {
	chat, err := c.repo.GetByID(ctx, id)

	return chat, handleNotFoundError(err, models.ErrChatNotFound)
}

func (c *ChatService) GetAll(ctx context.Context, pagination models.Pagination) ([]models.Chat, models.FullPagination, error) {
	chats, total, err := c.repo.GetAll(ctx, models.DBPagination{
		Offset: pagination.Offset(),
		Limit:  pagination.Limit(),
	})

	return chats, pagination.GetFull(total), err
}

func (c *ChatService) JoinUser(ctx context.Context, chatID int64, userID int64, password string) error {
	chat, err := c.repo.GetByID(ctx, chatID)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return models.ErrChatNotFound
		}
		return err
	}

	if chat.Type != models.ChatTypePrivate {
		return models.ErrChatNotPrivate
	}

	if err := hash.Compare(chat.PasswordHash, password); err != nil {
		return models.ErrChatWrongPassword
	}

	return c.repo.JoinUser(ctx, chatID, userID)
}
func (c *ChatService) LeaveUser(ctx context.Context, chatID int64, userID int64) error {
	return c.repo.LeaveUser(ctx, chatID, userID)
}
