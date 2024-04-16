package http

import (
	"errors"
	"net/http"
	"strconv"

	"spsu-chat/internal/models"

	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidPage      = errors.New("Invalid page query parameter")
	ErrInvalidPageLimit = errors.New("Invalid page_limit query parameter")
)

func (h *Handler) WithPagination() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			pageParam := c.QueryParam("page")
			if pageParam == "" {
				pageParam = "1"
			}
			page, err := strconv.ParseUint(pageParam, 10, 64)
			if err != nil {
				return h.newValidationErrorResponse(c, http.StatusBadRequest, ErrInvalidPage)
			}

			pageLimitParam := c.QueryParam("page_limit")
			if pageLimitParam == "" {
				pageLimitParam = strconv.Itoa(models.DefaultLimit)
			}
			pageLimit, err := strconv.ParseUint(pageLimitParam, 10, 64)
			if err != nil {
				return h.newValidationErrorResponse(c, http.StatusBadRequest, ErrInvalidPageLimit)
			}

			pagination, err := models.NewPagination(page, pageLimit)
			if err != nil {
				return h.newAppErrorResponse(c, err)
			}
			c.Set("pagination", pagination)

			if err := next(c); err != nil {
				c.Error(err)
			}
			h.logger.Debug("pagination", map[string]interface{}{
				"page":       page,
				"pageLimit":  pageLimit,
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			})

			return nil
		}
	}
}

func getPaginationFromContext(ctx echo.Context) (models.Pagination, error) {
	pagination, ok := ctx.Get("pagination").(models.Pagination)
	if !ok {
		pagination = models.Pagination{
			Page:      1,
			PageLimit: models.DefaultLimit,
		}
	}

	if pagination.PageLimit > models.MaxLimit {
		return models.Pagination{}, models.ErrInvalidLimit
	}

	return pagination, nil
}
