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
	categoryRepository := postgres.NewCategoryRepository(cfg.DB)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryController := http.NewCategoryController(categoryService)

	productRepository := postgres.NewProductRepository(cfg.DB)
	productService := service.NewProductService(productRepository, categoryRepository)
	productController := http.NewProductController(productService)

	routeConfig := routes.RouteConfig{
		App:                cfg.App,
		CategoryController: categoryController,
		ProductController:  productController,
	}

	routeConfig.Setup()
}
