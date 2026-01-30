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

func (r *categoryRepo) Create(ctx context.Context, c entity.Category) (entity.Category, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "CategoryRepository.Create"),
		zap.String("name", c.Name),
	)

	log.Info("in")

	var exist entity.Category
	err := r.db.WithContext(ctx).
		Where("name = ?", c.Name).
		First(&exist).Error

	if err == nil {
		log.Info("out", zap.String("result", "conflict"))
		return entity.Category{}, repository.ErrConflict
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return entity.Category{}, err
	}

	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		log.Error("out", zap.String("result", "insert_failed"), zap.Error(err))
		return entity.Category{}, err
	}

	log.Info("out", zap.String("result", "ok"), zap.Uint("category_id", c.ID))

	return c, nil
}

func (r *categoryRepo) FindByID(ctx context.Context, id uint) (entity.Category, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "CategoryRepository.FindByID"),
		zap.Uint("category_id", id),
	)

	log.Info("in")

	var c entity.Category
	err := r.db.WithContext(ctx).Take(&c, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Info("out", zap.String("result", "not_found"))
		return entity.Category{}, repository.ErrNotFound
	}

	if err != nil {
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return entity.Category{}, err
	}

	log.Info("out", zap.String("result", "ok"), zap.String("name", c.Name))

	return c, err
}

func (r *categoryRepo) FindAll(ctx context.Context) ([]entity.Category, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "CategoryRepository.FindAll"),
	)

	log.Info("in")

	var categories []entity.Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return nil, err
	}

	log.Info("out", zap.Int("count", len(categories)))
	return categories, nil
}

func (r *categoryRepo) Update(ctx context.Context, c entity.Category) (entity.Category, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "CategoryRepository.Update"),
		zap.Uint("category_id", c.ID),
	)

	log.Info("in")

	var current entity.Category
	if err := r.db.WithContext(ctx).First(&current, c.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("out", zap.String("result", "not_found"))
			return entity.Category{}, repository.ErrNotFound
		}
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return entity.Category{}, err
	}

	// build partial updates (lebih aman & jelas)
	updates := map[string]interface{}{}

	if c.Name != "" && c.Name != current.Name {
		updates["name"] = c.Name
	}
	if c.Description != "" && c.Description != current.Description {
		updates["description"] = c.Description
	}

	// nothing changed -> return current (updated_at juga gak perlu berubah)
	if len(updates) == 0 {
		log.Info("out", zap.String("result", "no_changes"))
		return current, nil
	}

	// update (autoUpdateTime akan set updated_at otomatis)
	if err := r.db.WithContext(ctx).
		Model(&entity.Category{}).
		Where("id = ?", c.ID).
		Updates(updates).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Info("out", zap.String("result", "conflict"))
			return entity.Category{}, repository.ErrConflict
		}

		log.Error("out", zap.String("result", "update_failed"), zap.Error(err))
		return entity.Category{}, err
	}

	if err := r.db.WithContext(ctx).First(&current, c.ID).Error; err != nil {
		log.Error("out", zap.String("result", "db_error_after_update"), zap.Error(err))
		return entity.Category{}, err
	}

	log.Info("out", zap.String("result", "ok"), zap.Uint("category_id", current.ID))

	return current, nil
}
