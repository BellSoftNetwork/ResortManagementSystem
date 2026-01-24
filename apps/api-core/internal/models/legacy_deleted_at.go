package models

import (
	"database/sql/driver"
	"time"
)

// LegacyDeletedAt is a custom type for handling legacy soft delete pattern
// where non-deleted records have '1970-01-01 00:00:00'
type LegacyDeletedAt time.Time

// DefaultDeletedAt returns the default value for non-deleted records
func DefaultDeletedAt() LegacyDeletedAt {
	return LegacyDeletedAt(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
}

// Scan implements the Scanner interface
func (ldt *LegacyDeletedAt) Scan(value interface{}) error {
	if value == nil {
		*ldt = DefaultDeletedAt()
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*ldt = LegacyDeletedAt(v)
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		*ldt = LegacyDeletedAt(t)
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		*ldt = LegacyDeletedAt(t)
	default:
		*ldt = DefaultDeletedAt()
	}
	return nil
}

// Value implements the driver Valuer interface
func (ldt LegacyDeletedAt) Value() (driver.Value, error) {
	t := time.Time(ldt)
	return t.Format("2006-01-02 15:04:05"), nil
}

// Time returns the time.Time value
func (ldt LegacyDeletedAt) Time() time.Time {
	return time.Time(ldt)
}

// IsDeleted returns true if the record is soft deleted
func (ldt LegacyDeletedAt) IsDeleted() bool {
	defaultTime := time.Time(DefaultDeletedAt())
	return !time.Time(ldt).Equal(defaultTime)
}
