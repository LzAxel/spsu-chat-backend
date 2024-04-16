package http

import (
	"errors"
	"net/http"
	"strconv"

	"spsu-chat/internal/models"

	"github.com/labstack/echo/v4"
)

type getAllUsersResponse struct {
	Users      []models.User         `json:"users"`
	Pagination models.FullPagination `json:"pagination"`
}

func (h *Handler) getAllUsers(ctx echo.Context) error {
	reqPagination, err := getPaginationFromContext(ctx)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}
	users, pagination, err := h.services.User.GetAll(ctx.Request().Context(), reqPagination)
	if err != nil {
		return h.newAppErrorResponse(ctx, err)
	}

	ctx.JSON(http.StatusOK, getAllUsersResponse{Users: users, Pagination: pagination})

	return nil
}

type getUserResponse struct {
	User models.User `json:"user"`
}

func (h *Handler) getUserByID(ctx echo.Context) error {
	id := ctx.Param("id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid user id"))
	}

	user, err := h.services.User.GetByID(ctx.Request().Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, err.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, getUserResponse{User: user})

	return nil
}

func (h *Handler) getSelfUser(ctx echo.Context) error {
	user, ok := ctx.Get("user").(models.User)
	if !ok {
		return h.newAppErrorResponse(ctx, errors.New("failed to get user from context"))
	}

	user, err := h.services.User.GetByID(ctx.Request().Context(), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			return h.newErrorResponse(ctx, http.StatusNotFound, err.Error())
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, getUserResponse{User: user})

	return nil
}
