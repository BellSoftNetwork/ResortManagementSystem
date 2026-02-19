package dto

import (
	"encoding/json"
)

// HistoryType represents the type of revision change
// Matches frontend's HistoryType: "CREATED" | "UPDATED" | "DELETED"
type HistoryType string

const (
	HistoryTypeCreated HistoryType = "CREATED"
	HistoryTypeUpdated HistoryType = "UPDATED"
	HistoryTypeDeleted HistoryType = "DELETED"
)

// ActionToHistoryType converts audit action to frontend-compatible HistoryType
func ActionToHistoryType(action string) HistoryType {
	switch action {
	case "CREATE":
		return HistoryTypeCreated
	case "UPDATE":
		return HistoryTypeUpdated
	case "DELETE":
		return HistoryTypeDeleted
	default:
		return HistoryTypeUpdated
	}
}

// RevisionResponse represents a revision entry matching frontend's Revision<T> type
// Frontend expects: { entity: T, historyType: HistoryType, historyCreatedAt: string, updatedFields: string[] }
type RevisionResponse[T any] struct {
	Entity           T           `json:"entity"`
	HistoryType      HistoryType `json:"historyType"`
	HistoryCreatedAt CustomTime  `json:"historyCreatedAt"`
	UpdatedFields    []string    `json:"updatedFields"`
}

// RoomRevisionResponse is a revision response for Room entity
type RoomRevisionResponse struct {
	Entity           RoomResponse `json:"entity"`
	HistoryType      HistoryType  `json:"historyType"`
	HistoryCreatedAt CustomTime   `json:"historyCreatedAt"`
	UpdatedFields    []string     `json:"updatedFields"`
}

// ReservationRevisionResponse is a revision response for Reservation entity
type ReservationRevisionResponse struct {
	Entity           ReservationResponse `json:"entity"`
	HistoryType      HistoryType         `json:"historyType"`
	HistoryCreatedAt CustomTime          `json:"historyCreatedAt"`
	UpdatedFields    []string            `json:"updatedFields"`
}

// ParseChangedFields parses JSON array of changed fields to string slice
func ParseChangedFields(changedFieldsJSON json.RawMessage) []string {
	if changedFieldsJSON == nil || len(changedFieldsJSON) == 0 {
		return []string{}
	}

	var fields []string
	if err := json.Unmarshal(changedFieldsJSON, &fields); err != nil {
		return []string{}
	}
	return fields
}

// RoomHistorySnapshot represents the room entity data stored in audit logs
type RoomHistorySnapshot struct {
	ID          uint   `json:"id"`
	Number      string `json:"number"`
	RoomGroupID uint   `json:"roomGroupId"`
	Note        string `json:"note"`
	Status      string `json:"status"`
	CreatedBy   uint   `json:"createdBy"`
	UpdatedBy   uint   `json:"updatedBy"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// RoomSnapshot represents a room in audit snapshot (minimal fields for display)
type RoomSnapshot struct {
	ID     uint   `json:"id"`
	Number string `json:"number"`
}

// PaymentMethodSnapshot represents a payment method in audit snapshot (minimal fields for display)
type PaymentMethodSnapshot struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ReservationHistorySnapshot represents the reservation entity data stored in audit logs
type ReservationHistorySnapshot struct {
	ID              uint                   `json:"id"`
	PaymentMethodID uint                   `json:"paymentMethodId"`
	Rooms           []RoomSnapshot         `json:"rooms"`
	PaymentMethod   *PaymentMethodSnapshot `json:"paymentMethod"`
	Name            string                 `json:"name"`
	Phone           string                 `json:"phone"`
	PeopleCount     int                    `json:"peopleCount"`
	StayStartAt     string                 `json:"stayStartAt"`
	StayEndAt       string                 `json:"stayEndAt"`
	CheckInAt       *string                `json:"checkInAt"`
	CheckOutAt      *string                `json:"checkOutAt"`
	Price           int                    `json:"price"`
	Deposit         int                    `json:"deposit"`
	PaymentAmount   int                    `json:"paymentAmount"`
	RefundAmount    int                    `json:"refundAmount"`
	BrokerFee       int                    `json:"brokerFee"`
	Note            string                 `json:"note"`
	CanceledAt      *string                `json:"canceledAt"`
	Status          string                 `json:"status"`
	Type            string                 `json:"type"`
	CreatedBy       uint                   `json:"createdBy"`
	UpdatedBy       uint                   `json:"updatedBy"`
	CreatedAt       string                 `json:"createdAt"`
	UpdatedAt       string                 `json:"updatedAt"`
}

// DateBlockHistorySnapshot represents the date_block entity data stored in audit logs
type DateBlockHistorySnapshot struct {
	ID        uint   `json:"id"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Reason    string `json:"reason"`
	CreatedBy uint   `json:"createdBy"`
	UpdatedBy uint   `json:"updatedBy"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// DateBlockRevisionResponse is a revision response for DateBlock entity
type DateBlockRevisionResponse struct {
	Entity           DateBlockResponse `json:"entity"`
	HistoryType      HistoryType       `json:"historyType"`
	HistoryCreatedAt CustomTime        `json:"historyCreatedAt"`
	UpdatedFields    []string          `json:"updatedFields"`
	HistoryUsername  string            `json:"historyUsername,omitempty"`
}
