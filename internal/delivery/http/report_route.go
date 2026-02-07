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

type ReportController struct {
	svc service.ReportService
}

func NewReportController(svc service.ReportService) *ReportController {
	return &ReportController{svc: svc}
}

func (h *ReportController) GetReport(ctx *fiber.Ctx) error {
	log := middleware.LoggerFromFiber(ctx).With(
		zap.String("layer", "controller"),
		zap.String("operation", "ReportController.GetReport"),
	)

	log.Info("in")

	startDate, err := helper.ParseDateQuery(ctx, "startDate")
	if err != nil {
		log.Warn("out", zap.String("result", "invalid_start_date"))
		return response.Error(ctx, http.StatusBadRequest, "Invalid input start date")
	}

	endDate, err := helper.ParseDateQuery(ctx, "endDate")
	if err != nil {
		log.Warn("out", zap.String("result", "invalid_end_date"))
		return response.Error(ctx, http.StatusBadRequest, "Invalid input end date")
	}

	reqCtx := middleware.RequestContext(ctx)

	res, err := h.svc.GetReport(reqCtx, startDate, endDate)
	if err != nil {
		log.Error("out", zap.Error(err))
		return helper.WriteServiceError(ctx, err)
	}

	log.Info("out")

	return response.Success(ctx, http.StatusOK, "Get report successfully", res)
}
