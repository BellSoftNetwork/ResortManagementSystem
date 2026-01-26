package dto

import (
	"gitlab.bellsoft.net/rms/api-core/internal/models"
)

type RoomResponse struct {
	ID          uint                 `json:"id"`
	Number      string               `json:"number"`
	RoomGroupID uint                 `json:"roomGroupId"`
	RoomGroup   *RoomGroupResponse   `json:"roomGroup,omitempty"`
	Note        string               `json:"note"`
	Status      string               `json:"status"`
	CreatedAt   CustomTime           `json:"createdAt"`
	UpdatedAt   CustomTime           `json:"updatedAt"`
	CreatedBy   *UserSummaryResponse `json:"createdBy"` // Spring Boot 호환성
	UpdatedBy   *UserSummaryResponse `json:"updatedBy"` // Spring Boot 호환성
}

type CreateRoomRequest struct {
	Number      string    `json:"number" binding:"required,min=2,max=20"`
	RoomGroupID uint      `json:"roomGroupId"`
	RoomGroup   *struct { // 프론트엔드 호환성을 위해 파싱하지만 무시
		ID uint `json:"id"`
	} `json:"roomGroup,omitempty"`
	Note   string `json:"note" binding:"max=200"`
	Status string `json:"status" binding:"omitempty,oneof=DAMAGED CONSTRUCTION INACTIVE NORMAL"`
}

type UpdateRoomRequest struct {
	Number      *string `json:"number" binding:"omitempty,min=2,max=20"`
	RoomGroupID *uint   `json:"roomGroupId" binding:"omitempty"`
	Note        *string `json:"note" binding:"omitempty,max=200"`
	Status      *string `json:"status" binding:"omitempty,oneof=DAMAGED CONSTRUCTION INACTIVE NORMAL"`
}

type RoomFilter struct {
	RoomGroupID          *uint   `form:"roomGroupId"`
	Status               *string `form:"status" binding:"omitempty,oneof=DAMAGED CONSTRUCTION INACTIVE NORMAL"`
	Search               string  `form:"search"`
	StayStartAt          *string `form:"stayStartAt"`
	StayEndAt            *string `form:"stayEndAt"`
	ExcludeReservationID *uint   `form:"excludeReservationId"`
}

type RoomFilterResponse struct {
	RoomGroupID          *uint   `json:"roomGroupId,omitempty"`
	Status               *string `json:"status,omitempty"`
	Search               string  `json:"search,omitempty"`
	StayStartAt          *string `json:"stayStartAt,omitempty"`
	StayEndAt            *string `json:"stayEndAt,omitempty"`
	ExcludeReservationID *uint   `json:"excludeReservationId,omitempty"`
}

type RoomRepositoryFilter struct {
	RoomGroupID *uint
	Status      *models.RoomStatus
	Search      string
}
