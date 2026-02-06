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

type TrxService interface {
	Checkout(ctx context.Context, req dto.Checkout) (dto.Transaction, error)
}

type trxService struct {
	productRepo repository.ProductRepository
	trxRepo     repository.TrxRepository
	trxDetRepo  repository.TrxDetailRepository
}

func NewTrxService(productRepo repository.ProductRepository, trxRepo repository.TrxRepository, trxDetRepo repository.TrxDetailRepository) TrxService {
	return &trxService{productRepo: productRepo, trxRepo: trxRepo, trxDetRepo: trxDetRepo}
}

func (s *trxService) Checkout(ctx context.Context, req dto.Checkout) (dto.Transaction, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "service"),
		zap.String("operation", "TrxService.Checkout"),
	)

	log.Info("in")

	if len(req.Items) <= 0 {
		log.Warn("out", zap.String("result", "invalid_items"))
		return dto.Transaction{}, InvalidInput("Items must be > 0")
	}

	var total int
	var details []dto.TransactionDetail
	// Loop item checkout
	items := req.Items
	for _, item := range items {
		if item.Quantity <= 0 {
			log.Warn("out", zap.String("result", "invalid_quantity"))
			return dto.Transaction{}, InvalidInput("Quantity must be > 0")
		}

		// Get product
		curProduct, err := s.productRepo.FindByID(ctx, item.ProductID)
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("out", zap.String("result", "not_found"))
			return dto.Transaction{}, NotFound("Product not found")
		}

		if curProduct.Stock < item.Quantity {
			return dto.Transaction{}, BadRequest("Stock not enough")
		}

		// Update product stock
		s.productRepo.Update(ctx, entity.Product{
			ID:    curProduct.ID,
			Stock: curProduct.Stock - item.Quantity,
		})

		// Calculate
		subtotal := curProduct.Price * item.Quantity
		total += subtotal

		// Save to struct
		detail := dto.TransactionDetail{
			ProductID:   curProduct.ID,
			ProductName: curProduct.Name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		}

		details = append(details, detail)
	}

	var res dto.Transaction

	return res, nil
}
