package audit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testEntity is a test implementation of Auditable interface
type testEntity struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100"`
	Age  int
}

func (e *testEntity) GetAuditEntityType() string {
	return "test_entity"
}

func (e *testEntity) GetAuditEntityID() uint {
	return e.ID
}

func (e *testEntity) GetAuditFields() map[string]interface{} {
	return map[string]interface{}{
		"id":   e.ID,
		"name": e.Name,
		"age":  e.Age,
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate test tables
	err = db.AutoMigrate(&AuditLog{}, &testEntity{})
	require.NoError(t, err)

	return db
}

func TestAuditService_LogCreate(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	entity := &testEntity{
		ID:   1,
		Name: "Test Entity",
		Age:  25,
	}

	ctx := SetUserContext(context.Background(), &[]uint{123}[0], "testuser")

	err := service.LogCreate(ctx, entity)
	assert.NoError(t, err)

	// 감사 로그가 정상적으로 생성되었는지 확인
	var auditLog AuditLog
	err = db.Where("entity_type = ? AND entity_id = ? AND action = ?",
		"test_entity", uint(1), ActionCreate).First(&auditLog).Error
	assert.NoError(t, err)

	assert.Equal(t, "test_entity", auditLog.EntityType)
	assert.Equal(t, uint(1), auditLog.EntityID)
	assert.Equal(t, ActionCreate, auditLog.Action)
	assert.Equal(t, uint(123), *auditLog.UserID)
	assert.Equal(t, "testuser", auditLog.Username)

	// NewValues 검증
	var newValues map[string]interface{}
	err = json.Unmarshal(auditLog.NewValues, &newValues)
	assert.NoError(t, err)
	assert.Equal(t, "Test Entity", newValues["name"])
	assert.Equal(t, float64(25), newValues["age"]) // JSON unmarshal converts int to float64
}

func TestAuditService_LogUpdate(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	entity := &testEntity{
		ID:   1,
		Name: "Updated Entity",
		Age:  30,
	}

	oldValues := map[string]interface{}{
		"id":   uint(1),
		"name": "Old Entity",
		"age":  25,
	}

	ctx := SetUserContext(context.Background(), &[]uint{123}[0], "testuser")

	err := service.LogUpdate(ctx, entity, oldValues)
	assert.NoError(t, err)

	// 감사 로그가 정상적으로 생성되었는지 확인
	var auditLog AuditLog
	err = db.Where("entity_type = ? AND entity_id = ? AND action = ?",
		"test_entity", uint(1), ActionUpdate).First(&auditLog).Error
	assert.NoError(t, err)

	assert.Equal(t, ActionUpdate, auditLog.Action)

	// ChangedFields 검증
	var changedFields []string
	err = json.Unmarshal(auditLog.ChangedFields, &changedFields)
	assert.NoError(t, err)
	assert.Contains(t, changedFields, "name")
	assert.Contains(t, changedFields, "age")
}

func TestAuditService_LogUpdate_NoChanges(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	entity := &testEntity{
		ID:   1,
		Name: "Same Entity",
		Age:  25,
	}

	oldValues := map[string]interface{}{
		"id":   uint(1),
		"name": "Same Entity",
		"age":  25,
	}

	ctx := SetUserContext(context.Background(), &[]uint{123}[0], "testuser")

	err := service.LogUpdate(ctx, entity, oldValues)
	assert.NoError(t, err)

	// 변경사항이 없으므로 로그가 생성되지 않아야 함
	var count int64
	db.Model(&AuditLog{}).Where("entity_type = ? AND entity_id = ? AND action = ?",
		"test_entity", uint(1), ActionUpdate).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestAuditService_LogDelete(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	entity := &testEntity{
		ID:   1,
		Name: "Entity to Delete",
		Age:  25,
	}

	ctx := SetUserContext(context.Background(), &[]uint{123}[0], "testuser")

	err := service.LogDelete(ctx, entity)
	assert.NoError(t, err)

	// 감사 로그가 정상적으로 생성되었는지 확인
	var auditLog AuditLog
	err = db.Where("entity_type = ? AND entity_id = ? AND action = ?",
		"test_entity", uint(1), ActionDelete).First(&auditLog).Error
	assert.NoError(t, err)

	assert.Equal(t, ActionDelete, auditLog.Action)
	assert.NotNil(t, auditLog.OldValues)
	assert.Nil(t, auditLog.NewValues)
}

func TestAuditService_GetHistory(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	ctx := SetUserContext(context.Background(), &[]uint{123}[0], "testuser")

	entity := &testEntity{ID: 1, Name: "Test", Age: 25}

	// 여러 개의 감사 로그 생성
	err := service.LogCreate(ctx, entity)
	require.NoError(t, err)

	entity.Name = "Updated"
	oldValues := map[string]interface{}{"id": uint(1), "name": "Test", "age": 25}
	err = service.LogUpdate(ctx, entity, oldValues)
	require.NoError(t, err)

	err = service.LogDelete(ctx, entity)
	require.NoError(t, err)

	// 히스토리 조회
	logs, total, err := service.GetHistory(context.Background(), "test_entity", uint(1), 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, logs, 3)

	// 최신순으로 정렬되었는지 확인
	assert.Equal(t, ActionDelete, logs[0].Action)
	assert.Equal(t, ActionUpdate, logs[1].Action)
	assert.Equal(t, ActionCreate, logs[2].Action)
}

func TestAuditService_GetHistory_WithPagination(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	ctx := SetUserContext(context.Background(), &[]uint{123}[0], "testuser")

	// 5개의 감사 로그 생성
	for i := 1; i <= 5; i++ {
		entity := &testEntity{ID: uint(i), Name: "Test", Age: 25}
		err := service.LogCreate(ctx, entity)
		require.NoError(t, err)
	}

	// 첫 번째 페이지 (크기: 2)
	logs, total, err := service.GetHistory(context.Background(), "test_entity", uint(1), 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total) // entity_id가 1인 것만
	assert.Len(t, logs, 1)

	// 전체 히스토리 조회
	logs, total, err = service.GetHistory(context.Background(), "test_entity", uint(0), 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total) // entity_id가 0인 것은 없음
	assert.Len(t, logs, 0)
}

func TestAuditService_WithoutUserContext(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)

	entity := &testEntity{
		ID:   1,
		Name: "Test Entity",
		Age:  25,
	}

	// 사용자 컨텍스트 없이 실행
	ctx := context.Background()

	err := service.LogCreate(ctx, entity)
	assert.NoError(t, err)

	// 사용자 정보가 빈 값으로 저장되었는지 확인
	var auditLog AuditLog
	err = db.Where("entity_type = ? AND entity_id = ?", "test_entity", uint(1)).First(&auditLog).Error
	assert.NoError(t, err)

	assert.Nil(t, auditLog.UserID)
	assert.Empty(t, auditLog.Username)
}
