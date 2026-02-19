package migrations

import (
	"gorm.io/gorm"
)

// Migration007AddDateBlocks creates the date_block table
var Migration007AddDateBlocks = Migration{
	ID:          "007_add_date_blocks",
	Description: "Create date_block table for managing blocked date ranges",
	Up: func(db *gorm.DB) error {
		return db.Exec(`
			CREATE TABLE date_block (
				id BIGINT PRIMARY KEY AUTO_INCREMENT,
				start_date DATE NOT NULL,
				end_date DATE NOT NULL,
				reason VARCHAR(200) NOT NULL,
				created_at DATETIME NOT NULL,
				updated_at DATETIME NOT NULL,
				deleted_at DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00',
				created_by BIGINT NOT NULL,
				updated_by BIGINT NOT NULL,
				INDEX idx_date_block_dates (start_date, end_date),
				INDEX idx_date_block_deleted_at (deleted_at),
				CONSTRAINT FK_DATE_BLOCK_ON_CREATED_BY FOREIGN KEY (created_by) REFERENCES user (id),
				CONSTRAINT FK_DATE_BLOCK_ON_UPDATED_BY FOREIGN KEY (updated_by) REFERENCES user (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`).Error
	},
	Down: func(db *gorm.DB) error {
		return db.Exec("DROP TABLE IF EXISTS date_block").Error
	},
}
