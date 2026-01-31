package repository

import (
	"context"
	"kasir-api/internal/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, p entity.Product) (entity.Product, error)
	// FindByID(ctx context.Context, id uint) (entity.Product, error)
	// FindAll(ctx context.Context) ([]entity.Product, error)
	// Update(ctx context.Context, p entity.Product) (entity.Product, error)
	// Delete(ctx context.Context, id uint) error
}
