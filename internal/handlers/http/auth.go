package http

import (
	"errors"
	"net/http"

	"spsu-chat/internal/jwt"
	"spsu-chat/internal/models"

	"github.com/labstack/echo/v4"
)

type signUpRequest struct {
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name,omitempty"`
	Password    string  `json:"password"`
}

func (h *Handler) signUp(ctx echo.Context) error {
	var req signUpRequest

	if err := ctx.Bind(&req); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid input"))
	}

	input, err := models.NewCreateUserInput(
		req.Username,
		req.DisplayName,
		req.Password,
	)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	err = h.services.Authorization.Register(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUsernameExists):
			return h.newAuthErrorResponse(ctx, http.StatusConflict, err)
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.NoContent(http.StatusCreated)

	return nil
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) signIn(ctx echo.Context) error {
	var req signInRequest

	if err := ctx.Bind(&req); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid input"))
	}

	input, err := models.NewLoginUserInput(
		req.Username,
		req.Password,
	)
	if err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, err)
	}

	tokenPair, err := h.services.Authorization.Login(ctx.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			return h.newAuthErrorResponse(ctx, http.StatusUnauthorized, err)
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})

	return nil
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) refreshTokens(ctx echo.Context) error {
	var req refreshRequest

	if err := ctx.Bind(&req); err != nil {
		return h.newValidationErrorResponse(ctx, http.StatusBadRequest, jwt.ErrInvalidToken)
	}

	tokenPair, err := h.services.Authorization.RefreshTokens(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrInvalidToken):
			return h.newValidationErrorResponse(ctx, http.StatusBadRequest, jwt.ErrInvalidToken)
		case errors.Is(err, jwt.ErrTokenExpired):
			return h.newAuthErrorResponse(ctx, http.StatusUnauthorized, jwt.ErrTokenExpired)
		default:
			return h.newAppErrorResponse(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})

	return nil
}
