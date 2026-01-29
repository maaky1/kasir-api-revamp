package repository

import (
	"context"
	"kasir-api/internal/entity"
)

type CategoryRepository interface {
	FindByID(ctx context.Context, id uint) (entity.Category, error)
}
