package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(v *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      v.GetString("app.name"),
		Prefork:      v.GetBool("web.prefork"),
		ErrorHandler: NewErrorHandler(),
	})

	return app
}
