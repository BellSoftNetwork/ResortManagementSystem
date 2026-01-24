package dto

import (
	"fmt"
	"strings"
	"time"
)

// CustomTime is a wrapper around time.Time that marshals to JSON without timezone
type CustomTime struct {
	time.Time
}

// MarshalJSON implements json.Marshaler
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.IsZero() {
		return []byte("null"), nil
	}
	// Format: "2025-06-07T04:05:50" without timezone
	formatted := fmt.Sprintf("\"%s\"", ct.Format("2006-01-02T15:04:05"))
	return []byte(formatted), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return nil
	}

	// Try parsing with various formats
	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
		time.RFC3339,
	}

	var err error
	for _, format := range formats {
		ct.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("cannot unmarshal %q into CustomTime", s)
}

// JSONDate is a wrapper around time.Time that marshals to JSON as date only (YYYY-MM-DD)
type JSONDate struct {
	time.Time
}

// MarshalJSON implements json.Marshaler
func (jd JSONDate) MarshalJSON() ([]byte, error) {
	if jd.IsZero() {
		return []byte("null"), nil
	}
	// Format: "2025-06-07" date only
	formatted := fmt.Sprintf("\"%s\"", jd.Format("2006-01-02"))
	return []byte(formatted), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (jd *JSONDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "null" {
		jd.Time = time.Time{}
		return nil
	}

	// Try parsing with various formats
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
		time.RFC3339,
	}

	var err error
	for _, format := range formats {
		jd.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("cannot unmarshal %q into JSONDate", s)
}
