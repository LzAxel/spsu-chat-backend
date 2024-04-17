package http

import (
	"errors"
	"net/http"
	"spsu-chat/internal/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type createChatRequest struct {
	Name     string  `json:"name"`
	Password *string `json:"password"`
}

func (h *Handler) createChat(ctx echo.Context) error {
	var req createChatRequest
	if err := ctx.Bind(&req); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return h.newAppErrorResponse(ctx, errors.New("invalid user in context"))
	}

	input, err := models.NewCreateChatInput(req.Name, user.ID, req.Password)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	err = h.services.Chat.Create(ctx.Request().Context(), input)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.NoContent(http.StatusCreated)
	return nil
}

type getAllChatsResponse struct {
	Chats      []models.Chat         `json:"chats"`
	Pagination models.FullPagination `json:"pagination"`
}

func (h *Handler) getAllChats(ctx echo.Context) error {
	reqPagination, err := getPaginationFromContext(ctx)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}
	chats, pagination, err := h.services.Chat.GetAll(ctx.Request().Context(), reqPagination)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.JSON(http.StatusOK, getAllChatsResponse{Chats: chats, Pagination: pagination})

	return nil
}

type getChatResponse struct {
	Chat models.Chat `json:"chat"`
}

func (h *Handler) getChatByID(ctx echo.Context) error {
	id := ctx.Param("id")
	chatID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid chat id"))
	}

	chat, err := h.services.Chat.GetByID(ctx.Request().Context(), chatID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrChatNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, err.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, getChatResponse{Chat: chat})

	return nil
}
