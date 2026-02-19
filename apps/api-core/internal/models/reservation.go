package models

import (
	"database/sql/driver"
	"sort"
	"time"

	"gorm.io/gorm"
)

type ReservationStatus int8

const (
	ReservationStatusRefund  ReservationStatus = -10
	ReservationStatusCancel  ReservationStatus = -1
	ReservationStatusPending ReservationStatus = 0
	ReservationStatusNormal  ReservationStatus = 1
)

func (s ReservationStatus) String() string {
	switch s {
	case ReservationStatusRefund:
		return "REFUND"
	case ReservationStatusCancel:
		return "CANCEL"
	case ReservationStatusPending:
		return "PENDING"
	case ReservationStatusNormal:
		return "NORMAL"
	default:
		return "UNKNOWN"
	}
}

func (s ReservationStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s *ReservationStatus) Scan(value interface{}) error {
	if value == nil {
		*s = ReservationStatusPending
		return nil
	}
	switch v := value.(type) {
	case int64:
		*s = ReservationStatus(v)
	case int8:
		*s = ReservationStatus(v)
	default:
		*s = ReservationStatusPending
	}
	return nil
}

type ReservationType int8

const (
	ReservationTypeStay        ReservationType = 0
	ReservationTypeMonthlyRent ReservationType = 10
)

func (t ReservationType) String() string {
	switch t {
	case ReservationTypeStay:
		return "STAY"
	case ReservationTypeMonthlyRent:
		return "MONTHLY_RENT"
	default:
		return "UNKNOWN"
	}
}

func (t ReservationType) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t *ReservationType) Scan(value interface{}) error {
	if value == nil {
		*t = ReservationTypeStay
		return nil
	}
	switch v := value.(type) {
	case int64:
		*t = ReservationType(v)
	case int8:
		*t = ReservationType(v)
	default:
		*t = ReservationTypeStay
	}
	return nil
}

type Reservation struct {
	BaseMustAuditEntity
	PaymentMethodID uint              `gorm:"column:payment_method_id;not null" json:"paymentMethodId"`
	PaymentMethod   *PaymentMethod    `gorm:"foreignKey:PaymentMethodID" json:"paymentMethod,omitempty"`
	Rooms           []ReservationRoom `gorm:"foreignKey:ReservationID" json:"rooms,omitempty"`
	Name            string            `gorm:"column:name;type:varchar(30);not null" json:"name"`
	Phone           string            `gorm:"column:phone;type:varchar(15);not null" json:"phone"`
	PeopleCount     int               `gorm:"column:people_count;not null;default:0" json:"peopleCount"`
	StayStartAt     time.Time         `gorm:"column:stay_start_at;type:date;not null" json:"stayStartAt"`
	StayEndAt       time.Time         `gorm:"column:stay_end_at;type:date;not null" json:"stayEndAt"`
	CheckInAt       *time.Time        `gorm:"column:check_in_at;type:datetime" json:"checkInAt,omitempty"`
	CheckOutAt      *time.Time        `gorm:"column:check_out_at;type:datetime" json:"checkOutAt,omitempty"`
	Price           int               `gorm:"not null" json:"price"`
	Deposit         int               `gorm:"not null;default:0" json:"deposit"`
	PaymentAmount   int               `gorm:"column:payment_amount;not null;default:0" json:"paymentAmount"`
	RefundAmount    int               `gorm:"column:refund_amount;not null;default:0" json:"refundAmount"`
	BrokerFee       int               `gorm:"column:broker_fee;not null;default:0" json:"brokerFee"`
	Note            string            `gorm:"type:varchar(200)" json:"note"`
	CanceledAt      *time.Time        `gorm:"column:canceled_at" json:"canceledAt,omitempty"`
	Status          ReservationStatus `gorm:"type:tinyint;not null;default:0" json:"status"`
	Type            ReservationType   `gorm:"type:tinyint;not null;default:0" json:"type"`
}

func (Reservation) TableName() string {
	return "reservation"
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) error {
	if err := r.BaseMustAuditEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if r.Status == 0 {
		r.Status = ReservationStatusPending
	}
	if r.Type == 0 {
		r.Type = ReservationTypeStay
	}
	if r.Note == "" {
		r.Note = ""
	}
	return nil
}

func (r *Reservation) IsActive() bool {
	return r.Status == ReservationStatusNormal || r.Status == ReservationStatusPending
}

func (r *Reservation) IsCanceled() bool {
	return r.Status == ReservationStatusCancel || r.Status == ReservationStatusRefund
}

func (r *Reservation) GetStayDays() int {
	return int(r.StayEndAt.Sub(r.StayStartAt).Hours() / 24)
}

// GetAuditEntityType implements audit.Auditable interface
func (r *Reservation) GetAuditEntityType() string {
	return "reservation"
}

// GetAuditEntityID implements audit.Auditable interface
func (r *Reservation) GetAuditEntityID() uint {
	return r.ID
}

func formatTimePtr(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format(time.RFC3339)
}

// GetAuditFields implements audit.Auditable interface
func (r *Reservation) GetAuditFields() map[string]interface{} {
	rooms := make([]map[string]interface{}, 0)
	for _, rr := range r.Rooms {
		roomData := map[string]interface{}{
			"id":     rr.RoomID,
			"number": "",
		}
		if rr.Room != nil {
			roomData["number"] = rr.Room.Number
		}
		rooms = append(rooms, roomData)
	}
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i]["id"].(uint) < rooms[j]["id"].(uint)
	})

	paymentMethod := map[string]interface{}{
		"id":   r.PaymentMethodID,
		"name": "",
	}
	if r.PaymentMethod != nil {
		paymentMethod["name"] = r.PaymentMethod.Name
	}

	return map[string]interface{}{
		"id":            r.ID,
		"rooms":         rooms,
		"paymentMethod": paymentMethod,
		"name":          r.Name,
		"phone":         r.Phone,
		"peopleCount":   r.PeopleCount,
		"stayStartAt":   r.StayStartAt.Format("2006-01-02"),
		"stayEndAt":     r.StayEndAt.Format("2006-01-02"),
		"checkInAt":     formatTimePtr(r.CheckInAt),
		"checkOutAt":    formatTimePtr(r.CheckOutAt),
		"price":         r.Price,
		"deposit":       r.Deposit,
		"paymentAmount": r.PaymentAmount,
		"refundAmount":  r.RefundAmount,
		"brokerFee":     r.BrokerFee,
		"note":          r.Note,
		"canceledAt":    formatTimePtr(r.CanceledAt),
		"status":        r.Status.String(),
		"type":          r.Type.String(),
		"createdBy":     r.CreatedBy,
		"updatedBy":     r.UpdatedBy,
		"createdAt":     r.CreatedAt,
		"updatedAt":     r.UpdatedAt,
	}
}
