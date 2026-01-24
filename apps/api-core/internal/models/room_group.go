package models

import (
	"gorm.io/gorm"
)

type RoomGroup struct {
	BaseMustAuditEntity
	Name          string `gorm:"type:varchar(20);not null;uniqueIndex:uc_room_group_name,where:deleted_at = '1970-01-01 00:00:00'" json:"name"`
	PeekPrice     int    `gorm:"column:peek_price;not null" json:"peekPrice"`
	OffPeekPrice  int    `gorm:"column:off_peek_price;not null" json:"offPeekPrice"`
	Description   string `gorm:"type:varchar(200);not null" json:"description"`
	Rooms         []Room `gorm:"foreignKey:RoomGroupID" json:"rooms,omitempty"`
	CreatedByUser *User  `gorm:"foreignKey:CreatedBy" json:"createdBy,omitempty"`
	UpdatedByUser *User  `gorm:"foreignKey:UpdatedBy" json:"updatedBy,omitempty"`
}

func (RoomGroup) TableName() string {
	return "room_group"
}

func (rg *RoomGroup) BeforeCreate(tx *gorm.DB) error {
	if err := rg.BaseMustAuditEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if rg.Description == "" {
		rg.Description = ""
	}
	return nil
}
