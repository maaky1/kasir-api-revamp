package postgres

import (
	"context"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type trxRepo struct {
	db *gorm.DB
}

func NewTrxRepository(db *gorm.DB) *trxRepo {
	return &trxRepo{db: db}
}

func (r *trxRepo) Create(ctx context.Context, trx entity.Transaction) (entity.Transaction, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "TrxRepository.Create"),
		zap.Uint("id", trx.ID),
	)

	log.Info("in")

	if err := r.db.WithContext(ctx).Create(&trx).Error; err != nil {
		log.Error("out", zap.String("result", "insert_failed"), zap.Error(err))
		return entity.Transaction{}, err
	}

	log.Info("out", zap.String("result", "ok"), zap.Uint("transaction_id", trx.ID))

	return trx, nil
}
