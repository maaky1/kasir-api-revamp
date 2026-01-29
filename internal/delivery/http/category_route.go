package http

import (
	"kasir-api/internal/helper"
	"kasir-api/internal/response"
	"kasir-api/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CategoryController struct {
	svc service.CategoryService
	log *zap.Logger
}

func NewCategoryController(svc service.CategoryService, log *zap.Logger) *CategoryController {
	return &CategoryController{svc: svc, log: log}
}

func (h *CategoryController) GetCategoryByID(ctx *fiber.Ctx) error {
	log := h.log.With(
		zap.String("layer", "controller"),
		zap.String("operation", "CategoryController.GetCategoryByID"),
	)

	log.Debug("Incoming request", zap.String("id", ctx.Params("id")))

	id, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		log.Warn("invalid id", zap.String("id", ctx.Params("id")), zap.Error(err))
		return response.Error(ctx, http.StatusBadRequest, "Invalid id")
	}

	res, err := h.svc.GetCategoryByID(ctx.Context(), id)
	if err != nil {
		log.Warn("Service error", zap.Error(err))
		return writeServiceError(ctx, err)
	}

	return response.Success(ctx, http.StatusOK, "Category found", res)
}

func writeServiceError(ctx *fiber.Ctx, err error) error {
	switch err {
	case service.ErrInvalidInput:
		return response.Error(ctx, http.StatusBadRequest, "Invalid input")
	case service.ErrNotFound:
		return response.Error(ctx, http.StatusNotFound, "Category not found")
	case service.ErrConflict:
		return response.Error(ctx, http.StatusConflict, "Conflict")
	default:
		return response.Error(ctx, http.StatusInternalServerError, "Internal error")
	}
}
