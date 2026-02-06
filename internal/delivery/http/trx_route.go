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

type TrxController struct {
	svc service.TrxService
}

func NewTrxController(svc service.TrxService) *TrxController {
	return &TrxController{svc: svc}
}

func (h *TrxController) Checkout(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "TrxController.Checkout"),
	)

	log.Info("in")

	var req dto.Checkout
	if err := ctx.BodyParser(&req); err != nil {
		log.Warn("out", zap.String("result", "invalid_request_body"), zap.Error(err))
		return response.Error(ctx, http.StatusBadRequest, "Invalid request body")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.Checkout(reqCtx, req)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusCreated, "Checkout successfully", res)
}
