# Database Migration System

This project uses a custom migration system to manage database schema changes that matches the api-legacy Liquibase schema.

## Usage

### Run migrations
```bash
go run cmd/migrate/main.go -action=migrate
```

### Check migration status
```bash
go run cmd/migrate/main.go -action=status
```

### Rollback migrations
```bash
# Rollback last migration
go run cmd/migrate/main.go -action=rollback -steps=1

# Rollback multiple migrations
go run cmd/migrate/main.go -action=rollback -steps=3
```

## Creating new migrations

1. Create a new file in `internal/migrations/` with the naming pattern: `YYYYMMDD_NNN_description.go`
2. Define your migration struct following the pattern in existing migrations
3. Add the migration to the `AllMigrations()` function in `internal/migrations/migrations.go`

## Migration Structure

Each migration must implement:
- `ID`: Unique identifier for the migration
- `Description`: Human-readable description
- `Up`: Function to apply the migration
- `Down`: Function to rollback the migration (optional)

## Database Schema

The migration system creates tables that are compatible with api-legacy's Liquibase schema:
- All tables use the same structure as defined in api-legacy
- Foreign key constraints are maintained
- Indexes match the original schema
- Charset and collation are consistent (utf8mb4_unicode_ci)

## Migration Table

The system tracks applied migrations in the `schema_migrations` table with:
- `id`: Migration identifier
- `description`: Migration description
- `checksum`: MD5 hash of migration ID and description
- `applied_at`: Timestamp when migration was applied
- `execution_time`: Time taken to execute migration in milliseconds