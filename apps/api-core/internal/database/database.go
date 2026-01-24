package database

import (
	"fmt"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/config"
	"gitlab.bellsoft.net/rms/api-core/internal/migrations"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
	)

	gormConfig := &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	if cfg.Environment == "local" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Hooks disabled - manually handle soft delete

	return db, nil
}

func Migrate(db *gorm.DB) error {
	// Create migration manager
	manager := migrations.NewMigrationManager(db)

	// Add all migrations
	for _, migration := range migrations.AllMigrations() {
		manager.AddMigration(migration)
	}

	// Run migrations
	return manager.Migrate()
}
