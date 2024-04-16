package middleware

import (
	"time"

	"spsu-chat/internal/logger"

	"github.com/labstack/echo/v4"
)

func Logger(logger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			reqSize := req.Header.Get(echo.HeaderContentLength)
			if reqSize == "" {
				reqSize = "0"
			}

			logger.Debug("request", map[string]interface{}{
				"request_id":   id,
				"method":       req.Method,
				"uri":          req.RequestURI,
				"status":       res.Status,
				"ip":           c.RealIP(),
				"request_size": reqSize,
				"duration":     stop.Sub(start).String(),
				"referer":      req.Referer(),
				"user_agent":   req.UserAgent(),
			})
			return err
		}
	}
}
