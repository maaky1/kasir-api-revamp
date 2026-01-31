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
