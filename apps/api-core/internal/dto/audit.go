package dto

import (
	"encoding/json"
	"time"
)

// AuditLogResponse represents an audit log entry in API responses
type AuditLogResponse struct {
	ID            uint            `json:"id"`
	EntityType    string          `json:"entityType"`
	EntityID      uint            `json:"entityId"`
	Action        string          `json:"action"`
	OldValues     json.RawMessage `json:"oldValues,omitempty"`
	NewValues     json.RawMessage `json:"newValues,omitempty"`
	ChangedFields json.RawMessage `json:"changedFields,omitempty"`
	UserID        *uint           `json:"userId,omitempty"`
	Username      string          `json:"username,omitempty"`
	CreatedAt     time.Time       `json:"createdAt"`
}
