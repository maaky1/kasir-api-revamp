package routes

import (
	"kasir-api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	CategoryController *http.CategoryController
}

func (c *RouteConfig) Setup() {
	c.SetupRegister()
}

func (c *RouteConfig) SetupRegister() {
	api := c.App.Group("/api")

	categories := api.Group("/category")

	categories.Get("/:id", c.CategoryController.GetCategoryByID)

	api.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"status": "Ok"})
	})
}
