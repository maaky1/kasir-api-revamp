package config

import (
	"kasir-api/internal/infra/gormzap"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDatabase(v *viper.Viper, log *zap.Logger) *gorm.DB {
	dsn := v.GetString("database.url")
	idleConnection := v.GetInt("database.pool.idle")
	maxConnection := v.GetInt("database.pool.max")
	maxLifeTimeConnection := v.GetInt("database.pool.lifetime")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormzap.New(log, glogger.Info, 700*time.Millisecond),
	})
	if err != nil {
		log.Fatal("Failed connect to NeonDB:", zap.Error(err))
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal("Failed get sql.DB:", zap.Error(err))
	}

	sqlDb.SetMaxIdleConns(idleConnection)
	sqlDb.SetMaxOpenConns(maxConnection)
	sqlDb.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	log.Info("Success connect to NeonDB")
	return db
}
