package migrations

import (
	"crypto/md5"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID          string
	Description string
	AppliedAt   *time.Time
	Up          func(*gorm.DB) error
	Down        func(*gorm.DB) error
}

// AppliedMigration represents a migration that has been applied to the database
type AppliedMigration struct {
	ID            string    `gorm:"primaryKey;type:varchar(255)"`
	Description   string    `gorm:"type:varchar(255)"`
	Checksum      string    `gorm:"type:varchar(32)"`
	AppliedAt     time.Time `gorm:"not null"`
	ExecutionTime int64     `gorm:"not null"` // in milliseconds
}

func (AppliedMigration) TableName() string {
	return "schema_migrations"
}

// MigrationManager manages database migrations
type MigrationManager struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: []Migration{},
	}
}

// AddMigration adds a migration to the manager
func (m *MigrationManager) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// Initialize creates the migration table if it doesn't exist
func (m *MigrationManager) Initialize() error {
	return m.db.AutoMigrate(&AppliedMigration{})
}

// Migrate runs all pending migrations
func (m *MigrationManager) Migrate() error {
	if err := m.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize migration table: %w", err)
	}

	// Sort migrations by ID to ensure consistent order
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].ID < m.migrations[j].ID
	})

	// Get applied migrations
	var applied []AppliedMigration
	if err := m.db.Find(&applied).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[string]AppliedMigration)
	for _, a := range applied {
		appliedMap[a.ID] = a
	}

	// Run pending migrations
	for _, migration := range m.migrations {
		if _, ok := appliedMap[migration.ID]; ok {
			continue // Migration already applied
		}

		start := time.Now()

		// Run migration in a transaction
		err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.ID, err)
			}

			// Record migration
			checksum := m.calculateChecksum(migration)
			record := AppliedMigration{
				ID:            migration.ID,
				Description:   migration.Description,
				Checksum:      checksum,
				AppliedAt:     time.Now(),
				ExecutionTime: time.Since(start).Milliseconds(),
			}

			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.ID, err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		fmt.Printf("Applied migration: %s - %s (%dms)\n", migration.ID, migration.Description, time.Since(start).Milliseconds())
	}

	return nil
}

// Rollback rolls back the last n migrations
func (m *MigrationManager) Rollback(n int) error {
	if n <= 0 {
		return nil
	}

	// Get applied migrations in reverse order
	var applied []AppliedMigration
	if err := m.db.Order("id DESC").Limit(n).Find(&applied).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create a map of migrations by ID
	migrationMap := make(map[string]Migration)
	for _, migration := range m.migrations {
		migrationMap[migration.ID] = migration
	}

	// Rollback migrations
	for _, record := range applied {
		migration, ok := migrationMap[record.ID]
		if !ok {
			return fmt.Errorf("migration %s not found", record.ID)
		}

		if migration.Down == nil {
			return fmt.Errorf("migration %s does not support rollback", record.ID)
		}

		start := time.Now()

		// Run rollback in a transaction
		err := m.db.Transaction(func(tx *gorm.DB) error {
			if err := migration.Down(tx); err != nil {
				return fmt.Errorf("rollback %s failed: %w", migration.ID, err)
			}

			// Remove migration record
			if err := tx.Delete(&record).Error; err != nil {
				return fmt.Errorf("failed to remove migration record %s: %w", migration.ID, err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		fmt.Printf("Rolled back migration: %s - %s (%dms)\n", migration.ID, migration.Description, time.Since(start).Milliseconds())
	}

	return nil
}

// Status returns the migration status
func (m *MigrationManager) Status() ([]MigrationStatus, error) {
	// Get applied migrations
	var applied []AppliedMigration
	if err := m.db.Find(&applied).Error; err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[string]AppliedMigration)
	for _, a := range applied {
		appliedMap[a.ID] = a
	}

	// Sort migrations by ID
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].ID < m.migrations[j].ID
	})

	// Build status list
	var status []MigrationStatus
	for _, migration := range m.migrations {
		s := MigrationStatus{
			ID:          migration.ID,
			Description: migration.Description,
			Applied:     false,
		}

		if applied, ok := appliedMap[migration.ID]; ok {
			s.Applied = true
			s.AppliedAt = &applied.AppliedAt
			s.ExecutionTime = applied.ExecutionTime
		}

		status = append(status, s)
	}

	return status, nil
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	ID            string
	Description   string
	Applied       bool
	AppliedAt     *time.Time
	ExecutionTime int64 // in milliseconds
}

func (m *MigrationManager) calculateChecksum(migration Migration) string {
	h := md5.New()
	h.Write([]byte(migration.ID))
	h.Write([]byte(migration.Description))
	return fmt.Sprintf("%x", h.Sum(nil))
}
