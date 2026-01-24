package models

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type RoomStatus int8

const (
	RoomStatusDamaged      RoomStatus = -10
	RoomStatusConstruction RoomStatus = -1
	RoomStatusInactive     RoomStatus = 0
	RoomStatusNormal       RoomStatus = 1
)

func (s RoomStatus) String() string {
	switch s {
	case RoomStatusDamaged:
		return "DAMAGED"
	case RoomStatusConstruction:
		return "CONSTRUCTION"
	case RoomStatusInactive:
		return "INACTIVE"
	case RoomStatusNormal:
		return "NORMAL"
	default:
		return "UNKNOWN"
	}
}

func (s RoomStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s *RoomStatus) Scan(value interface{}) error {
	if value == nil {
		*s = RoomStatusInactive
		return nil
	}
	switch v := value.(type) {
	case int64:
		*s = RoomStatus(v)
	case int8:
		*s = RoomStatus(v)
	default:
		*s = RoomStatusInactive
	}
	return nil
}

type Room struct {
	BaseMustAuditEntity
	Number        string     `gorm:"type:varchar(10);not null;uniqueIndex:uc_room_number,where:deleted_at = '1970-01-01 00:00:00'" json:"number"`
	RoomGroupID   uint       `gorm:"column:room_group_id;not null" json:"roomGroupId"`
	RoomGroup     *RoomGroup `gorm:"foreignKey:RoomGroupID" json:"roomGroup,omitempty"`
	Note          string     `gorm:"type:varchar(200);not null" json:"note"`
	Status        RoomStatus `gorm:"type:tinyint;not null" json:"status"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBy" json:"createdBy,omitempty"`
	UpdatedByUser *User      `gorm:"foreignKey:UpdatedBy" json:"updatedBy,omitempty"`
}

func (Room) TableName() string {
	return "room"
}

func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if err := r.BaseMustAuditEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if r.Status == 0 {
		r.Status = RoomStatusInactive
	}
	if r.Note == "" {
		r.Note = ""
	}
	return nil
}

func (r *Room) IsAvailable() bool {
	return r.Status == RoomStatusNormal
}

// GetAuditEntityType implements audit.Auditable interface
func (r *Room) GetAuditEntityType() string {
	return "room"
}

// GetAuditEntityID implements audit.Auditable interface
func (r *Room) GetAuditEntityID() uint {
	return r.ID
}

// GetAuditFields implements audit.Auditable interface
func (r *Room) GetAuditFields() map[string]interface{} {
	return map[string]interface{}{
		"id":          r.ID,
		"number":      r.Number,
		"roomGroupId": r.RoomGroupID,
		"note":        r.Note,
		"status":      r.Status.String(),
		"createdBy":   r.CreatedBy,
		"updatedBy":   r.UpdatedBy,
		"createdAt":   r.CreatedAt,
		"updatedAt":   r.UpdatedAt,
	}
}
