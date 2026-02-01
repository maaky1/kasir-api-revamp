package http

import (
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/dto"
	"kasir-api/internal/helper"
	"kasir-api/internal/response"
	"kasir-api/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ProductController struct {
	svc service.ProductService
}

func NewProductController(svc service.ProductService) *ProductController {
	return &ProductController{svc: svc}
}

func (h *ProductController) CreateProduct(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "ProductController.CreateProduct"),
	)

	log.Info("in")

	var req dto.Product
	if err := ctx.BodyParser(&req); err != nil {
		log.Warn("out", zap.String("result", "invalid_request_body"), zap.Error(err))
		return response.Error(ctx, http.StatusBadRequest, "Invalid request body")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.CreateProduct(reqCtx, req)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusCreated, "Product created", res)
}

func (h *ProductController) GetProductByID(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "ProductController.GetProductByID"),
	)

	log.Info("in")

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("out", zap.String("result", "invalid_product_id"))
		return response.Error(ctx, http.StatusBadRequest, "Invalid product ID")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.GetProductByID(reqCtx, id)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")
	return response.Success(ctx, http.StatusOK, "Product found", res)
}

func (h *ProductController) GetAllProduct(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "ProductController.FindAllProduct"),
	)

	log.Info("in")

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.GetAllProduct(reqCtx)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out", zap.Int("count", len(res)))
	return response.Success(ctx, http.StatusOK, "Products found", res)
}

func (h *ProductController) UpdateProductByID(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "ProductController.UpdateProduct"),
	)

	log.Info("in")

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("out", zap.String("result", "invalid_product_id"))
		return response.Error(ctx, http.StatusBadRequest, "Invalid product ID")
	}

	var req dto.UpdateProduct
	if err := ctx.BodyParser(&req); err != nil {
		log.Warn("out", zap.String("result", "invalid_request_body"), zap.Error(err))
		return response.Error(ctx, http.StatusBadRequest, "Invalid request body")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.UpdateProductByID(reqCtx, id, req)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusOK, "Product updated", res)
}

func (h *ProductController) DeleteProductByID(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "ProductController.DeleteProductByID"),
	)

	log.Info("in")

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("out", zap.String("result", "invalid_product_id"))
		return response.Error(ctx, http.StatusBadRequest, "Invalid product ID")
	}

	reqCtx := middleware.RequestContext(ctx)

	if err := h.svc.DeleteProductByID(reqCtx, id); err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusOK, "Product deleted", nil)
}

func (h *ProductController) GetProductDetailByID(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "ProductController.GetProductDetailByID"),
	)

	log.Info("in")

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("out", zap.String("result", "invalid_product_id"))
		return response.Error(ctx, http.StatusBadRequest, "Invalid product ID")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.GetProductDetailByID(reqCtx, id)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusOK, "Product detail found", res)
}
