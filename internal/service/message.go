package service

import (
	"context"
	"errors"
	"spsu-chat/internal/apperror"
	"spsu-chat/internal/models"
	"spsu-chat/internal/repository"
	"spsu-chat/pkg/clock"
)

type MessageService struct {
	repo     repository.Message
	chatRepo repository.Chat
}

func NewMessageService(repo repository.Message, chatRepo repository.Chat) *MessageService {
	return &MessageService{
		repo:     repo,
		chatRepo: chatRepo,
	}
}

func (m *MessageService) Create(ctx context.Context, message models.CreateMessageInput) error {
	chat, err := m.chatRepo.GetByID(ctx, message.ChatID)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return models.ErrChatNotFound
		}
		return err
	}

	if chat.Type == models.ChatTypePrivate {
		isJoined, err := m.chatRepo.IsUserInChat(ctx, chat.ID, message.SenderID)
		if err != nil {
			return err
		}
		if !isJoined {
			return models.ErrChatNotJoined
		}
	}

	input := models.CreateMessageRecord{
		ChatID:    message.ChatID,
		SenderID:  message.SenderID,
		Text:      message.Text,
		CreatedAt: clock.Now(),
	}
	return m.repo.Create(ctx, input)
}
func (m *MessageService) GetAll(ctx context.Context, pagination models.Pagination, filters models.GetMessagesFilters, userID int64) ([]models.Message, uint64, error) {
	chat, err := m.chatRepo.GetByID(ctx, filters.ChatID)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, 0, models.ErrChatNotFound
		}
		return nil, 0, err
	}

	if chat.Type == models.ChatTypePrivate {
		isJoined, err := m.chatRepo.IsUserInChat(ctx, chat.ID, userID)
		if err != nil {
			return nil, 0, err
		}
		if !isJoined {
			return nil, 0, models.ErrChatNotJoined
		}
	}

	return m.repo.GetAll(ctx, models.DBPagination{
		Offset: pagination.Offset(),
		Limit:  pagination.Limit(),
	}, filters)
}
func (m *MessageService) Delete(ctx context.Context, userID int64, messageID int64) error {
	message, err := m.repo.GetByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return models.ErrMessageNotFound
		}
		return err
	}
	if message.SenderID != userID {
		return models.ErrNotYourMessage
	}

	return m.repo.Delete(ctx, messageID)
}
