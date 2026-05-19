package config

import (
	"fmt"

	"backend-travelku/pkg/logger"

	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// InitDatabase menginisialisasi koneksi database
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})

	if err != nil {
		logger.Error(
			"Failed to connect to database",
			zap.String("db_host", cfg.DBHost),
			zap.String("db_port", cfg.DBPort),
			zap.String("db_name", cfg.DBName),
			zap.String("db_user", cfg.DBUser),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get database instance")
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	logger.Info("Database connected successfully")
	return db, nil
}
