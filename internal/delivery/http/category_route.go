package http

import (
	"kasir-api/internal/delivery/http/middleware"
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

func (h *CategoryController) GetCategoryByID(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "CategoryController.GetCategoryByID"),
	)

	log.Info("handle_get_category_by_id", zap.String("id_param", ctx.Params("id")))

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("invalid_category_id", zap.String("id_param", ctx.Params("id")), zap.Error(err))
		return response.Error(ctx, http.StatusBadRequest, "Invalid category ID")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.GetCategoryByID(reqCtx, id)
	if err != nil {
		log.Error("get_category_by_id_failed", zap.Uint("category_id", id), zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("get_category_by_id_succeeded", zap.Uint("category_id", id))
	return response.Success(ctx, http.StatusOK, "Category found", res)
}
