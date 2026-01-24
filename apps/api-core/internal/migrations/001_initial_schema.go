package migrations

import (
	"gorm.io/gorm"
)

// Migration001InitialSchema creates the complete initial database schema
// This migration creates all tables with the exact structure matching api-legacy
var Migration001InitialSchema = Migration{
	ID:          "001_initial_schema",
	Description: "Create initial database schema matching api-legacy structure",
	Up: func(db *gorm.DB) error {
		// 1. Create revision_info table (for audit history)
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS revision_info (
				id bigint NOT NULL AUTO_INCREMENT,
				timestamp bigint NOT NULL,
				user_id bigint DEFAULT NULL,
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		// 2. Create user table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS user (
				id bigint NOT NULL AUTO_INCREMENT,
				email varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
				password varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
				role tinyint NOT NULL,
				name varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
				status tinyint NOT NULL,
				created_at datetime NOT NULL,
				updated_at datetime NOT NULL,
				deleted_at datetime NOT NULL,
				user_id varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL,
				PRIMARY KEY (id),
				UNIQUE KEY uc_user_user_id (user_id,deleted_at),
				UNIQUE KEY uc_user_email (email,deleted_at)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		// 3. Create login_attempts table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS login_attempts (
				id bigint NOT NULL AUTO_INCREMENT,
				username varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '계정 ID',
				ip_address varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP 주소',
				successful bit(1) NOT NULL COMMENT '로그인 성공 여부',
				attempt_at datetime NOT NULL COMMENT '로그인 시도 시각',
				os_info varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '운영체제 정보',
				language_info varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '언어 설정 정보',
				user_agent varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '사용자 에이전트 정보',
				device_fingerprint varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '디바이스 정보',
				PRIMARY KEY (id),
				KEY idx_login_attempts_ip_address_attempt_at (ip_address,attempt_at),
				KEY idx_login_attempts_username_attempt_at (username,attempt_at),
				KEY idx_login_attempts_username_ip_address_attempt_at (username,ip_address,attempt_at)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='로그인 시도 이력'
		`).Error; err != nil {
			return err
		}

		// 4. Create payment_method table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS payment_method (
				id bigint NOT NULL AUTO_INCREMENT,
				name varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
				commission_rate double NOT NULL,
				status tinyint NOT NULL,
				created_at datetime NOT NULL,
				updated_at datetime NOT NULL,
				deleted_at datetime NOT NULL,
				required_unpaid_amount_check bit(1) NOT NULL,
				is_default_select bit(1) NOT NULL,
				PRIMARY KEY (id),
				UNIQUE KEY uc_payment_method_name (name,deleted_at)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		// 5. Create room_group table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS room_group (
				id bigint NOT NULL AUTO_INCREMENT,
				name varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL,
				description varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
				created_by bigint NOT NULL,
				updated_by bigint NOT NULL,
				created_at datetime NOT NULL,
				updated_at datetime NOT NULL,
				deleted_at datetime NOT NULL,
				off_peek_price int NOT NULL,
				peek_price int NOT NULL,
				PRIMARY KEY (id),
				UNIQUE KEY uc_room_group_name (name,deleted_at),
				KEY FK_ROOM_GROUP_ON_CREATED_BY (created_by),
				KEY FK_ROOM_GROUP_ON_UPDATED_BY (updated_by),
				CONSTRAINT FK_ROOM_GROUP_ON_CREATED_BY FOREIGN KEY (created_by) REFERENCES user (id),
				CONSTRAINT FK_ROOM_GROUP_ON_UPDATED_BY FOREIGN KEY (updated_by) REFERENCES user (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		// 6. Create room table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS room (
				id bigint NOT NULL AUTO_INCREMENT,
				number varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
				status tinyint NOT NULL,
				note varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
				created_at datetime NOT NULL,
				updated_at datetime NOT NULL,
				deleted_at datetime NOT NULL,
				created_by bigint NOT NULL,
				updated_by bigint NOT NULL,
				room_group_id bigint NOT NULL,
				PRIMARY KEY (id),
				UNIQUE KEY uc_room_number (number,deleted_at),
				KEY FK_ROOM_ON_CREATED_BY (created_by),
				KEY FK_ROOM_ON_UPDATED_BY (updated_by),
				KEY FK_ROOM_ON_ROOM_GROUP (room_group_id),
				CONSTRAINT FK_ROOM_ON_CREATED_BY FOREIGN KEY (created_by) REFERENCES user (id),
				CONSTRAINT FK_ROOM_ON_ROOM_GROUP FOREIGN KEY (room_group_id) REFERENCES room_group (id),
				CONSTRAINT FK_ROOM_ON_UPDATED_BY FOREIGN KEY (updated_by) REFERENCES user (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		// 7. Create reservation table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS reservation (
				id bigint NOT NULL AUTO_INCREMENT,
				payment_method_id bigint NOT NULL,
				name varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL,
				phone varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
				people_count int NOT NULL,
				stay_start_at date NOT NULL,
				stay_end_at date NOT NULL,
				check_in_at datetime DEFAULT NULL,
				check_out_at datetime DEFAULT NULL,
				price int NOT NULL,
				payment_amount int NOT NULL,
				refund_amount int NOT NULL,
				broker_fee int NOT NULL,
				note varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
				status tinyint NOT NULL,
				canceled_at datetime DEFAULT NULL,
				created_at datetime NOT NULL,
				updated_at datetime NOT NULL,
				deleted_at datetime NOT NULL,
				created_by bigint NOT NULL,
				updated_by bigint NOT NULL,
				type tinyint NOT NULL,
				deposit int NOT NULL,
				PRIMARY KEY (id),
				KEY FK_RESERVATION_ON_RESERVATION_METHOD (payment_method_id),
				KEY FK_RESERVATION_ON_CREATED_BY (created_by),
				KEY FK_RESERVATION_ON_UPDATED_BY (updated_by),
				CONSTRAINT FK_RESERVATION_ON_CREATED_BY FOREIGN KEY (created_by) REFERENCES user (id),
				CONSTRAINT FK_RESERVATION_ON_RESERVATION_METHOD FOREIGN KEY (payment_method_id) REFERENCES payment_method (id),
				CONSTRAINT FK_RESERVATION_ON_UPDATED_BY FOREIGN KEY (updated_by) REFERENCES user (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		// 8. Create reservation_room table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS reservation_room (
				id bigint NOT NULL AUTO_INCREMENT,
				reservation_id bigint NOT NULL,
				room_id bigint NOT NULL,
				created_by bigint NOT NULL,
				updated_by bigint NOT NULL,
				created_at datetime NOT NULL,
				updated_at datetime NOT NULL,
				deleted_at datetime NOT NULL,
				PRIMARY KEY (id),
				KEY FK_RESERVATION_ROOM_ON_RESERVATION (reservation_id),
				KEY FK_RESERVATION_ROOM_ON_ROOM (room_id),
				KEY FK_RESERVATION_ROOM_ON_CREATED_BY (created_by),
				KEY FK_RESERVATION_ROOM_ON_UPDATED_BY (updated_by),
				CONSTRAINT FK_RESERVATION_ROOM_ON_CREATED_BY FOREIGN KEY (created_by) REFERENCES user (id),
				CONSTRAINT FK_RESERVATION_ROOM_ON_RESERVATION FOREIGN KEY (reservation_id) REFERENCES reservation (id),
				CONSTRAINT FK_RESERVATION_ROOM_ON_ROOM FOREIGN KEY (room_id) REFERENCES room (id),
				CONSTRAINT FK_RESERVATION_ROOM_ON_UPDATED_BY FOREIGN KEY (updated_by) REFERENCES user (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		// 9. Create room_history table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS room_history (
				rev bigint NOT NULL,
				revtype tinyint DEFAULT NULL,
				number varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
				number_mod bit(1) DEFAULT NULL,
				note varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
				note_mod bit(1) DEFAULT NULL,
				status tinyint DEFAULT NULL,
				status_mod bit(1) DEFAULT NULL,
				id bigint NOT NULL,
				created_by bigint DEFAULT NULL,
				updated_by bigint DEFAULT NULL,
				created_at datetime DEFAULT NULL,
				deleted_at datetime DEFAULT NULL,
				updated_at datetime DEFAULT NULL,
				created_by_mod bit(1) DEFAULT NULL,
				updated_by_mod bit(1) DEFAULT NULL,
				room_group_id bigint DEFAULT NULL,
				room_group_id_mod bit(1) DEFAULT NULL,
				PRIMARY KEY (rev,id),
				CONSTRAINT FK_ROOM_HISTORY_ON_REV FOREIGN KEY (rev) REFERENCES revision_info (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		// 10. Create reservation_history table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS reservation_history (
				rev bigint NOT NULL,
				revtype tinyint DEFAULT NULL,
				status tinyint DEFAULT NULL,
				status_mod bit(1) DEFAULT NULL,
				id bigint NOT NULL,
				user_id bigint DEFAULT NULL,
				user_id_mod bit(1) DEFAULT NULL,
				payment_method_id bigint DEFAULT NULL,
				payment_method_id_mod bit(1) DEFAULT NULL,
				name varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
				name_mod bit(1) DEFAULT NULL,
				phone varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
				phone_mod bit(1) DEFAULT NULL,
				people_count int DEFAULT NULL,
				people_count_mod bit(1) DEFAULT NULL,
				stay_start_at date DEFAULT NULL,
				stay_start_at_mod bit(1) DEFAULT NULL,
				stay_end_at date DEFAULT NULL,
				stay_end_at_mod bit(1) DEFAULT NULL,
				check_in_at datetime DEFAULT NULL,
				check_in_at_mod bit(1) DEFAULT NULL,
				check_out_at datetime DEFAULT NULL,
				check_out_at_mod bit(1) DEFAULT NULL,
				price int DEFAULT NULL,
				price_mod bit(1) DEFAULT NULL,
				payment_amount int DEFAULT NULL,
				payment_amount_mod bit(1) DEFAULT NULL,
				refund_amount int DEFAULT NULL,
				refund_amount_mod bit(1) DEFAULT NULL,
				broker_fee int DEFAULT NULL,
				broker_fee_mod bit(1) DEFAULT NULL,
				note varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
				note_mod bit(1) DEFAULT NULL,
				canceled_at datetime DEFAULT NULL,
				canceled_at_mod bit(1) DEFAULT NULL,
				created_by bigint DEFAULT NULL,
				updated_by bigint DEFAULT NULL,
				created_at datetime DEFAULT NULL,
				deleted_at datetime DEFAULT NULL,
				updated_at datetime DEFAULT NULL,
				type tinyint DEFAULT NULL,
				type_mod bit(1) DEFAULT NULL,
				deposit int DEFAULT NULL,
				deposit_mod bit(1) DEFAULT NULL,
				PRIMARY KEY (rev,id),
				CONSTRAINT FK_RESERVATION_HISTORY_ON_REV FOREIGN KEY (rev) REFERENCES revision_info (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		// 11. Create reservation_room_history table
		if err := db.Exec(`
			CREATE TABLE IF NOT EXISTS reservation_room_history (
				rev bigint NOT NULL,
				revtype tinyint DEFAULT NULL,
				id bigint NOT NULL,
				reservation_id bigint DEFAULT NULL,
				reservation_id_mod bit(1) DEFAULT NULL,
				room_id bigint DEFAULT NULL,
				room_id_mod bit(1) DEFAULT NULL,
				created_by bigint DEFAULT NULL,
				updated_by bigint DEFAULT NULL,
				created_at datetime DEFAULT NULL,
				deleted_at datetime DEFAULT NULL,
				updated_at datetime DEFAULT NULL,
				PRIMARY KEY (rev,id),
				CONSTRAINT FK_RESERVATION_ROOM_HISTORY_ON_REV FOREIGN KEY (rev) REFERENCES revision_info (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`).Error; err != nil {
			return err
		}

		return nil
	},
	Down: func(db *gorm.DB) error {
		// Drop tables in reverse order of dependencies
		tables := []string{
			"reservation_room_history",
			"reservation_history",
			"room_history",
			"reservation_room",
			"reservation",
			"room",
			"room_group",
			"payment_method",
			"user",
			"login_attempts",
			"revision_info",
		}

		for _, table := range tables {
			if err := db.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
				return err
			}
		}

		return nil
	},
}
