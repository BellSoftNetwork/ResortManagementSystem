package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// service implements AuditService
type service struct {
	db *gorm.DB
}

// NewService creates a new audit service
func NewService(db *gorm.DB) AuditService {
	return &service{db: db}
}

// LogCreate logs a creation action
func (s *service) LogCreate(ctx context.Context, entity Auditable) error {
	userCtx := GetUserContext(ctx)
	newValues, err := json.Marshal(entity.GetAuditFields())
	if err != nil {
		return fmt.Errorf("failed to marshal new values: %w", err)
	}

	auditLog := AuditLog{
		EntityType: entity.GetAuditEntityType(),
		EntityID:   entity.GetAuditEntityID(),
		Action:     ActionCreate,
		NewValues:  newValues,
		Username:   userCtx.Username,
	}

	if userCtx.UserID != 0 {
		auditLog.UserID = &userCtx.UserID
	}

	if err := s.db.Create(&auditLog).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// LogUpdate logs an update action with old and new values
func (s *service) LogUpdate(ctx context.Context, entity Auditable, oldValues map[string]interface{}) error {
	userCtx := GetUserContext(ctx)
	newFields := entity.GetAuditFields()

	// Find changed fields
	changedFields := findChangedFields(oldValues, newFields)
	if len(changedFields) == 0 {
		// No changes detected, skip logging
		return nil
	}

	oldValuesJSON, err := json.Marshal(oldValues)
	if err != nil {
		return fmt.Errorf("failed to marshal old values: %w", err)
	}

	newValuesJSON, err := json.Marshal(newFields)
	if err != nil {
		return fmt.Errorf("failed to marshal new values: %w", err)
	}

	changedFieldsJSON, err := json.Marshal(changedFields)
	if err != nil {
		return fmt.Errorf("failed to marshal changed fields: %w", err)
	}

	auditLog := AuditLog{
		EntityType:    entity.GetAuditEntityType(),
		EntityID:      entity.GetAuditEntityID(),
		Action:        ActionUpdate,
		OldValues:     oldValuesJSON,
		NewValues:     newValuesJSON,
		ChangedFields: changedFieldsJSON,
		Username:      userCtx.Username,
	}

	if userCtx.UserID != 0 {
		auditLog.UserID = &userCtx.UserID
	}

	if err := s.db.Create(&auditLog).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// LogDelete logs a deletion action
func (s *service) LogDelete(ctx context.Context, entity Auditable) error {
	userCtx := GetUserContext(ctx)
	oldValues, err := json.Marshal(entity.GetAuditFields())
	if err != nil {
		return fmt.Errorf("failed to marshal old values: %w", err)
	}

	auditLog := AuditLog{
		EntityType: entity.GetAuditEntityType(),
		EntityID:   entity.GetAuditEntityID(),
		Action:     ActionDelete,
		OldValues:  oldValues,
		Username:   userCtx.Username,
	}

	if userCtx.UserID != 0 {
		auditLog.UserID = &userCtx.UserID
	}

	if err := s.db.Create(&auditLog).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetHistory retrieves audit history for an entity
func (s *service) GetHistory(ctx context.Context, entityType string, entityID uint, page, size int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := s.db.Model(&AuditLog{}).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Apply pagination and get records
	offset := page * size
	if err := query.Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, total, nil
}

// findChangedFields compares old and new values to find changed fields
func findChangedFields(oldValues, newValues map[string]interface{}) []string {
	var changedFields []string

	for key, newValue := range newValues {
		oldValue, exists := oldValues[key]
		if !exists {
			// New field
			changedFields = append(changedFields, key)
		} else if !reflect.DeepEqual(oldValue, newValue) {
			// Changed field
			changedFields = append(changedFields, key)
		}
	}

	// Check for deleted fields
	for key := range oldValues {
		if _, exists := newValues[key]; !exists {
			changedFields = append(changedFields, key)
		}
	}

	return changedFields
}
