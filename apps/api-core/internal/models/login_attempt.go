package models

import (
	"time"

	"gorm.io/gorm"
)

type LoginAttempt struct {
	BaseEntity
	Username          string    `gorm:"type:varchar(50);not null;index:idx_login_attempts_username_attempt_at" json:"username"`
	IPAddress         string    `gorm:"type:varchar(50);not null;index:idx_login_attempts_ip_address_attempt_at" json:"ipAddress"`
	Successful        bool      `gorm:"not null" json:"successful"`
	AttemptAt         time.Time `gorm:"not null;autoCreateTime:false;index:idx_login_attempts_username_attempt_at,idx_login_attempts_ip_address_attempt_at" json:"attemptAt"`
	OSInfo            *string   `gorm:"type:varchar(50)" json:"osInfo,omitempty"`
	LanguageInfo      *string   `gorm:"type:varchar(50)" json:"languageInfo,omitempty"`
	UserAgent         *string   `gorm:"type:varchar(500)" json:"userAgent,omitempty"`
	DeviceFingerprint *string   `gorm:"type:varchar(50)" json:"deviceFingerprint,omitempty"`
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}

func (la *LoginAttempt) BeforeCreate(tx *gorm.DB) error {
	if la.AttemptAt.IsZero() {
		la.AttemptAt = time.Now()
	}
	return nil
}
