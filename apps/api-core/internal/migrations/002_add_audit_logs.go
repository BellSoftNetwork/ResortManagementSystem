package migrations

import (
	"gorm.io/gorm"
)

// Migration002AddAuditLogs creates the audit_logs table
var Migration002AddAuditLogs = Migration{
	ID:          "002_add_audit_logs",
	Description: "Create audit_logs table for tracking entity changes",
	Up: func(db *gorm.DB) error {
		return db.Exec(`
			CREATE TABLE audit_logs (
				id BIGINT PRIMARY KEY AUTO_INCREMENT,
				entity_type VARCHAR(100) NOT NULL,
				entity_id BIGINT NOT NULL,
				action VARCHAR(20) NOT NULL,
				old_values JSON,
				new_values JSON,
				changed_fields JSON,
				user_id BIGINT,
				username VARCHAR(100),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				INDEX idx_audit_logs_entity (entity_type, entity_id),
				INDEX idx_audit_logs_created_at (created_at),
				INDEX idx_audit_logs_user_id (user_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`).Error
	},
	Down: func(db *gorm.DB) error {
		return db.Exec("DROP TABLE IF EXISTS audit_logs").Error
	},
}
