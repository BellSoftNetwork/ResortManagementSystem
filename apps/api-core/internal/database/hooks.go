package database

import (
	"time"

	"gorm.io/gorm"
)

var DefaultDeletedAt = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

func RegisterHooks(db *gorm.DB) {
	// Register global scopes for soft delete
	db.Callback().Query().Before("gorm:query").Register("soft_delete_query", applySoftDeleteQuery)
	db.Callback().Delete().Replace("gorm:delete", softDelete)
}

func applySoftDeleteQuery(db *gorm.DB) {
	if db.Statement.Unscoped {
		return
	}

	// Apply soft delete condition for tables that have deleted_at column
	if db.Statement.Schema != nil {
		if field := db.Statement.Schema.LookUpField("deleted_at"); field != nil {
			db.Where("deleted_at = ?", DefaultDeletedAt)
		}
	}
}

func softDelete(db *gorm.DB) {
	if db.Statement.Schema != nil {
		if field := db.Statement.Schema.LookUpField("deleted_at"); field != nil {
			// Soft delete by setting deleted_at to current time
			db.UpdateColumn("deleted_at", time.Now())
			return
		}
	}

	// If no deleted_at field, perform hard delete
	db.Exec("DELETE FROM ? WHERE ?", db.Statement.Table, db.Statement.Clauses["WHERE"])
}
