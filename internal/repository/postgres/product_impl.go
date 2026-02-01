package postgres

import (
	"context"
	"errors"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/dto"
	"kasir-api/internal/entity"
	"kasir-api/internal/repository"

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

func (r *productRepo) FindByID(ctx context.Context, id uint) (entity.Product, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "ProductRepository.FindByID"),
		zap.Uint("product_id", id),
	)

	log.Info("in")

	var p entity.Product
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("out", zap.String("result", "not_found"))
			return entity.Product{}, repository.ErrNotFound
		}
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return entity.Product{}, err
	}

	log.Info("out", zap.String("result", "ok"))

	return p, nil
}

func (r *productRepo) FindAll(ctx context.Context) ([]entity.Product, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "ProductRepository.FindAll"),
	)

	log.Info("in")

	var products []entity.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return nil, err
	}

	log.Info("out", zap.Int("count", len(products)))

	return products, nil
}

func (r *productRepo) Update(ctx context.Context, p entity.Product) (entity.Product, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "ProductRepository.Update"),
		zap.Uint("product_id", p.ID),
	)

	log.Info("in")

	// ensure exists
	var current entity.Product
	if err := r.db.WithContext(ctx).First(&current, p.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("out", zap.String("result", "not_found"))
			return entity.Product{}, repository.ErrNotFound
		}
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return entity.Product{}, err
	}

	updates := map[string]interface{}{}

	// CategoryID: 0 berarti "tidak diubah" di level entity, jadi service yang ngatur.
	// Kalau repo dipanggil dengan CategoryID != 0, berarti mau update.
	if p.CategoryID != 0 && p.CategoryID != current.CategoryID {
		updates["category_id"] = p.CategoryID
	}

	if p.Name != "" && p.Name != current.Name {
		updates["name"] = p.Name
	}

	if p.Price != 0 && p.Price != current.Price {
		updates["price"] = p.Price
	}
	// stock bisa 0 valid, jadi kita perlu sentinel lain.
	// Cara gampang: service yang kirim field stock via pointer, atau gunakan dto khusus.
	// Untuk sekarang: kalau p.Stock != current.Stock (tetap akan update termasuk 0)
	if p.Stock != current.Stock {
		updates["stock"] = p.Stock
	}

	if len(updates) == 0 {
		log.Info("out", zap.String("result", "no_changes"))
		return current, nil
	}

	if err := r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Where("id = ?", p.ID).
		Updates(updates).Error; err != nil {
		log.Error("out", zap.String("result", "update_failed"), zap.Error(err))
		return entity.Product{}, err
	}

	if err := r.db.WithContext(ctx).First(&current, p.ID).Error; err != nil {
		log.Error("out", zap.String("result", "db_error_after_update"), zap.Error(err))
		return entity.Product{}, err
	}

	log.Info("out", zap.String("result", "ok"))

	return current, nil
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "ProductRepository.Delete"),
		zap.Uint("product_id", id),
	)

	log.Info("in")

	// ensure exists
	var current entity.Product
	if err := r.db.WithContext(ctx).First(&current, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("out", zap.String("result", "not_found"))
			return repository.ErrNotFound
		}
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error; err != nil {
		log.Error("out", zap.String("result", "delete_failed"), zap.Error(err))
		return err
	}

	log.Info("out", zap.String("result", "ok"))

	return nil
}

func (r *productRepo) FindDetailByID(ctx context.Context, id uint) (dto.ProductDetailResponse, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "ProductRepository.FindDetailByID"),
		zap.Uint("product_id", id),
	)

	log.Info("in")

	var out dto.ProductDetailResponse

	err := r.db.WithContext(ctx).
		Table("product p").
		Select(`
			p.id,
			p.category_id,
			c.name AS category_name,
			p.name,
			p.price,
			p.stock,
			p.created_at,
			p.updated_at
		`).
		Joins("JOIN category c ON c.id = p.category_id").
		Where("p.id = ?", id).
		Scan(&out).Error

	if err != nil {
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return dto.ProductDetailResponse{}, err
	}

	// Scan() biasanya tidak return ErrRecordNotFound, jadi check ID
	if out.ID == 0 {
		log.Info("out", zap.String("result", "not_found"))
		return dto.ProductDetailResponse{}, repository.ErrNotFound
	}

	log.Info("out", zap.String("result", "ok"), zap.String("name", out.Name))

	return out, nil
}
