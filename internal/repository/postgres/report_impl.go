package postgres

import (
	"context"
	"kasir-api/internal/delivery/http/middleware"
	"kasir-api/internal/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type reportRepo struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *reportRepo {
	return &reportRepo{db: db}
}

func (r *reportRepo) GetReport(ctx context.Context, startDate string, endDate string) (entity.Report, error) {
	log := middleware.LoggerFromCtx(ctx).With(
		zap.String("layer", "repository"),
		zap.String("operation", "ReportRepository.GetReport"),
	)

	log.Info("in")

	sd := nullIfEmpty(startDate)
	ed := nullIfEmpty(endDate)

	var result entity.Report
	if err := r.db.WithContext(ctx).
		Table("transaction").
		Select(`
			COALESCE(SUM(total_amount), 0) AS total_revenue, 
			COUNT(id) AS total_transaction
		`).
		Where(`
			created_at >= COALESCE(?::date, CURRENT_DATE) AND 
			created_at < COALESCE(?::date, CURRENT_DATE) + INTERVAL '1 day'
		`, sd, ed).
		Scan(&result).
		Error; err != nil {
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return result, err
	}

	var bestProduct []entity.BestProduct
	if err := r.db.WithContext(ctx).
		Table("transaction_detail td").
		Select(`
			p.name AS name, 
			SUM(td.quantity) AS quantity, 
			SUM(td.quantity * p.price) AS subtotal
		`).
		Joins(`
			JOIN product p ON td.product_id = p.id 
			JOIN "transaction" t ON td.transaction_id = t.id 
		`).
		Where(`
			t.created_at >= COALESCE(?::date, CURRENT_DATE) AND 
			t.created_at < COALESCE(?::date, CURRENT_DATE) + INTERVAL '1 day'
		`, sd, ed).
		Group("p.id, p.name").
		Order("quantity DESC").
		Scan(&bestProduct).
		Error; err != nil {
		log.Error("out", zap.String("result", "db_error"), zap.Error(err))
		return result, err
	}

	result.BestProduct = bestProduct

	log.Info("out", zap.String("result", "ok"))

	return result, nil
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
