package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// metaFieldsToExclude lists fields that should not appear in changed_fields
// These are system-managed fields that always change on update
var metaFieldsToExclude = []string{
	"id",
	"updatedAt",
	"createdAt",
	"updatedBy",
	"createdBy",
}

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
	fields := entity.GetAuditFields()
	newValues, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("failed to marshal new values: %w", err)
	}

	nonEmptyFields := findNonEmptyFields(fields)
	changedFields := filterMetaFields(nonEmptyFields)

	changedFieldsJSON, err := json.Marshal(changedFields)
	if err != nil {
		return fmt.Errorf("failed to marshal changed fields: %w", err)
	}

	auditLog := AuditLog{
		EntityType:    entity.GetAuditEntityType(),
		EntityID:      entity.GetAuditEntityID(),
		Action:        ActionCreate,
		NewValues:     newValues,
		ChangedFields: changedFieldsJSON,
		Username:      userCtx.Username,
	}

	if userCtx.UserID != 0 {
		auditLog.UserID = &userCtx.UserID
	}

	if err := s.db.Session(&gorm.Session{SkipHooks: true}).Create(&auditLog).Error; err != nil {
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
	changedFields = filterMetaFields(changedFields)
	if len(changedFields) == 0 {
		// No business changes detected, skip logging
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

	if err := s.db.Session(&gorm.Session{SkipHooks: true}).Create(&auditLog).Error; err != nil {
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

	if err := s.db.Session(&gorm.Session{SkipHooks: true}).Create(&auditLog).Error; err != nil {
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

func (s *service) GetAllHistory(ctx context.Context, filter AuditLogFilter, page, size int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := s.db.Model(&AuditLog{}).Where("entity_type != ?", "audit_log")

	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.EntityID != nil {
		query = query.Where("entity_id = ?", *filter.EntityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	offset := page * size
	if err := query.Order("created_at DESC, id DESC").
		Limit(size).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, total, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*AuditLog, error) {
	var log AuditLog
	if err := s.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
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

// filterMetaFields removes meta fields from the changed fields list
func filterMetaFields(fields []string) []string {
	var filtered []string
	for _, field := range fields {
		exclude := false
		for _, metaField := range metaFieldsToExclude {
			if field == metaField {
				exclude = true
				break
			}
		}
		if !exclude {
			filtered = append(filtered, field)
		}
	}
	return filtered
}

func findNonEmptyFields(fields map[string]interface{}) []string {
	var nonEmpty []string
	for key, value := range fields {
		if !isEmptyValue(value) {
			nonEmpty = append(nonEmpty, key)
		}
	}
	return nonEmpty
}

func isEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(val).Int() == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(val).Uint() == 0
	case float32, float64:
		return reflect.ValueOf(val).Float() == 0
	case []interface{}:
		return len(val) == 0
	case []map[string]interface{}:
		return len(val) == 0
	case map[string]interface{}:
		return len(val) == 0
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr && rv.IsNil() {
			return true
		}
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			return rv.Len() == 0
		}
		return false
	}
}
