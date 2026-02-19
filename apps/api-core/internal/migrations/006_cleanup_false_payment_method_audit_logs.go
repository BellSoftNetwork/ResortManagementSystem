package migrations

import "gorm.io/gorm"

var Migration006CleanupFalsePaymentMethodAuditLogs = Migration{
	ID:          "006_cleanup_false_payment_method_audit_logs",
	Description: "Delete duplicate payment_method CREATE audit logs from FullSaveAssociations re-saves",
	Up: func(db *gorm.DB) error {
		return db.Exec(`
			DELETE a FROM audit_logs a
			LEFT JOIN (
				SELECT MIN(id) as keep_id
				FROM audit_logs
				WHERE entity_type = 'payment_method' AND action = 'CREATE'
				GROUP BY entity_id
			) b ON a.id = b.keep_id
			WHERE a.entity_type = 'payment_method'
			AND a.action = 'CREATE'
			AND b.keep_id IS NULL
		`).Error
	},
	Down: func(db *gorm.DB) error {
		// Deleted records were junk data from FullSaveAssociations re-saves — cannot be restored
		return nil
	},
}
