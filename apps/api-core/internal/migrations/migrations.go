package migrations

// AllMigrations returns all available migrations
func AllMigrations() []Migration {
	return []Migration{
		Migration001InitialSchema,
		Migration002AddAuditLogs,
	}
}
