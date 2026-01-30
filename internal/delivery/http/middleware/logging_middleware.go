package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const (
	LoggerKey contextKey = "logger"
)

func LoggingMiddleware(baseLogger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		if !strings.HasPrefix(c.Path(), "/api") {
			return c.Next()
		}

		start := time.Now()
		requestID := uuid.NewString()

		logger := baseLogger.With(
			zap.String("transport", "http"),
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		ctx := context.WithValue(c.UserContext(), LoggerKey, logger)
		c.SetUserContext(ctx)

		logger.Info("request_started")

		err := c.Next()
		status := c.Response().StatusCode()
		duration := time.Since(start).Truncate(time.Millisecond)

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("duration", duration.String()),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		switch {
		case status >= 500:
			logger.Error("request_error", fields...)
		case status >= 400:
			logger.Warn("request_completed", fields...)
		default:
			logger.Info("request_completed", fields...)
		}

		return err
	}
}
