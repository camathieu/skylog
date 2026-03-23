package models

import (
	"fmt"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB opens (or creates) the SQLite database in dataDir and auto-migrates all models.
func InitDB(dataDir string) (*gorm.DB, error) {
	dbPath := filepath.Join(dataDir, "skylog.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// SQLite performance tuning
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	if _, err := sqlDB.Exec("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;"); err != nil {
		return nil, fmt.Errorf("set pragmas: %w", err)
	}

	// Auto-migrate schema
	if err := db.AutoMigrate(&Jump{}); err != nil {
		return nil, fmt.Errorf("migrate schema: %w", err)
	}

	return db, nil
}
