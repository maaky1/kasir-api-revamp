package service

import (
	"context"
	"errors"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/dto"
	"kasir-api/internal/entity"
	"kasir-api/internal/repository"

	"go.uber.org/zap"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req dto.Category) (dto.CategoryResponse, error)
	GetCategoryByID(ctx context.Context, id uint) (dto.CategoryResponse, error)
	GetAllCategory(ctx context.Context) ([]dto.CategoryResponse, error)
	UpdateCategoryByID(ctx context.Context, id uint, req dto.Category) (dto.CategoryResponse, error)
	DeleteCategoryByID(ctx context.Context, id uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) CreateCategory(ctx context.Context, req dto.Category) (dto.CategoryResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "CategoryService.CreateCategory"),
	)

	log.Info("in")

	if req.Name == "" {
		log.Warn("out", zap.String("result", "name_is_required"))
		return dto.CategoryResponse{}, InvalidInput("Name is required")
	}

	if req.Description == "" {
		log.Warn("out", zap.String("result", "description_is_required"))
		return dto.CategoryResponse{}, InvalidInput("Description is required")
	}

	entity := entity.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	created, err := s.categoryRepo.Create(ctx, entity)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			log.Warn("out", zap.String("result", "conflict"))
			return dto.CategoryResponse{}, Conflict("Category already exists")
		}

		log.Error("out", zap.String("result", "repository_error"), zap.Error(err))
		return dto.CategoryResponse{}, err
	}

	res := dto.CategoryResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}

	log.Info("out", zap.String("result", "ok"))
	log.Debug("category_created", zap.String("name", res.Name))

	return res, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uint) (dto.CategoryResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "CategoryService.GetCategoryByID"),
	)

	log.Info("in")

	c, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("out", zap.String("result", "not_found"))
			return dto.CategoryResponse{}, NotFound("Category not found")
		}

		log.Error("out", zap.String("result", "repository_error"), zap.Error(err))
		return dto.CategoryResponse{}, Internal("Internal server error")
	}

	log.Info("out", zap.String("result", "ok"))
	log.Debug("category_loaded", zap.String("name", c.Name))

	return dto.CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}, nil
}

func (s *categoryService) GetAllCategory(ctx context.Context) ([]dto.CategoryResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "CategoryService.GetAllCategory"),
	)

	log.Info("in")

	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		log.Error("out", zap.Error(err))
		return nil, err
	}

	res := make([]dto.CategoryResponse, 0, len(categories))
	for _, c := range categories {
		res = append(res, dto.CategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}

	log.Info("out", zap.Int("count", len(res)))
	return res, nil
}

func (s *categoryService) UpdateCategoryByID(ctx context.Context, id uint, req dto.Category) (dto.CategoryResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "CategoryService.UpdateCategoryByID"),
	)

	log.Info("in")

	if req.Name == "" && req.Description == "" {
		log.Warn("out", zap.String("result", "no_fields_to_update"))
		return dto.CategoryResponse{}, InvalidInput("Nothing to update")
	}

	cat := entity.Category{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	updated, err := s.categoryRepo.Update(ctx, cat)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("out", zap.String("result", "not_found"))
			return dto.CategoryResponse{}, NotFound("Category not found")
		}
		if errors.Is(err, repository.ErrConflict) {
			log.Warn("out", zap.String("result", "conflict"))
			return dto.CategoryResponse{}, Conflict("Category name already exists")
		}

		log.Error("out", zap.Error(err))
		return dto.CategoryResponse{}, err
	}

	res := dto.CategoryResponse{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}

	log.Info("out", zap.String("result", "ok"), zap.Uint("category_id", updated.ID))

	return res, nil
}

func (s *categoryService) DeleteCategoryByID(ctx context.Context, id uint) error {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "CategoryService.DeleteCategoryByID"),
	)

	log.Info("in")

	if id == 0 {
		log.Warn("out", zap.String("result", "invalid_category_id"))
		return InvalidInput("Invalid category ID")
	}

	err := s.categoryRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("out", zap.String("result", "not_found"))
			return NotFound("Category not found")
		}

		if errors.Is(err, repository.ErrForbidden) {
			log.Warn("out", zap.String("result", "forbidden"))
			return Forbidden("Default category cannot be deleted")
		}

		log.Error("out", zap.Error(err))
		return err
	}

	log.Info("out", zap.String("result", "ok"))

	return nil
}
