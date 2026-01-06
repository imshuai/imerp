package config

import (
	"fmt"
	"erp/models"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	var err error

	// 连接SQLite数据库
	dsn := "database/erp.db"
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		// 只记录错误级别的SQL日志，避免泄露敏感数据
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		Error("Failed to connect database",
			zap.Error(err),
			zap.String("dsn", dsn))
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移数据表
	err = DB.AutoMigrate(
		&models.Person{},
		&models.Customer{},
		&models.Task{},
		&models.Agreement{},
		&models.Payment{},
		&models.AdminUser{},
		&models.AuditLog{},
	)
	if err != nil {
		Error("Failed to migrate database", zap.Error(err))
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	Info("Database initialized successfully",
		zap.String("path", dsn),
		zap.String("driver", "sqlite (pure Go)"))
	return nil
}
