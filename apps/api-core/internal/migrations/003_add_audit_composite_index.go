package migrations

import "gorm.io/gorm"

var Migration003AddAuditCompositeIndex = Migration{
	ID:          "003_add_audit_composite_index",
	Description: "Add composite index for global audit log queries",
	Up: func(db *gorm.DB) error {
		return db.Exec(`
			CREATE INDEX idx_audit_logs_entity_type_created_at 
			ON audit_logs (entity_type, created_at, id)
		`).Error
	},
	Down: func(db *gorm.DB) error {
		return db.Exec("DROP INDEX idx_audit_logs_entity_type_created_at ON audit_logs").Error
	},
}
