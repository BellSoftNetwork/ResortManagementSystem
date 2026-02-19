package models

import (
	"time"

	"gorm.io/gorm"
)

type DateBlock struct {
	BaseMustAuditEntity
	StartDate     time.Time `gorm:"column:start_date;type:date;not null" json:"startDate"`
	EndDate       time.Time `gorm:"column:end_date;type:date;not null" json:"endDate"`
	Reason        string    `gorm:"type:varchar(200);not null" json:"reason"`
	CreatedByUser *User     `gorm:"foreignKey:CreatedBy" json:"createdBy,omitempty"`
	UpdatedByUser *User     `gorm:"foreignKey:UpdatedBy" json:"updatedBy,omitempty"`
}

func (DateBlock) TableName() string {
	return "date_block"
}

func (d *DateBlock) BeforeCreate(tx *gorm.DB) error {
	if err := d.BaseMustAuditEntity.BeforeCreate(tx); err != nil {
		return err
	}
	return nil
}

// GetAuditEntityType implements audit.Auditable interface
func (d *DateBlock) GetAuditEntityType() string {
	return "date_block"
}

// GetAuditEntityID implements audit.Auditable interface
func (d *DateBlock) GetAuditEntityID() uint {
	return d.ID
}

// GetAuditFields implements audit.Auditable interface
func (d *DateBlock) GetAuditFields() map[string]interface{} {
	return map[string]interface{}{
		"id":        d.ID,
		"startDate": d.StartDate.Format("2006-01-02"),
		"endDate":   d.EndDate.Format("2006-01-02"),
		"reason":    d.Reason,
		"createdBy": d.CreatedBy,
		"updatedBy": d.UpdatedBy,
		"createdAt": d.CreatedAt,
		"updatedAt": d.UpdatedAt,
	}
}
