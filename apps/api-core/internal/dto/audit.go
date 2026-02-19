package dto

import (
	"encoding/json"
	"time"
)

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

type AuditLogQuery struct {
	EntityType string `form:"entityType"`
	StartDate  string `form:"startDate"`
	EndDate    string `form:"endDate"`
	Action     string `form:"action"`
	UserID     *uint  `form:"userId"`
	EntityID   *uint  `form:"entityId"`
}

type AuditLogListResponse struct {
	ID            uint      `json:"id"`
	EntityType    string    `json:"entityType"`
	EntityID      uint      `json:"entityId"`
	Action        string    `json:"action"`
	ChangedFields []string  `json:"changedFields"`
	UserID        *uint     `json:"userId,omitempty"`
	Username      string    `json:"username,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}
