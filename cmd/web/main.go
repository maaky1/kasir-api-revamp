package main

import (
	"fmt"
	"kasir-api/internal/config"
	"kasir-api/internal/database"

	"go.uber.org/zap"
)

func main() {
	// Load config
	v := config.NewViper()

	// Init logger
	logger := config.NewLogger(v)
	defer func() { _ = logger.Sync() }()

	// Get dsn
	dsn := v.GetString("database.url")

	// Run migrations
	database.RunMigration(dsn, logger)

	// Init Gorm
	db := config.NewDatabase(v, logger)

	// Init Fiber
	app := config.NewFiber(v)

	// Register routes / wiring
	config.Bootstrap(&config.BootstrapConfig{
		Config: v,
		Logger: logger,
		DB:     db,
		App:    app,
	})

	// Start server
	port := v.GetInt("web.port")
	logger.Info("Server starting", zap.Int("port", port))

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
