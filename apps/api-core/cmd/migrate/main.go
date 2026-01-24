package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gitlab.bellsoft.net/rms/api-core/internal/config"
	"gitlab.bellsoft.net/rms/api-core/internal/database"
	"gitlab.bellsoft.net/rms/api-core/internal/migrations"
)

func main() {
	var (
		action  = flag.String("action", "migrate", "Action to perform: migrate, rollback, or status")
		steps   = flag.Int("steps", 1, "Number of migrations to rollback (only for rollback action)")
		profile = flag.String("profile", "local", "Configuration profile to use")
	)
	flag.Parse()

	// Set profile environment variable
	os.Setenv("PROFILE", *profile)

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create migration manager
	manager := migrations.NewMigrationManager(db)

	// Register all migrations
	for _, migration := range migrations.AllMigrations() {
		manager.AddMigration(migration)
	}

	switch *action {
	case "migrate":
		fmt.Println("Running migrations...")
		if err := manager.Migrate(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "rollback":
		fmt.Printf("Rolling back %d migration(s)...\n", *steps)
		if err := manager.Rollback(*steps); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("Rollback completed successfully")

	case "status":
		fmt.Println("Migration status:")
		status, err := manager.Status()
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

		fmt.Printf("\n%-30s %-50s %-10s %-20s\n", "ID", "Description", "Applied", "Applied At")
		fmt.Println(String(110, "-"))

		for _, s := range status {
			appliedAt := "-"
			if s.AppliedAt != nil {
				appliedAt = s.AppliedAt.Format("2006-01-02 15:04:05")
			}

			applied := "No"
			if s.Applied {
				applied = "Yes"
			}

			fmt.Printf("%-30s %-50s %-10s %-20s\n", s.ID, s.Description, applied, appliedAt)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", *action)
		flag.Usage()
		os.Exit(1)
	}
}

func String(n int, char string) string {
	result := ""
	for i := 0; i < n; i++ {
		result += char
	}
	return result
}
