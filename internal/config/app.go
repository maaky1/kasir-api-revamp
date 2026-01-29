package config

import (
	"kasir-api/internal/delivery/http"
	"kasir-api/internal/delivery/http/routes"
	"kasir-api/internal/repository/postgres"
	"kasir-api/internal/service"

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
	categoryRepository := postgres.NewCategoryRepository(cfg.DB, cfg.Logger)
	categoryService := service.NewCategoryService(categoryRepository, cfg.Logger)
	categoryController := http.NewCategoryController(categoryService, cfg.Logger)

	routeConfig := routes.RouteConfig{
		App:                cfg.App,
		CategoryController: categoryController,
	}

	routeConfig.Setup()
}
