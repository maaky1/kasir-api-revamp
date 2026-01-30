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
			zap.String("layer", "http"),
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		ctx := context.WithValue(context.Background(), LoggerKey, logger)
		c.Locals(string(LoggerKey), ctx)

		logger.Info("request_started")

		err := c.Next()
		status := c.Response().StatusCode()
		duration := time.Since(start).Truncate(time.Millisecond)

		switch {
		case status >= 500:
			logger.Error("request_error",
				zap.Int("status", status),
				zap.String("duration", duration.String()),
				zap.Error(err),
			)

		case status >= 400:
			logger.Warn("request_completed",
				zap.Int("status", status),
				zap.String("duration", duration.String()),
			)

		default:
			logger.Info("request_completed",
				zap.Int("status", status),
				zap.String("duration", duration.String()),
			)
		}

		return err
	}
}
