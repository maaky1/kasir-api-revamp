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

type CategoryController struct {
	svc service.CategoryService
}

func NewCategoryController(svc service.CategoryService) *CategoryController {
	return &CategoryController{svc: svc}
}

func (h *CategoryController) CreateCategory(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "CategoryController.CreateCategory"),
	)

	log.Info("in")

	var req dto.Category
	if err := ctx.BodyParser(&req); err != nil {
		log.Warn("out", zap.String("result", "invalid_request_body"), zap.Error(err))
		return response.Error(ctx, http.StatusBadRequest, "Invalid request body")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.CreateCategory(reqCtx, req)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusCreated, "Category created", res)
}

func (h *CategoryController) GetCategoryByID(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "CategoryController.GetCategoryByID"),
	)

	log.Info("in")

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("out", zap.String("result", "invalid_category_id"))
		return response.Error(ctx, http.StatusBadRequest, "Invalid category ID")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.GetCategoryByID(reqCtx, id)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusOK, "Category found", res)
}

func (h *CategoryController) GetAllCategory(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "CategoryController.GetAllCategory"),
	)

	log.Info("in")

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.GetAllCategory(reqCtx)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out", zap.Int("count", len(res)))

	return response.Success(ctx, http.StatusOK, "Categories list", res)
}

func (h *CategoryController) UpdateCategoryByID(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "CategoryController.UpdateCategoryByID"),
	)

	log.Info("in")

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("out", zap.String("result", "invalid_category_id"))
		return response.Error(ctx, http.StatusBadRequest, "Invalid category ID")
	}

	var req dto.Category
	if err := ctx.BodyParser(&req); err != nil {
		log.Warn("out", zap.String("result", "invalid_request_body"), zap.Error(err))
		return response.Error(ctx, http.StatusBadRequest, "Invalid request body")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.UpdateCategoryByID(reqCtx, id, req)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusOK, "Category updated", res)
}

func (h *CategoryController) DeleteCategoryByID(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "CategoryController.DeleteCategoryByID"),
	)

	log.Info("in")

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("out", zap.String("result", "invalid_category_id"))
		return response.Error(ctx, http.StatusBadRequest, "Invalid category ID")
	}

	reqCtx := middleware.RequestContext(ctx)

	if err := h.svc.DeleteCategoryByID(reqCtx, id); err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusOK, "Category deleted", nil)
}
