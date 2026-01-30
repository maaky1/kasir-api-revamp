package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// dipakai di CONTROLLER
func LoggerFromFiber(c *fiber.Ctx) *zap.Logger {
	return LoggerFromCtx(c.UserContext())
}

// dipakai buat pass ke SERVICE / REPO
func RequestContext(c *fiber.Ctx) context.Context {
	return c.UserContext()
}

// dipakai di SERVICE / REPO
func LoggerFromCtx(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.NewNop()
	}
	if v := ctx.Value(LoggerKey); v != nil {
		if log, ok := v.(*zap.Logger); ok && log != nil {
			return log
		}
	}
	return zap.NewNop()
}
