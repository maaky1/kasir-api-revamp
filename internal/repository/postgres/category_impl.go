package postgres

import (
	"context"
	"errors"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/entity"
	"kasir-api/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) FindByID(ctx context.Context, id uint) (entity.Category, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "CategoryRepository.FindByID"),
		zap.Uint("category_id", id),
	)

	log.Debug("Querying category")

	var c entity.Category
	err := r.db.WithContext(ctx).Take(&c, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warn("not found")
		return entity.Category{}, repository.ErrNotFound
	}

	if err != nil {
		log.Error("Database error", zap.Error(err))
		return entity.Category{}, err
	}

	log.Debug("Category found", zap.String("name", c.Name))
	return c, err
}
