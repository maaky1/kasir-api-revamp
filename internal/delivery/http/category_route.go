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

	log.Info("Request received", zap.String("id", ctx.Params("id")))

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("Invalid category id", zap.String("id", ctx.Params("id")), zap.Error(err))
		return response.Error(ctx, http.StatusBadRequest, "Invalid category id")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.GetCategoryByID(reqCtx, id)
	if err != nil {
		log.Warn("Service error", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	return response.Success(ctx, http.StatusOK, "Category found", res)
}
