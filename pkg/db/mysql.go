package db

import (
	"fmt"
	"multitenant/pkg/config"
	"time"

	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewMySQL(cfg *config.Config) (*gorm.DB, error) {
	// Set a GORM logger with better logging levels
	gormLogger := logger.Default.LogMode(logger.Warn)

	var db *gorm.DB
	var err error

	// Retry connection with exponential backoff
	err = retryWithBackoff(3, 2*time.Second, func() error {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   gormLogger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
				NoLowerCase:   false,
			},
		})
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL after retries: %w", err)
	}

	return db, err
}
