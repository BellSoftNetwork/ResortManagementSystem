package dto

import (
	"time"
)

type ReservationResponse struct {
	ID              uint                   `json:"id"`
	PaymentMethodID uint                   `json:"paymentMethodId"`
	PaymentMethod   *PaymentMethodResponse `json:"paymentMethod,omitempty"`
	Rooms           []RoomResponse         `json:"rooms"` // Spring Boot 호환성을 위해 RoomResponse 직접 사용
	Name            string                 `json:"name"`
	Phone           string                 `json:"phone"`
	PeopleCount     int                    `json:"peopleCount"`
	StayStartAt     JSONDate               `json:"stayStartAt"` // 날짜만 반환
	StayEndAt       JSONDate               `json:"stayEndAt"`   // 날짜만 반환
	CheckInAt       *CustomTime            `json:"checkInAt,omitempty"`
	CheckOutAt      *CustomTime            `json:"checkOutAt,omitempty"`
	Price           int                    `json:"price"`
	Deposit         int                    `json:"deposit"`
	PaymentAmount   int                    `json:"paymentAmount"`
	RefundAmount    int                    `json:"refundAmount"`
	BrokerFee       int                    `json:"brokerFee"`
	Note            string                 `json:"note"`
	CanceledAt      *CustomTime            `json:"canceledAt,omitempty"`
	Status          string                 `json:"status"`
	Type            string                 `json:"type"`
	CreatedAt       CustomTime             `json:"createdAt"`
	UpdatedAt       CustomTime             `json:"updatedAt"`
	CreatedBy       *UserSummaryResponse   `json:"createdBy"` // Spring Boot 호환성
	UpdatedBy       *UserSummaryResponse   `json:"updatedBy"` // Spring Boot 호환성
}

// ReservationRoomResponse는 더 이상 사용하지 않음 - Spring Boot 호환성을 위해 제거

type CreateReservationRequest struct {
	PaymentMethodID uint              `json:"paymentMethodId"`
	PaymentMethod   *EntityReference  `json:"paymentMethod,omitempty"`
	RoomIDs         []uint            `json:"roomIds,omitempty"`
	Rooms           []EntityReference `json:"rooms,omitempty"`
	Name            string            `json:"name" binding:"required,min=2,max=30"`
	Phone           string            `json:"phone" binding:"omitempty,max=20"`
	PeopleCount     int               `json:"peopleCount" binding:"min=0"`
	StayStartAt     JSONTime          `json:"stayStartAt" binding:"required"`
	StayEndAt       JSONTime          `json:"stayEndAt" binding:"required"`
	Price           int               `json:"price" binding:"min=0"`
	Deposit         int               `json:"deposit" binding:"min=0"`
	PaymentAmount   int               `json:"paymentAmount" binding:"min=0"`
	BrokerFee       int               `json:"brokerFee" binding:"min=0"`
	Note            string            `json:"note" binding:"max=200"`
	Status          string            `json:"status,omitempty"`
	Type            string            `json:"type" binding:"omitempty,oneof=STAY MONTHLY_RENT"`
}

func (r *CreateReservationRequest) GetRoomIDs() []uint {
	if len(r.RoomIDs) > 0 {
		return r.RoomIDs
	}
	if len(r.Rooms) > 0 {
		ids := make([]uint, len(r.Rooms))
		for i, room := range r.Rooms {
			ids[i] = room.ID
		}
		return ids
	}
	return nil
}

func (r *CreateReservationRequest) GetPaymentMethodID() uint {
	if r.PaymentMethodID != 0 {
		return r.PaymentMethodID
	}
	if r.PaymentMethod != nil && r.PaymentMethod.ID != 0 {
		return r.PaymentMethod.ID
	}
	return 0
}

// EntityReference is used for frontend compatibility (e.g., {id: 1})
type EntityReference struct {
	ID uint `json:"id"`
}

type UpdateReservationRequest struct {
	PaymentMethodID *uint              `json:"paymentMethodId"`
	PaymentMethod   *EntityReference   `json:"paymentMethod,omitempty"` // 프론트엔드 호환성
	RoomIDs         []uint             `json:"roomIds,omitempty"`
	Rooms           *[]EntityReference `json:"rooms,omitempty"` // 프론트엔드 호환성을 위해 추가
	Name            *string            `json:"name" binding:"omitempty,min=2,max=30"`
	Phone           *string            `json:"phone" binding:"omitempty,max=20"`
	PeopleCount     *int               `json:"peopleCount" binding:"omitempty,min=0"`
	StayStartAt     *time.Time         `json:"stayStartAt"`
	StayEndAt       *time.Time         `json:"stayEndAt"`
	CheckInAt       *time.Time         `json:"checkInAt"`
	CheckOutAt      *time.Time         `json:"checkOutAt"`
	Price           *int               `json:"price" binding:"omitempty,min=0"`
	Deposit         *int               `json:"deposit" binding:"omitempty,min=0"`
	PaymentAmount   *int               `json:"paymentAmount" binding:"omitempty,min=0"`
	RefundAmount    *int               `json:"refundAmount" binding:"omitempty,min=0"`
	Note            *string            `json:"note" binding:"omitempty,max=200"`
	Status          *string            `json:"status" binding:"omitempty,oneof=REFUND CANCEL PENDING NORMAL"`
	Type            *string            `json:"type" binding:"omitempty,oneof=STAY MONTHLY_RENT"`
}

// GetRoomIDs extracts room IDs from either RoomIDs or Rooms field
func (r *UpdateReservationRequest) GetRoomIDs() []uint {
	if len(r.RoomIDs) > 0 {
		return r.RoomIDs
	}
	if r.Rooms != nil {
		ids := make([]uint, len(*r.Rooms))
		for i, room := range *r.Rooms {
			ids[i] = room.ID
		}
		return ids
	}
	return nil
}

// GetPaymentMethodID extracts payment method ID from either PaymentMethodID or PaymentMethod field
func (r *UpdateReservationRequest) GetPaymentMethodID() *uint {
	if r.PaymentMethodID != nil {
		return r.PaymentMethodID
	}
	if r.PaymentMethod != nil && r.PaymentMethod.ID != 0 {
		return &r.PaymentMethod.ID
	}
	return nil
}

// HasRoomsUpdate returns true if rooms were explicitly updated (even if empty)
func (r *UpdateReservationRequest) HasRoomsUpdate() bool {
	return len(r.RoomIDs) > 0 || r.Rooms != nil
}

type ReservationFilter struct {
	Status      *string    `form:"status" binding:"omitempty,oneof=REFUND CANCEL PENDING NORMAL"`
	Type        *string    `form:"type" binding:"omitempty,oneof=STAY MONTHLY_RENT"`
	RoomID      *uint      `form:"roomId"`
	StayStartAt *time.Time `form:"stayStartAt" time_format:"2006-01-02"`
	StayEndAt   *time.Time `form:"stayEndAt" time_format:"2006-01-02"`
	Search      string     `form:"search"`
}

type ReservationStatisticsQuery struct {
	StartDate  time.Time `form:"startDate" binding:"required" time_format:"2006-01-02"`
	EndDate    time.Time `form:"endDate" binding:"required" time_format:"2006-01-02"`
	PeriodType string    `form:"periodType" binding:"omitempty,oneof=DAILY MONTHLY YEARLY"`
}

type ReservationFilterResponse struct {
	Status      *string   `json:"status,omitempty"`
	Type        *string   `json:"type,omitempty"`
	RoomID      *uint     `json:"roomId,omitempty"`
	StayStartAt *JSONDate `json:"stayStartAt,omitempty"`
	StayEndAt   *JSONDate `json:"stayEndAt,omitempty"`
	Search      string    `json:"search,omitempty"`
}

// ReservationStatisticsResponse represents the response for reservation statistics
type ReservationStatisticsResponse struct {
	PeriodType   string           `json:"periodType"`
	Stats        []StatisticsData `json:"stats"`
	MonthlyStats []MonthlyStats   `json:"monthlyStats"`
}

// StatisticsData represents statistics for a specific period
type StatisticsData struct {
	Period            string `json:"period"`
	TotalSales        int    `json:"totalSales"`
	TotalReservations int    `json:"totalReservations"`
	TotalGuests       int    `json:"totalGuests"`
}

// MonthlyStats represents monthly statistics (for backward compatibility)
type MonthlyStats struct {
	YearMonth         string `json:"yearMonth"`
	TotalSales        int    `json:"totalSales"`
	TotalReservations int    `json:"totalReservations"`
	TotalGuests       int    `json:"totalGuests"`
}
