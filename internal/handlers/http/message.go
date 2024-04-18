package http

import (
	"errors"
	"net/http"
	"spsu-chat/internal/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type getAllMessagesResponse struct {
	Messages   []models.Message      `json:"messages"`
	Pagination models.FullPagination `json:"pagination"`
}

func (h *Handler) getAllMessages(ctx echo.Context) error {
	var filters models.GetMessagesFilters

	err := ctx.Bind(&filters)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid filters"))
	}
	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return h.newAppErrorResponse(ctx, errors.New("invalid user in context"))
	}

	reqPagination, err := getPaginationFromContext(ctx)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}
	messages, count, err := h.services.Message.GetAll(ctx.Request().Context(), reqPagination, filters, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrChatNotJoined):
			return h.newErrorResponse(ctx, http.StatusForbidden, models.ErrChatNotJoined.Error())
		}
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.JSON(http.StatusOK, getAllMessagesResponse{
		Messages:   messages,
		Pagination: reqPagination.GetFull(count),
	})

	return nil
}

type sendMessageRequest struct {
	Text   string `json:"text"`
	ChatID int64  `json:"chat_id"`
}

func (h *Handler) SendMessage(ctx echo.Context) error {
	var message sendMessageRequest

	err := ctx.Bind(&message)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}
	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return h.newAppErrorResponse(ctx, errors.New("invalid user in context"))
	}

	input := models.NewCreateMessageInput(message.ChatID, user.ID, message.Text)

	err = h.services.Message.Create(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrChatNotJoined):
			return h.newErrorResponse(ctx, http.StatusForbidden, models.ErrChatNotJoined.Error())
		case errors.Is(err, models.ErrChatNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, models.ErrChatNotFound.Error())
		}
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.NoContent(http.StatusCreated)

	return nil
}

func (h *Handler) DeleteMessage(ctx echo.Context) error {
	messageID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid message id"))
	}

	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return h.newAppErrorResponse(ctx, errors.New("invalid user in context"))
	}

	err = h.services.Message.Delete(ctx.Request().Context(), user.ID, messageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotYourMessage):
			return h.newErrorResponse(ctx, http.StatusForbidden, models.ErrNotYourMessage.Error())
		}
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.NoContent(http.StatusCreated)

	return nil
}
