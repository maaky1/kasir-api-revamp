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

		// 🚫 skip non-API
		if !strings.HasPrefix(c.Path(), "/api") {
			return c.Next()
		}

		start := time.Now()
		requestID := uuid.NewString()

		// request-scoped logger
		logger := baseLogger.With(
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		// simpan context + logger
		ctx := context.WithValue(context.Background(), LoggerKey, logger)
		c.Locals(string(LoggerKey), ctx)

		logger.Info("http_request_in")

		err := c.Next()

		logger.Info("http_request_out",
			zap.Int("status", c.Response().StatusCode()),
			zap.String("latency", time.Since(start).Truncate(time.Millisecond).String()),
		)

		return err
	}
}
