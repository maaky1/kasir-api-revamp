package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// dipakai di CONTROLLER
func LoggerFromFiber(c *fiber.Ctx) *zap.Logger {
	if ctx, ok := c.Locals(string(LoggerKey)).(context.Context); ok {
		if log, ok := ctx.Value(LoggerKey).(*zap.Logger); ok {
			return log
		}
	}
	return zap.L()
}

// dipakai buat pass ke SERVICE / REPO
func RequestContext(c *fiber.Ctx) context.Context {
	if ctx, ok := c.Locals(string(LoggerKey)).(context.Context); ok {
		return ctx
	}
	return context.Background()
}

// dipakai di SERVICE / REPO
func LoggerFromCtx(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(LoggerKey).(*zap.Logger); ok {
		return log
	}
	return zap.L()
}
