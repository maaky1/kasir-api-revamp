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

type ProductService interface {
	CreateProduct(ctx context.Context, req dto.Product) (dto.ProductResponse, error)
	GetProductByID(ctx context.Context, id uint) (dto.ProductResponse, error)
	GetAllProduct(ctx context.Context) ([]dto.ProductResponse, error)
	UpdateProductByID(ctx context.Context, id uint, req dto.UpdateProduct) (dto.ProductResponse, error)
	DeleteProductByID(ctx context.Context, id uint) error
	GetProductDetailByID(ctx context.Context, id uint) (dto.ProductDetailResponse, error)
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

func (s *productService) GetProductByID(ctx context.Context, id uint) (dto.ProductResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "ProductService.GetProductByID"),
	)

	log.Info("in")

	p, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("out", zap.String("result", "not_found"))
			return dto.ProductResponse{}, NotFound("Product not found")
		}
		log.Error("out", zap.Error(err))
		return dto.ProductResponse{}, err
	}

	res := dto.ProductResponse{
		ID:         p.ID,
		CategoryID: p.CategoryID,
		Name:       p.Name,
		Price:      p.Price,
		Stock:      p.Stock,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}

	log.Info("out", zap.String("result", "ok"))

	return res, nil
}

func (s *productService) GetAllProduct(ctx context.Context) ([]dto.ProductResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "ProductService.FindAllProduct"),
	)

	log.Info("in")

	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		log.Error("out", zap.Error(err))
		return nil, err
	}

	res := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		res = append(res, dto.ProductResponse{
			ID:         p.ID,
			CategoryID: p.CategoryID,
			Name:       p.Name,
			Price:      p.Price,
			Stock:      p.Stock,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
		})
	}

	log.Info("out", zap.Int("count", len(res)))

	return res, nil
}

func (s *productService) UpdateProductByID(ctx context.Context, id uint, req dto.UpdateProduct) (dto.ProductResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "ProductService.UpdateProduct"),
	)

	log.Info("in")

	// minimal 1 field
	if req.CategoryID == nil && req.Name == nil && req.Price == nil && req.Stock == nil {
		log.Warn("out", zap.String("result", "no_fields_to_update"))
		return dto.ProductResponse{}, InvalidInput("Nothing to update")
	}

	// prepare entity for repo
	update := entity.Product{ID: id}

	// category optional
	if req.CategoryID != nil {
		if *req.CategoryID == 0 {
			log.Warn("out", zap.String("result", "invalid_category_id"))
			return dto.ProductResponse{}, InvalidInput("Invalid category ID")
		}

		// optional: validate category exists (kalau bukan default 1 pun boleh dicek)
		if _, err := s.categoryRepo.FindByID(ctx, *req.CategoryID); err != nil {
			return dto.ProductResponse{}, err
		}

		update.CategoryID = *req.CategoryID
	}

	if req.Name != nil {
		if *req.Name == "" {
			return dto.ProductResponse{}, InvalidInput("Name cannot be empty")
		}
		update.Name = *req.Name
	}

	if req.Price != nil {
		if *req.Price <= 0 {
			return dto.ProductResponse{}, InvalidInput("Price must be greater than 0")
		}
		update.Price = *req.Price
	}

	if req.Stock != nil {
		if *req.Stock < 0 {
			return dto.ProductResponse{}, InvalidInput("Stock must be >= 0")
		}
		update.Stock = *req.Stock
	}

	updated, err := s.productRepo.Update(ctx, update)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("out", zap.String("result", "not_found"))
			return dto.ProductResponse{}, NotFound("Product not found")
		}
		log.Error("out", zap.Error(err))
		return dto.ProductResponse{}, err
	}

	res := dto.ProductResponse{
		ID:         updated.ID,
		CategoryID: updated.CategoryID,
		Name:       updated.Name,
		Price:      updated.Price,
		Stock:      updated.Stock,
		CreatedAt:  updated.CreatedAt,
		UpdatedAt:  updated.UpdatedAt,
	}

	log.Info("out", zap.String("result", "ok"))

	return res, nil
}

func (s *productService) DeleteProductByID(ctx context.Context, id uint) error {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "ProductService.DeleteProductByID"),
	)

	log.Info("in")

	if id == 0 {
		log.Warn("out", zap.String("result", "invalid_product_id"))
		return InvalidInput("Invalid product ID")
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("out", zap.String("result", "not_found"))
			return NotFound("Product not found")
		}

		log.Error("out", zap.Error(err))
		return err
	}

	log.Info("out", zap.String("result", "ok"))

	return nil
}

func (s *productService) GetProductDetailByID(ctx context.Context, id uint) (dto.ProductDetailResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "ProductService.GetProductDetailByID"),
	)

	log.Info("in")

	res, err := s.productRepo.FindDetailByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("out", zap.String("result", "not_found"))
			return dto.ProductDetailResponse{}, NotFound("Product not found")
		}

		log.Error("out", zap.Error(err))
		return dto.ProductDetailResponse{}, err
	}

	log.Info("out", zap.String("result", "ok"))

	return res, nil
}
