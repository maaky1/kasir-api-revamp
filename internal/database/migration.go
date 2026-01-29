package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func RunMigration(dsn string, log *zap.Logger) {
	m, err := migrate.New(
		"file://internal/database/migrations",
		dsn,
	)

	if err != nil {
		log.Fatal("Migration init error", zap.Error(err))
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		log.Info("Migration skipped (no changes)")
		return
	}

	if err != nil {
		log.Fatal("Migration run error", zap.Error(err))
	}

	log.Info("Migration applied successfully")
}
