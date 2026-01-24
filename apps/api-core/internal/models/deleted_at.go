package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CustomDeletedAt represents the legacy soft delete pattern
// Non-deleted records have '1970-01-01 00:00:00'
// Deleted records have the actual deletion timestamp
type CustomDeletedAt struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not the zero value
}

// DefaultDeletedAtTime is the timestamp used for non-deleted records
var DefaultDeletedAtTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

// NewCustomDeletedAt creates a new CustomDeletedAt for non-deleted records
func NewCustomDeletedAt() CustomDeletedAt {
	return CustomDeletedAt{
		Time:  DefaultDeletedAtTime,
		Valid: true,
	}
}

// Scan implements the Scanner interface
func (ct *CustomDeletedAt) Scan(value interface{}) error {
	if value == nil {
		ct.Time = DefaultDeletedAtTime
		ct.Valid = true
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ct.Time = v
		ct.Valid = true
	case []byte:
		if string(v) == "NULL" {
			ct.Time = DefaultDeletedAtTime
			ct.Valid = true
			return nil
		}
		t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		ct.Time = t
		ct.Valid = true
	case string:
		if v == "NULL" {
			ct.Time = DefaultDeletedAtTime
			ct.Valid = true
			return nil
		}
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		ct.Time = t
		ct.Valid = true
	default:
		return fmt.Errorf("cannot scan type %T into CustomDeletedAt", value)
	}

	return nil
}

// Value implements the driver Valuer interface
func (ct CustomDeletedAt) Value() (driver.Value, error) {
	if !ct.Valid {
		return DefaultDeletedAtTime, nil
	}
	return ct.Time, nil
}

// DeleteClauses implements the soft delete interface for GORM
func (ct CustomDeletedAt) DeleteClauses(db *gorm.DB) []clause.Interface {
	return []clause.Interface{}
}

// IsDeleted returns true if the record is soft deleted
func (ct CustomDeletedAt) IsDeleted() bool {
	return ct.Valid && !ct.Time.Equal(DefaultDeletedAtTime)
}

// CustomSoftDeleteClause is the implementation of soft delete
type CustomSoftDeleteClause struct{}

func (CustomSoftDeleteClause) Name() string {
	return "soft_delete"
}

func (CustomSoftDeleteClause) Build(clause.Builder) {
}

func (CustomSoftDeleteClause) MergeClause(*clause.Clause) {
}

func (CustomSoftDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.String() == "" {
		return
	}

	// For DELETE operations, update deleted_at instead of actually deleting
	if stmt.SQL.String() != "" && stmt.SQL.String()[0:6] == "DELETE" {
		stmt.SQL.Reset()
		stmt.SQL.WriteString("UPDATE ")
		stmt.SQL.WriteString(stmt.Table)
		stmt.SQL.WriteString(" SET deleted_at = ? WHERE ")
		stmt.AddVar(stmt, time.Now().Format("2006-01-02 15:04:05"))

		// Add the original WHERE conditions
		if whereClause, ok := stmt.DB.Statement.Clauses["WHERE"]; ok {
			if where, ok := whereClause.Expression.(clause.Where); ok {
				stmt.SQL.WriteString("(")
				where.Build(stmt)
				stmt.SQL.WriteString(")")
			}
		}

		// Add the soft delete condition
		stmt.SQL.WriteString(" AND deleted_at = ?")
		stmt.AddVar(stmt, DefaultDeletedAtTime.Format("2006-01-02 15:04:05"))
	}
}

// QueryClauses implements the soft delete query interface
func (ct CustomDeletedAt) QueryClauses(db *gorm.DB) []clause.Interface {
	return []clause.Interface{}
}

// CustomSoftDeleteQueryClause adds the soft delete condition to queries
type CustomSoftDeleteQueryClause struct{}

func (CustomSoftDeleteQueryClause) Name() string {
	return "soft_delete_query"
}

func (CustomSoftDeleteQueryClause) Build(clause.Builder) {
}

func (CustomSoftDeleteQueryClause) MergeClause(*clause.Clause) {
}

func (CustomSoftDeleteQueryClause) ModifyStatement(stmt *gorm.Statement) {
	if _, ok := stmt.Clauses["soft_delete_query"]; ok {
		return
	}

	// Add WHERE deleted_at = '1970-01-01 00:00:00' for queries
	stmt.AddClause(clause.Where{
		Exprs: []clause.Expression{
			clause.Eq{
				Column: clause.Column{Table: stmt.Table, Name: "deleted_at"},
				Value:  DefaultDeletedAtTime.Format("2006-01-02 15:04:05"),
			},
		},
	})
	stmt.Clauses["soft_delete_query"] = clause.Clause{}
}
