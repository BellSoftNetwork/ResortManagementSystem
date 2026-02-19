package migrations

// AllMigrations returns all available migrations
func AllMigrations() []Migration {
	return []Migration{
		Migration001InitialSchema,
		Migration002AddAuditLogs,
		Migration003AddAuditCompositeIndex,
		Migration004DeleteAuditLogSelfReferences,
		Migration005ConvertKstTimestampsToUtc,
		Migration006CleanupFalsePaymentMethodAuditLogs,
		Migration007AddDateBlocks,
	}
}
