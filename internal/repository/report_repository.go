package repository

import (
	"context"
	"kasir-api/internal/entity"
)

type ReportRepository interface {
	GetReport(ctx context.Context, startDate string, endDate string) (entity.Report, error)
}
