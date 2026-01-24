package models

import (
	"database/sql/driver"
	"errors"
)

// BitBool is a custom type to handle MySQL bit(1) fields
type BitBool bool

// Scan implements the sql.Scanner interface
func (b *BitBool) Scan(value interface{}) error {
	if value == nil {
		*b = false
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if len(v) == 1 {
			*b = BitBool(v[0] != 0)
		} else {
			*b = false
		}
	case int64:
		*b = BitBool(v != 0)
	case bool:
		*b = BitBool(v)
	default:
		return errors.New("cannot scan type into BitBool")
	}
	return nil
}

// Value implements the driver.Valuer interface
func (b BitBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// Bool returns the boolean value
func (b BitBool) Bool() bool {
	return bool(b)
}
