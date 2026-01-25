package models

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type UserRole int8

const (
	UserRoleNormal     UserRole = 0
	UserRoleAdmin      UserRole = 100
	UserRoleSuperAdmin UserRole = 127
)

func (r UserRole) String() string {
	switch r {
	case UserRoleNormal:
		return "NORMAL"
	case UserRoleAdmin:
		return "ADMIN"
	case UserRoleSuperAdmin:
		return "SUPER_ADMIN"
	default:
		return "UNKNOWN"
	}
}

func (r UserRole) Value() (driver.Value, error) {
	return int64(r), nil
}

func (r *UserRole) Scan(value interface{}) error {
	if value == nil {
		*r = UserRoleNormal
		return nil
	}
	switch v := value.(type) {
	case int64:
		*r = UserRole(v)
	case int8:
		*r = UserRole(v)
	default:
		*r = UserRoleNormal
	}
	return nil
}

type UserStatus int8

const (
	UserStatusInactive UserStatus = -1
	UserStatusActive   UserStatus = 1
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusInactive:
		return "INACTIVE"
	case UserStatusActive:
		return "ACTIVE"
	default:
		return "UNKNOWN"
	}
}

func (s UserStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s *UserStatus) Scan(value interface{}) error {
	if value == nil {
		*s = UserStatusInactive
		return nil
	}
	switch v := value.(type) {
	case int64:
		*s = UserStatus(v)
	case int8:
		*s = UserStatus(v)
	default:
		*s = UserStatusInactive
	}
	return nil
}

type User struct {
	BaseTimeEntity
	Email    *string    `gorm:"type:varchar(100);uniqueIndex:uc_user_email"`
	Name     string     `gorm:"type:varchar(20);not null" json:"name"`
	Password string     `gorm:"type:varchar(100);not null" json:"-"`
	Status   UserStatus `gorm:"type:tinyint;not null" json:"status"`
	Role     UserRole   `gorm:"type:tinyint;not null" json:"role"`
	UserID   string     `gorm:"column:user_id;type:varchar(30);not null" json:"userId"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.BaseTimeEntity.BeforeCreate(tx); err != nil {
		return err
	}
	if u.Status == 0 {
		u.Status = UserStatusInactive
	}
	if u.Role == 0 {
		u.Role = UserRoleNormal
	}
	return nil
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *User) HasRole(roles ...string) bool {
	userRole := u.Role.String()
	for _, role := range roles {
		if userRole == role {
			return true
		}
	}
	return false
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleSuperAdmin
}

func (u *User) IsSuperAdmin() bool {
	return u.Role == UserRoleSuperAdmin
}
