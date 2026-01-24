package dto

type RoomGroupResponse struct {
	ID           uint                         `json:"id"`
	Name         string                       `json:"name"`
	PeekPrice    int                          `json:"peekPrice"`
	OffPeekPrice int                          `json:"offPeekPrice"`
	Description  string                       `json:"description"`
	Rooms        []RoomLastStayDetailResponse `json:"rooms"`
	CreatedAt    CustomTime                   `json:"createdAt"`
	CreatedBy    *UserSummaryResponse         `json:"createdBy"`
	UpdatedAt    CustomTime                   `json:"updatedAt"`
	UpdatedBy    *UserSummaryResponse         `json:"updatedBy"`
}

// RoomLastStayDetailResponse represents a room with its last reservation
type RoomLastStayDetailResponse struct {
	Room            RoomResponse         `json:"room"`
	LastReservation *ReservationResponse `json:"lastReservation"`
}

type CreateRoomGroupRequest struct {
	Name         string `json:"name" binding:"required,min=2,max=20"`
	PeekPrice    int    `json:"peekPrice" binding:"min=0"`
	OffPeekPrice int    `json:"offPeekPrice" binding:"min=0"`
	Description  string `json:"description" binding:"max=200"`
}

type UpdateRoomGroupRequest struct {
	Name         *string `json:"name" binding:"omitempty,min=2,max=20"`
	PeekPrice    *int    `json:"peekPrice" binding:"omitempty,min=0"`
	OffPeekPrice *int    `json:"offPeekPrice" binding:"omitempty,min=0"`
	Description  *string `json:"description" binding:"omitempty,max=200"`
}

// RoomGroupRoomFilter는 룸그룹 내 객실 조회 시 필터링 조건
type RoomGroupRoomFilter struct {
	Status               *string `form:"status"`
	StayStartAt          *string `form:"stayStartAt"`
	StayEndAt            *string `form:"stayEndAt"`
	ExcludeReservationID *uint   `form:"excludeReservationId"`
}
