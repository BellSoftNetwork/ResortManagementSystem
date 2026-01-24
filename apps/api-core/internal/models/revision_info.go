package models

import (
	"time"
)

type RevisionInfo struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt"`
}

func (RevisionInfo) TableName() string {
	return "revision_info"
}
