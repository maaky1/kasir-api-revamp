package routes

import (
	"kasir-api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	CategoryController *http.CategoryController
	ProductController  *http.ProductController
	TrxController      *http.TrxController
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
	product.Get("/:id", c.ProductController.GetProductByID)
	product.Get("", c.ProductController.GetAllProduct)
	product.Put("/:id", c.ProductController.UpdateProductByID)
	product.Delete("/:id", c.ProductController.DeleteProductByID)
	product.Get("/:id/detail", c.ProductController.GetProductDetailByID)

	trx := api.Group("/transaction")
	trx.Post("/checkout", c.TrxController.Checkout)

	api.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"status": "Ok"})
	})
}
