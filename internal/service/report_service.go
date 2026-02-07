package service

import (
	"context"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/dto"
	"kasir-api/internal/repository"
	"time"

	"go.uber.org/zap"
)

type ReportService interface {
	GetReport(ctx context.Context, startDate string, endDate string) (dto.Report, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{reportRepo: reportRepo}
}

func (s *reportService) GetReport(ctx context.Context, startDate string, endDate string) (dto.Report, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "ReportService.GetReport"),
	)

	log.Info("in")

	now := time.Now().In(time.FixedZone("WIB", 7*3600))
	nowStr := now.Format("2006-01-02")

	effEnd := endDate
	if effEnd == "" {
		if startDate != "" {
			effEnd = nowStr // start - now
		} else {
			effEnd = nowStr // now
		}
	}

	if startDate != "" && startDate > effEnd {
		return dto.Report{}, InvalidInput("Invalid input date")
	}

	entityRes, err := s.reportRepo.GetReport(ctx, startDate, effEnd)
	if err != nil {
		log.Error("out", zap.Error(err))
		return dto.Report{}, err
	}

	bestProducts := make([]dto.BestProduct, len(entityRes.BestProduct))
	for i, v := range entityRes.BestProduct {
		bestProducts[i] = dto.BestProduct{
			Name:     v.Name,
			Quantity: v.Quantity,
			Subtotal: v.Subtotal,
		}
	}

	var rangeStr string
	if startDate == "" {
		rangeStr = nowStr
	} else if endDate == "" {
		rangeStr = startDate + " - " + nowStr
	} else {
		rangeStr = startDate + " - " + endDate
	}

	res := dto.Report{
		ReportRange:      rangeStr,
		TotalRevenue:     entityRes.TotalRevenue,
		TotalTransaction: entityRes.TotalTransaction,
		BestProduct:      bestProducts,
	}

	log.Info("out", zap.String("result", "ok"))

	return res, nil
}
