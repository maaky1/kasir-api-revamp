package config

import (
	httpmw "kasir-api/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewFiber(v *viper.Viper, log *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      v.GetString("app.name"),
		Prefork:      v.GetBool("web.prefork"),
		ErrorHandler: NewErrorHandler(),
	})

	app.Use(httpmw.LoggingMiddleware(log))

	return app
}
