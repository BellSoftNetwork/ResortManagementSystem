package migrations

import "gorm.io/gorm"

var Migration004DeleteAuditLogSelfReferences = Migration{
	ID:          "004_delete_audit_log_self_references",
	Description: "Delete self-referencing audit_log records created before hooks fix",
	Up: func(db *gorm.DB) error {
		return db.Exec("DELETE FROM audit_logs WHERE entity_type = 'audit_log'").Error
	},
	Down: func(db *gorm.DB) error {
		// Self-referencing records are junk data and cannot be restored
		return nil
	},
}
