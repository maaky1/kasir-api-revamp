package postgres

import (
	"context"
	"errors"
	"kasir-api/internal/entity"
	"kasir-api/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type categoryRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCategoryRepository(db *gorm.DB, log *zap.Logger) *categoryRepo {
	return &categoryRepo{db: db, log: log}
}

func (r *categoryRepo) FindByID(ctx context.Context, id uint) (entity.Category, error) {
	log := r.log.With(
		zap.String("layer", "repository"),
		zap.String("operation", "CategoryRepository.FindByID"),
		zap.Uint("category_id", id),
	)

	log.Debug("Querying category")

	var c entity.Category
	err := r.db.WithContext(ctx).First(&c, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warn(repository.ErrCategoryNotFound.Error())
		return entity.Category{}, repository.ErrCategoryNotFound
	}

	if err != nil {
		log.Error("Database error", zap.Error(err))
		return entity.Category{}, err
	}

	log.Debug("Category found", zap.String("name", c.Name))
	return c, err
}
