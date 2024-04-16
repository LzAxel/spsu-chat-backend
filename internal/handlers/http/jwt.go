package http

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"spsu-chat/internal/jwt"
	"spsu-chat/internal/models"

	"github.com/labstack/echo/v4"
)

type JWTValidator interface {
	ValidateToken(token string) (jwt.Claims, error)
}

func (h *Handler) Authorized() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			splitAuth := strings.Split(authHeader, " ")
			if len(splitAuth) != 2 {
				return h.newAuthErrorResponse(c, http.StatusUnauthorized, jwt.ErrInvalidToken)
			}

			claims, err := h.jwtValidator.ValidateToken(splitAuth[1])
			if err != nil {
				return h.newAuthErrorResponse(c, http.StatusUnauthorized, err)
			}
			user, err := h.services.User.GetByID(c.Request().Context(), claims.Subject)
			if err != nil {
				return h.newAppErrorResponse(c, errors.New("failed to get user type"))
			}

			c.Set("user", user)

			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}

func (h *Handler) RequireUserType(userTypes ...models.UserType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(models.User)
			if !ok {
				return h.newAppErrorResponse(c, errors.New("Handler.RequireUserType:failed to get user from context (forgot to use Authorize middleware)"))
			}

			if !slices.Contains(userTypes, user.Type) {
				return h.newAuthErrorResponse(c, http.StatusForbidden, errors.New("access denied"))
			}

			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}
