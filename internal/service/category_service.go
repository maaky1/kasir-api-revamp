package service

import (
	"context"
	"errors"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/dto"
	"kasir-api/internal/repository"

	"go.uber.org/zap"
)

type CategoryService interface {
	GetCategoryByID(ctx context.Context, id uint) (dto.CategoryResponse, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uint) (dto.CategoryResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "CategoryService.GetCategoryByID"),
		zap.Uint("category_id", id),
	)

	log.Debug("Calling repository")

	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("Category not found")
			return dto.CategoryResponse{}, NotFound("Category not found")
		}

		log.Error("Repository error", zap.Error(err))
		return dto.CategoryResponse{}, Internal("Internal server error")
	}

	log.Debug("Service success", zap.String("name", c.Name))
	return dto.CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}, nil
}
