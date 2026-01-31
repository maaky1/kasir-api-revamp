package postgres

import (
	"context"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepo {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, p entity.Product) (entity.Product, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "ProductRepository.Create"),
		zap.String("name", p.Name),
	)

	log.Info("in")

	if p.CategoryID == 0 {
		p.CategoryID = 1
	}

	if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
		log.Error("out", zap.String("result", "insert_failed"), zap.Error(err))
		return entity.Product{}, err
	}

	log.Info("out", zap.String("result", "ok"), zap.Uint("product_id", p.ID), zap.Uint("category_id", p.CategoryID))

	return p, nil
}
