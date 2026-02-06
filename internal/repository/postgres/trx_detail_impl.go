package postgres

import (
	"context"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type trxDetailRepository struct {
	db *gorm.DB
}

func NewTrxDetailRepository(db *gorm.DB) *trxDetailRepository {
	return &trxDetailRepository{db: db}
}

func (r *trxDetailRepository) Create(ctx context.Context, trx entity.TransactionDetail) (entity.TransactionDetail, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "TrxDetailRepository.Create"),
		zap.Uint("id", trx.ID),
	)

	log.Info("in")

	if err := r.db.WithContext(ctx).Create(&trx).Error; err != nil {
		log.Error("out", zap.String("result", "insert_failed"), zap.Error(err))
		return entity.TransactionDetail{}, err
	}

	log.Info("out", zap.String("result", "ok"), zap.Uint("transaction_detail_id", trx.ID))

	return trx, nil
}
