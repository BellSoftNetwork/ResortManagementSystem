package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

var Migration005ConvertKstTimestampsToUtc = Migration{
	ID:          "005_convert_kst_timestamps_to_utc",
	Description: "Convert legacy KST timestamps to UTC across all tables",
	Up: func(db *gorm.DB) error {
		tables := []string{"user", "room", "room_group", "reservation", "reservation_room", "payment_method"}
		for _, table := range tables {
			var count int64
			if err := db.Raw(fmt.Sprintf(
				"SELECT COUNT(*) FROM `%s` WHERE created_at BETWEEN '2026-01-24 20:02:11' AND '2026-01-25 05:02:11'",
				table,
			)).Scan(&count).Error; err != nil {
				return fmt.Errorf("failed to check ambiguous window for %s: %w", table, err)
			}
			if count > 0 {
				fmt.Printf("[WARNING] Table %s has %d records in ambiguous 9h window. Proceeding with migration.\n", table, count)
			}
		}

		historyTables := []string{"room_history", "reservation_history", "reservation_room_history"}
		historyColumns := []string{"created_at", "updated_at", "deleted_at"}
		for _, table := range historyTables {
			for _, col := range historyColumns {
				if err := db.Exec(fmt.Sprintf(
					"UPDATE `%s` SET `%s` = DATE_SUB(`%s`, INTERVAL 9 HOUR) WHERE `%s` > '1000-01-01 00:00:00' AND `%s` != '1970-01-01 00:00:00'",
					table, col, col, col, col,
				)).Error; err != nil {
					return fmt.Errorf("failed to convert %s.%s: %w", table, col, err)
				}
			}
		}

		mainTables := []string{"user", "room", "room_group", "reservation", "reservation_room", "payment_method"}
		for _, table := range mainTables {
			if err := db.Exec(fmt.Sprintf(
				"UPDATE `%s` SET created_at = DATE_SUB(created_at, INTERVAL 9 HOUR) WHERE created_at <= '2026-01-25 05:02:11' AND created_at != '1970-01-01 00:00:00' AND created_at > '1000-01-01 00:00:00'",
				table,
			)).Error; err != nil {
				return fmt.Errorf("failed to convert %s.created_at: %w", table, err)
			}
		}

		if err := db.Exec(
			"UPDATE login_attempts SET attempt_at = DATE_SUB(attempt_at, INTERVAL 9 HOUR) WHERE attempt_at <= '2026-01-25 05:02:11' AND attempt_at > '1000-01-01 00:00:00'",
		).Error; err != nil {
			return fmt.Errorf("failed to convert login_attempts.attempt_at: %w", err)
		}

		for _, table := range mainTables {
			if err := db.Exec(fmt.Sprintf(
				"UPDATE `%s` SET updated_at = DATE_SUB(updated_at, INTERVAL 9 HOUR) WHERE updated_at <= '2026-01-24 20:02:11' AND updated_at != '1970-01-01 00:00:00' AND updated_at > '1000-01-01 00:00:00'",
				table,
			)).Error; err != nil {
				return fmt.Errorf("failed to convert %s.updated_at: %w", table, err)
			}
		}

		deletedAtTables := []string{"user", "room", "room_group", "reservation", "reservation_room", "payment_method"}
		for _, table := range deletedAtTables {
			if err := db.Exec(fmt.Sprintf(
				"UPDATE `%s` SET deleted_at = DATE_SUB(deleted_at, INTERVAL 9 HOUR) WHERE deleted_at <= '2026-01-25 05:02:11' AND deleted_at != '1970-01-01 00:00:00' AND deleted_at > '1000-01-01 00:00:00'",
				table,
			)).Error; err != nil {
				return fmt.Errorf("failed to convert %s.deleted_at: %w", table, err)
			}
		}

		nullableCols := []string{"check_in_at", "check_out_at", "canceled_at"}
		for _, col := range nullableCols {
			if err := db.Exec(fmt.Sprintf(
				"UPDATE reservation SET `%s` = DATE_SUB(`%s`, INTERVAL 9 HOUR) WHERE `%s` IS NOT NULL AND `%s` <= '2026-01-24 20:02:11' AND `%s` > '1000-01-01 00:00:00'",
				col, col, col, col, col,
			)).Error; err != nil {
				return fmt.Errorf("failed to convert reservation.%s: %w", col, err)
			}
		}

		return nil
	},
	Down: func(db *gorm.DB) error {
		historyTables := []string{"room_history", "reservation_history", "reservation_room_history"}
		historyColumns := []string{"created_at", "updated_at", "deleted_at"}
		for _, table := range historyTables {
			for _, col := range historyColumns {
				if err := db.Exec(fmt.Sprintf(
					"UPDATE `%s` SET `%s` = DATE_ADD(`%s`, INTERVAL 9 HOUR) WHERE `%s` > '1000-01-01 00:00:00' AND `%s` != '1970-01-01 00:00:00'",
					table, col, col, col, col,
				)).Error; err != nil {
					return fmt.Errorf("failed to rollback %s.%s: %w", table, col, err)
				}
			}
		}

		mainTables := []string{"user", "room", "room_group", "reservation", "reservation_room", "payment_method"}
		for _, table := range mainTables {
			if err := db.Exec(fmt.Sprintf(
				"UPDATE `%s` SET created_at = DATE_ADD(created_at, INTERVAL 9 HOUR) WHERE created_at <= '2026-01-24 20:02:11' AND created_at != '1970-01-01 00:00:00' AND created_at > '1000-01-01 00:00:00'",
				table,
			)).Error; err != nil {
				return fmt.Errorf("failed to rollback %s.created_at: %w", table, err)
			}
		}

		if err := db.Exec(
			"UPDATE login_attempts SET attempt_at = DATE_ADD(attempt_at, INTERVAL 9 HOUR) WHERE attempt_at <= '2026-01-24 20:02:11' AND attempt_at > '1000-01-01 00:00:00'",
		).Error; err != nil {
			return fmt.Errorf("failed to rollback login_attempts.attempt_at: %w", err)
		}

		for _, table := range mainTables {
			if err := db.Exec(fmt.Sprintf(
				"UPDATE `%s` SET updated_at = DATE_ADD(updated_at, INTERVAL 9 HOUR) WHERE updated_at <= '2026-01-24 11:02:11' AND updated_at != '1970-01-01 00:00:00' AND updated_at > '1000-01-01 00:00:00'",
				table,
			)).Error; err != nil {
				return fmt.Errorf("failed to rollback %s.updated_at: %w", table, err)
			}
		}

		deletedAtTables := []string{"user", "room", "room_group", "reservation", "reservation_room", "payment_method"}
		for _, table := range deletedAtTables {
			if err := db.Exec(fmt.Sprintf(
				"UPDATE `%s` SET deleted_at = DATE_ADD(deleted_at, INTERVAL 9 HOUR) WHERE deleted_at <= '2026-01-24 20:02:11' AND deleted_at != '1970-01-01 00:00:00' AND deleted_at > '1000-01-01 00:00:00'",
				table,
			)).Error; err != nil {
				return fmt.Errorf("failed to rollback %s.deleted_at: %w", table, err)
			}
		}

		nullableCols := []string{"check_in_at", "check_out_at", "canceled_at"}
		for _, col := range nullableCols {
			if err := db.Exec(fmt.Sprintf(
				"UPDATE reservation SET `%s` = DATE_ADD(`%s`, INTERVAL 9 HOUR) WHERE `%s` IS NOT NULL AND `%s` <= '2026-01-24 11:02:11' AND `%s` > '1000-01-01 00:00:00'",
				col, col, col, col, col,
			)).Error; err != nil {
				return fmt.Errorf("failed to rollback reservation.%s: %w", col, err)
			}
		}

		return nil
	},
}
