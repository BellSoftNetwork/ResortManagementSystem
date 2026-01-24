package audit

import (
	"context"
	"encoding/json"
	"time"
)

// Action represents the type of audit action
type Action string

const (
	ActionCreate Action = "CREATE"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"
)

// Auditable interface should be implemented by models that need audit logging
type Auditable interface {
	// GetAuditEntityType returns the entity type for audit logging
	GetAuditEntityType() string
	// GetAuditEntityID returns the entity ID for audit logging
	GetAuditEntityID() uint
	// GetAuditFields returns a map of field names to values for audit logging
	GetAuditFields() map[string]interface{}
}

// AuditService provides methods for audit logging and retrieval
type AuditService interface {
	// LogCreate logs a creation action
	LogCreate(ctx context.Context, entity Auditable) error
	// LogUpdate logs an update action with old and new values
	LogUpdate(ctx context.Context, entity Auditable, oldValues map[string]interface{}) error
	// LogDelete logs a deletion action
	LogDelete(ctx context.Context, entity Auditable) error
	// GetHistory retrieves audit history for an entity
	GetHistory(ctx context.Context, entityType string, entityID uint, page, size int) ([]AuditLog, int64, error)
}

// UserContext represents user information for audit logging
type UserContext struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID            uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	EntityType    string          `gorm:"column:entity_type;size:100;not null" json:"entityType"`
	EntityID      uint            `gorm:"column:entity_id;not null" json:"entityId"`
	Action        Action          `gorm:"column:action;size:20;not null" json:"action"`
	OldValues     json.RawMessage `gorm:"column:old_values;type:json" json:"oldValues,omitempty"`
	NewValues     json.RawMessage `gorm:"column:new_values;type:json" json:"newValues,omitempty"`
	ChangedFields json.RawMessage `gorm:"column:changed_fields;type:json" json:"changedFields,omitempty"`
	UserID        *uint           `gorm:"column:user_id" json:"userId,omitempty"`
	Username      string          `gorm:"column:username;size:100" json:"username,omitempty"`
	CreatedAt     time.Time       `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

// TableName returns the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// GetAuditEntityType implements Auditable interface
func (a AuditLog) GetAuditEntityType() string {
	return "audit_log"
}

// GetAuditEntityID implements Auditable interface
func (a AuditLog) GetAuditEntityID() uint {
	return a.ID
}

// GetAuditFields implements Auditable interface
func (a AuditLog) GetAuditFields() map[string]interface{} {
	return map[string]interface{}{
		"entityType":    a.EntityType,
		"entityId":      a.EntityID,
		"action":        string(a.Action),
		"oldValues":     a.OldValues,
		"newValues":     a.NewValues,
		"changedFields": a.ChangedFields,
		"userId":        a.UserID,
		"username":      a.Username,
		"createdAt":     a.CreatedAt,
	}
}
