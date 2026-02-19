package audit

import (
	"context"
	"fmt"
	"reflect"

	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

// AuditableWithOldValues extends Auditable to store old values for updates
type AuditableWithOldValues interface {
	Auditable
	SetOldValues(map[string]interface{})
	GetOldValues() map[string]interface{}
}

// RegisterHooks registers GORM hooks for audit logging on a model
func RegisterHooks(db *gorm.DB, auditService AuditService) {
	// BeforeCreate hook
	db.Callback().Create().Before("gorm:create").Register("audit:before_create", func(db *gorm.DB) {
		if auditable, ok := db.Statement.Dest.(Auditable); ok {
			// Skip AuditLog entities to prevent infinite recursion
			if auditable.GetAuditEntityType() == "audit_log" {
				return
			}
			// Store the context and service for later use
			if db.Statement.Context != nil {
				db.Set("audit_context", db.Statement.Context)
				db.Set("audit_service", auditService)
			}
		}
	})

	// AfterCreate hook
	db.Callback().Create().After("gorm:create").Register("audit:after_create", func(db *gorm.DB) {
		// Get entity directly from Statement.Dest - this is the CURRENT entity being created
		auditable, ok := db.Statement.Dest.(Auditable)
		if !ok {
			return // Not auditable (like testChildEntity), skip
		}

		// Skip AuditLog entities to prevent infinite recursion
		if auditable.GetAuditEntityType() == "audit_log" {
			return
		}

		// Get context and service (these can still use db.Get)
		ctx, ctxExists := db.Get("audit_context")
		service, svcExists := db.Get("audit_service")
		if !ctxExists || !svcExists {
			return
		}

		if auditSvc, ok := service.(AuditService); ok {
			if reqCtx, ok := ctx.(context.Context); ok {
				if err := auditSvc.LogCreate(reqCtx, auditable); err != nil {
					fmt.Printf("Audit log create error: %v\n", err)
				}
			}
		}
	})

	// BeforeUpdate hook - capture old values
	db.Callback().Update().Before("gorm:update").Register("audit:before_update", func(db *gorm.DB) {
		if auditable, ok := db.Statement.Dest.(Auditable); ok {
			// Skip AuditLog entities to prevent infinite recursion
			if auditable.GetAuditEntityType() == "audit_log" {
				return
			}
			// Get old values before update - skip hooks to prevent infinite recursion
			var oldEntity interface{}
			switch auditable.(type) {
			default:
				// Create a new instance of the same type
				oldEntity = reflect.New(reflect.TypeOf(auditable).Elem()).Interface()
			}

			// Use a fresh session without callbacks to avoid infinite recursion
			freshDB := db.Session(&gorm.Session{SkipHooks: true, SkipDefaultTransaction: true})

			if _, isReservation := auditable.(*models.Reservation); isReservation {
				freshDB = freshDB.Preload("PaymentMethod").Preload("Rooms.Room")
			}

			if err := freshDB.Where("id = ?", auditable.GetAuditEntityID()).First(oldEntity).Error; err == nil {
				if oldAuditable, ok := oldEntity.(Auditable); ok {
					oldValues := oldAuditable.GetAuditFields()

					// Store context and old values for after hook
					if db.Statement.Context != nil {
						db.Set("audit_context", db.Statement.Context)
						db.Set("audit_service", auditService)
						db.Set("audit_old_values", oldValues)
					}
				}
			}
		}
	})

	// AfterUpdate hook
	db.Callback().Update().After("gorm:update").Register("audit:after_update", func(db *gorm.DB) {
		// Get entity directly from Statement.Dest - this is the CURRENT entity being updated
		auditable, ok := db.Statement.Dest.(Auditable)
		if !ok {
			return
		}

		// Skip AuditLog entities to prevent infinite recursion
		if auditable.GetAuditEntityType() == "audit_log" {
			return
		}

		// Get context and service (these can still use db.Get)
		ctx, ctxExists := db.Get("audit_context")
		service, svcExists := db.Get("audit_service")
		if !ctxExists || !svcExists {
			return
		}

		// Get old values (this is fine, it's entity-specific)
		oldValues, oldValuesExists := db.Get("audit_old_values")
		if !oldValuesExists {
			return
		}

		if auditSvc, ok := service.(AuditService); ok {
			if reqCtx, ok := ctx.(context.Context); ok {
				if oldVals, ok := oldValues.(map[string]interface{}); ok {
					if err := auditSvc.LogUpdate(reqCtx, auditable, oldVals); err != nil {
						// Log error but don't fail the transaction
						fmt.Printf("Audit log update error: %v\n", err)
					}
				}
			}
		}
	})

	// BeforeDelete hook - capture values before deletion
	db.Callback().Delete().Before("gorm:delete").Register("audit:before_delete", func(db *gorm.DB) {
		if db.Statement.Schema != nil && db.Statement.Context != nil {
			// For delete operations, the entity should be passed via Dest
			if auditable, ok := db.Statement.Dest.(Auditable); ok {
				// Skip AuditLog entities to prevent infinite recursion
				if auditable.GetAuditEntityType() == "audit_log" {
					return
				}
				db.Set("audit_context", db.Statement.Context)
				db.Set("audit_service", auditService)
				db.Set("audit_delete_entity", auditable)
			}
		}
	})

	// AfterDelete hook
	db.Callback().Delete().After("gorm:delete").Register("audit:after_delete", func(db *gorm.DB) {
		if ctx, exists := db.Get("audit_context"); exists {
			if service, exists := db.Get("audit_service"); exists {
				if entity, exists := db.Get("audit_delete_entity"); exists {
					if auditable, ok := entity.(Auditable); ok {
						// Skip AuditLog entities to prevent infinite recursion
						if auditable.GetAuditEntityType() == "audit_log" {
							return
						}
						if auditSvc, ok := service.(AuditService); ok {
							if reqCtx, ok := ctx.(context.Context); ok {
								if err := auditSvc.LogDelete(reqCtx, auditable); err != nil {
									// Log error but don't fail the transaction
									fmt.Printf("Audit log delete error: %v\n", err)
								}
							}
						}
					}
				}
			}
		}
	})
}
