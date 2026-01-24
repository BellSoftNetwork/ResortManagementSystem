package models

import (
	"gorm.io/gorm"
)

type ReservationRoom struct {
	BaseMustAuditEntity
	ReservationID uint         `gorm:"column:reservation_id;not null;uniqueIndex:uc_reservation_room_reservation_id_and_room_id,where:deleted_at = '1970-01-01 00:00:00'" json:"reservationId"`
	Reservation   *Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	RoomID        uint         `gorm:"column:room_id;not null;uniqueIndex:uc_reservation_room_reservation_id_and_room_id,where:deleted_at = '1970-01-01 00:00:00'" json:"roomId"`
	Room          *Room        `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

func (ReservationRoom) TableName() string {
	return "reservation_room"
}

func (rr *ReservationRoom) BeforeCreate(tx *gorm.DB) error {
	return rr.BaseMustAuditEntity.BeforeCreate(tx)
}
