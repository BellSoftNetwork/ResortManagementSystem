package dto

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// GenerateDataType represents the type of data to generate
type GenerateDataType string

const (
	GenerateEssentialData   GenerateDataType = "essential"
	GenerateReservationData GenerateDataType = "reservation"
	GenerateAllData         GenerateDataType = "all"
)

// TestDataType is an alias for GenerateDataType (for backward compatibility)
type TestDataType = GenerateDataType

const (
	TestDataTypeEssential   TestDataType = "essential"
	TestDataTypeReservation TestDataType = "reservation"
	TestDataTypeAll         TestDataType = "all"
)

// GenerateTestDataRequest represents the request for generating test data
type GenerateTestDataRequest struct {
	Type GenerateDataType `json:"type" binding:"required,oneof=essential reservation all" example:"all"`

	// Reservation specific options (optional, only used when type is "reservation" or "all")
	ReservationOptions *ReservationGenerationOptions `json:"reservationOptions,omitempty"`
}

// ReservationGenerationOptions represents options for generating reservation data
type ReservationGenerationOptions struct {
	StartDate           *time.Time `json:"startDate,omitempty" example:"2024-01-01T00:00:00Z"`
	EndDate             *time.Time `json:"endDate,omitempty" example:"2024-12-31T23:59:59Z"`
	RegularReservations *int       `json:"regularReservations,omitempty" example:"20"`
	MonthlyReservations *int       `json:"monthlyReservations,omitempty" example:"5"`
}

// Validate performs custom validation for GenerateTestDataRequest
func (r *GenerateTestDataRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	// Validate date range if reservation options are provided
	if r.ReservationOptions != nil {
		if r.ReservationOptions.StartDate != nil && r.ReservationOptions.EndDate != nil {
			if r.ReservationOptions.StartDate.After(*r.ReservationOptions.EndDate) {
				return fmt.Errorf("startDate must be before endDate")
			}
		}
	}

	return nil
}

// GenerateTestDataResponse represents the response after generating test data
type GenerateTestDataResponse struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
