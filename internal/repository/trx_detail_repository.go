package repository

import (
	"context"
	"kasir-api/internal/entity"
)

type TrxDetailRepository interface {
	Create(ctx context.Context, trx entity.TransactionDetail) (entity.TransactionDetail, error)
}
