package repository

import (
	"context"
	"kasir-api/internal/entity"
)

type TrxRepository interface {
	Create(ctx context.Context, trx entity.Transaction) (entity.Transaction, error)
}
