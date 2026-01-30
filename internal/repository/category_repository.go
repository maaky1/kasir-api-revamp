package repository

import (
	"context"
	"kasir-api/internal/entity"
)

type CategoryRepository interface {
	Create(ctx context.Context, c entity.Category) (entity.Category, error)
	FindByID(ctx context.Context, id uint) (entity.Category, error)
	FindAll(ctx context.Context) ([]entity.Category, error)
	Update(ctx context.Context, c entity.Category) (entity.Category, error)
}
