package models

import (
	"time"

	appContext "gitlab.bellsoft.net/rms/api-core/internal/context"
	"gorm.io/gorm"
)

type BaseEntity struct {
	ID uint `gorm:"primarykey;autoIncrement" json:"id"`
}

type BaseTimeEntity struct {
	BaseEntity
	CreatedAt time.Time       `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time       `gorm:"not null" json:"updatedAt"`
	DeletedAt LegacyDeletedAt `gorm:"column:deleted_at;type:datetime;not null;index" json:"-"`
}

type BaseMustAuditEntity struct {
	BaseTimeEntity
	CreatedBy uint `gorm:"column:created_by;not null" json:"createdBy"`
	UpdatedBy uint `gorm:"column:updated_by;not null" json:"updatedBy"`
}

type BaseOptionalAuditEntity struct {
	BaseTimeEntity
	CreatedByID *uint `json:"createdById,omitempty"`
	CreatedBy   *User `gorm:"foreignKey:CreatedByID" json:"-"`
	UpdatedByID *uint `json:"updatedById,omitempty"`
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID" json:"-"`
}

func (b *BaseTimeEntity) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	// Initialize DeletedAt with the default value for non-deleted records
	b.DeletedAt = DefaultDeletedAt()
	return nil
}

func (b *BaseTimeEntity) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}

func (b *BaseMustAuditEntity) BeforeCreate(tx *gorm.DB) error {
	if err := b.BaseTimeEntity.BeforeCreate(tx); err != nil {
		return err
	}

	userID := GetUserIDFromContext(tx)
	if userID == 0 {
		// Use default user ID 1 when userID is not found in context
		// This is a fallback to prevent creation failures
		userID = 1
	}

	b.CreatedBy = userID
	b.UpdatedBy = userID
	return nil
}

func (b *BaseMustAuditEntity) BeforeUpdate(tx *gorm.DB) error {
	if err := b.BaseTimeEntity.BeforeUpdate(tx); err != nil {
		return err
	}

	userID := GetUserIDFromContext(tx)
	if userID == 0 {
		// Use default user ID 1 when userID is not found in context
		// This is a fallback to prevent update failures
		userID = 1
	}

	b.UpdatedBy = userID
	return nil
}

func (b *BaseOptionalAuditEntity) BeforeCreate(tx *gorm.DB) error {
	if err := b.BaseTimeEntity.BeforeCreate(tx); err != nil {
		return err
	}

	if userID := GetUserIDFromContext(tx); userID != 0 {
		b.CreatedByID = &userID
		b.UpdatedByID = &userID
	}

	return nil
}

func (b *BaseOptionalAuditEntity) BeforeUpdate(tx *gorm.DB) error {
	if err := b.BaseTimeEntity.BeforeUpdate(tx); err != nil {
		return err
	}

	if userID := GetUserIDFromContext(tx); userID != 0 {
		b.UpdatedByID = &userID
	}

	return nil
}

func GetUserIDFromContext(tx *gorm.DB) uint {
	if tx.Statement.Context != nil {
		if userID, ok := appContext.GetUserID(tx.Statement.Context); ok {
			return userID
		}
	}
	return 0
}
