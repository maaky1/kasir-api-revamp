package service

import (
	"context"
	"errors"
	"kasir-api/internal/dto"
	"kasir-api/internal/repository"

	"go.uber.org/zap"
)

type CategoryService interface {
	GetCategoryByID(ctx context.Context, id uint) (dto.CategoryResponse, error)
}

type categoryService struct {
	repo repository.CategoryRepository
	log  *zap.Logger
}

var (
	ErrInvalidInput = errors.New("Invalid input")
	ErrConflict     = errors.New("Conflict")
	ErrNotFound     = errors.New("Category not found")
	ErrInternal     = errors.New("Internal error")
)

func NewCategoryService(repo repository.CategoryRepository, log *zap.Logger) CategoryService {
	return &categoryService{repo: repo, log: log}
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id uint) (dto.CategoryResponse, error) {
	log := s.log.With(
		zap.String("layer", "service"),
		zap.String("operation", "CategoryService.GetCategoryByID"),
		zap.Uint("category_id", id),
	)

	log.Debug("Calling repository")

	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			log.Warn(ErrNotFound.Error())
			return dto.CategoryResponse{}, ErrNotFound
		}

		log.Error("Repository error", zap.Error(err))
		return dto.CategoryResponse{}, ErrInternal
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
