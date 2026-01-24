package models

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type PaymentMethodStatus int8

const (
	PaymentMethodStatusInactive PaymentMethodStatus = -1
	PaymentMethodStatusActive   PaymentMethodStatus = 1
)

func (s PaymentMethodStatus) String() string {
	switch s {
	case PaymentMethodStatusInactive:
		return "INACTIVE"
	case PaymentMethodStatusActive:
		return "ACTIVE"
	default:
		return "UNKNOWN"
	}
}

func (s PaymentMethodStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s *PaymentMethodStatus) Scan(value interface{}) error {
	if value == nil {
		*s = PaymentMethodStatusInactive
		return nil
	}
	switch v := value.(type) {
	case int64:
		*s = PaymentMethodStatus(v)
	case int8:
		*s = PaymentMethodStatus(v)
	default:
		*s = PaymentMethodStatusInactive
	}
	return nil
}

type PaymentMethod struct {
	BaseTimeEntity
	Name                     string              `gorm:"type:varchar(20);not null;uniqueIndex:uc_payment_method_name,where:deleted_at = '1970-01-01 00:00:00'" json:"name"`
	CommissionRate           float64             `gorm:"column:commission_rate;not null" json:"commissionRate"`
	RequireUnpaidAmountCheck BitBool             `gorm:"column:required_unpaid_amount_check;type:bit(1);not null" json:"requireUnpaidAmountCheck"`
	IsDefaultSelect          BitBool             `gorm:"column:is_default_select;type:bit(1);not null" json:"isDefaultSelect"`
	Status                   PaymentMethodStatus `gorm:"type:tinyint;not null" json:"status"`
}

func (PaymentMethod) TableName() string {
	return "payment_method"
}

func (pm *PaymentMethod) BeforeCreate(tx *gorm.DB) error {
	if err := pm.BaseTimeEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if pm.Status == 0 {
		pm.Status = PaymentMethodStatusInactive
	}
	return nil
}

func (pm *PaymentMethod) BeforeUpdate(tx *gorm.DB) error {
	return pm.BaseTimeEntity.BeforeUpdate(tx)
}

func (pm *PaymentMethod) IsActive() bool {
	return pm.Status == PaymentMethodStatusActive
}
