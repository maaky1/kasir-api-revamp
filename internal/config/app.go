package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Config *viper.Viper
	Logger *zap.Logger
	DB     *gorm.DB
	App    *fiber.App
}

func Bootstrap(cfg *BootstrapConfig) {
	cfg.Logger.Info("Settingnya disini semua")

	// Minimal route biar bisa test server hidup
	// cfg.App.Get("/health", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{
	// 		"status": "ok",
	// 	})
	// })
}
