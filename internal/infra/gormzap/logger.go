package gormzap

import (
	"context"
	"errors"
	"time"

	httpmw "kasir-api/internal/delivery/http/middleware"

	"go.uber.org/zap"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type Logger struct {
	base          *zap.Logger
	level         glogger.LogLevel
	slowThreshold time.Duration
}

func New(base *zap.Logger, level glogger.LogLevel, slowThreshold time.Duration) *Logger {
	return &Logger{
		base:          base,
		level:         level,
		slowThreshold: slowThreshold,
	}
}

func (l *Logger) LogMode(level glogger.LogLevel) glogger.Interface {
	n := *l
	n.level = level
	return &n
}

// gak dipakai, SQL keluar dari Trace()
func (l *Logger) Info(ctx context.Context, msg string, data ...interface{})  {}
func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{})  {}
func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level == glogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	log := httpmw.LoggerFromCtx(ctx)
	if log == zap.L() {
		log = l.base
	}

	fields := []zap.Field{
		zap.String("layer", "db"),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Int64("latency_ms", elapsed.Milliseconds()),
	}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.level >= glogger.Error:
		log.Error("db_query", append(fields, zap.Error(err))...)
	case l.slowThreshold > 0 && elapsed > l.slowThreshold && l.level >= glogger.Warn:
		log.Warn("db_query_slow", fields...)
	default:
		if l.level >= glogger.Info {
			log.Info("db_query", fields...)
		}
	}
}
