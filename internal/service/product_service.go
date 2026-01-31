package service

import (
	"context"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/dto"
	"kasir-api/internal/entity"
	"kasir-api/internal/repository"

	"go.uber.org/zap"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req dto.Product) (dto.ProductResponse, error)
	// GetProductByID(ctx context.Context, id uint) (dto.ProductResponse, error)
	// GetAllProduct(ctx context.Context) ([]dto.ProductResponse, error)
	// UpdateProductByID(ctx context.Context, id uint, req dto.Product) (dto.ProductResponse, error)
	// DeleteProductByID(ctx context.Context, id uint) error
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{productRepo: productRepo, categoryRepo: categoryRepo}
}

func (s *productService) CreateProduct(ctx context.Context, req dto.Product) (dto.ProductResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "ProductService.CreateProduct"),
	)

	log.Info("in")

	if req.Name == "" {
		log.Warn("out", zap.String("result", "name_is_required"))
		return dto.ProductResponse{}, InvalidInput("Name is required")
	}

	if req.Price <= 0 {
		log.Warn("out", zap.String("result", "invalid_price"))
		return dto.ProductResponse{}, InvalidInput("Price must be greater than 0")
	}

	if req.Stock < 0 {
		log.Warn("out", zap.String("result", "invalid_stock"))
		return dto.ProductResponse{}, InvalidInput("Stock must be >= 0")
	}

	catID := uint(1)
	if req.CategoryID != nil && *req.CategoryID != 0 {
		catID = *req.CategoryID

		if _, err := s.categoryRepo.FindByID(ctx, catID); err != nil {
			log.Warn("out", zap.String("result", "category_not_found"))
			return dto.ProductResponse{}, err
		}
	}

	p := entity.Product{
		CategoryID: catID,
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
	}

	created, err := s.productRepo.Create(ctx, p)
	if err != nil {
		log.Error("out", zap.Error(err))
		return dto.ProductResponse{}, err
	}

	res := dto.ProductResponse{
		ID:         created.ID,
		CategoryID: created.CategoryID,
		Name:       created.Name,
		Price:      created.Price,
		Stock:      created.Stock,
		CreatedAt:  created.CreatedAt,
		UpdatedAt:  created.UpdatedAt,
	}

	log.Info("out", zap.String("result", "ok"), zap.Uint("product_id", created.ID))

	return res, nil
}
