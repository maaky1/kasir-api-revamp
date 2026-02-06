package service

import (
	"context"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/dto"
	"kasir-api/internal/repository"

	"go.uber.org/zap"
)

type TrxService interface {
	Checkout(ctx context.Context, req dto.Checkout) (dto.Transaction, error)
}

type trxService struct {
	trxRepo    repository.TrxRepository
	trxDetRepo repository.TrxDetailRepository
}

func NewTrxService(trxRepo repository.TrxRepository, trxDetRepo repository.TrxDetailRepository) TrxService {
	return &trxService{trxRepo: trxRepo, trxDetRepo: trxDetRepo}
}

func (s *trxService) Checkout(ctx context.Context, req dto.Checkout) (dto.Transaction, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "TrxService.CreateTrx"),
	)

	log.Info("in")

	log.Info("Cek", zap.Any("Hasil:", len(req.Items)))
	if len(req.Items) <= 0 {
		log.Warn("out", zap.String("result", "invalid_items"))
		return dto.Transaction{}, InvalidInput("Items must be > 0")
	}

	return dto.Transaction{}, nil
}
