package routes

import (
	"kasir-api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	CategoryController *http.CategoryController
	ProductController  *http.ProductController
}

func (c *RouteConfig) Setup() {
	c.SetupRegister()
}

func (c *RouteConfig) SetupRegister() {
	api := c.App.Group("/api")

	category := api.Group("/category")
	category.Post("", c.CategoryController.CreateCategory)
	category.Get("/:id", c.CategoryController.GetCategoryByID)
	category.Get("", c.CategoryController.GetAllCategory)
	category.Put("/:id", c.CategoryController.UpdateCategoryByID)
	category.Delete("/:id", c.CategoryController.DeleteCategoryByID)

	product := api.Group("/product")
	product.Post("", c.ProductController.CreateProduct)

	api.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"status": "Ok"})
	})
}
