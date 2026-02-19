package audit

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testParentEntity is a test implementation of Auditable interface (simulates Reservation)
type testParentEntity struct {
	ID     uint              `gorm:"primaryKey"`
	Name   string            `gorm:"size:100"`
	Childs []testChildEntity `gorm:"foreignKey:ParentID"`
}

func (testParentEntity) TableName() string {
	return "test_parent_entities"
}

func (e *testParentEntity) GetAuditEntityType() string {
	return "test_parent"
}

func (e *testParentEntity) GetAuditEntityID() uint {
	return e.ID
}

func (e *testParentEntity) GetAuditFields() map[string]interface{} {
	return map[string]interface{}{
		"id":   e.ID,
		"name": e.Name,
	}
}

// testChildEntity does NOT implement Auditable (simulates ReservationRoom)
// This is critical for reproducing the bug.
type testChildEntity struct {
	ID       uint   `gorm:"primaryKey"`
	ParentID uint   `gorm:"not null"`
	Name     string `gorm:"size:100"`
}

func (testChildEntity) TableName() string {
	return "test_child_entities"
}

// mockAuditService tracks all LogCreate calls to detect duplicates
type mockAuditService struct {
	mu          sync.Mutex
	createCalls []Auditable
}

func (m *mockAuditService) LogCreate(ctx context.Context, entity Auditable) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createCalls = append(m.createCalls, entity)
	return nil
}

func (m *mockAuditService) LogUpdate(ctx context.Context, entity Auditable, oldValues map[string]interface{}) error {
	return nil
}

func (m *mockAuditService) LogDelete(ctx context.Context, entity Auditable) error {
	return nil
}

func (m *mockAuditService) GetHistory(ctx context.Context, entityType string, entityID uint, page, size int) ([]AuditLog, int64, error) {
	return nil, 0, nil
}

func (m *mockAuditService) GetAllHistory(ctx context.Context, filter AuditLogFilter, page, size int) ([]AuditLog, int64, error) {
	return nil, 0, nil
}

func (m *mockAuditService) GetByID(ctx context.Context, id uint) (*AuditLog, error) {
	return nil, nil
}

func (m *mockAuditService) getCreateCallsByEntityType(entityType string) []Auditable {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []Auditable
	for _, call := range m.createCalls {
		if call.GetAuditEntityType() == entityType {
			result = append(result, call)
		}
	}
	return result
}

func setupTestDBWithMockHooks(t *testing.T, mockService *mockAuditService) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate test tables
	err = db.AutoMigrate(&testParentEntity{}, &testChildEntity{})
	require.NoError(t, err)

	// Register hooks with mock audit service
	RegisterHooks(db, mockService)

	return db
}

func TestAfterCreate_FullSaveAssociations_기존엔티티업데이트시_중복CREATE방지(t *testing.T) {
	// Setup with mock service to track LogCreate calls
	mockService := &mockAuditService{}
	db := setupTestDBWithMockHooks(t, mockService)
	ctx := SetUserContext(context.Background(), &[]uint{123}[0], "testuser")

	// Given: 새로운 부모 엔티티와 자식들을 한 번에 생성
	// testChildEntity does NOT implement Auditable, so:
	// 1. Parent's BeforeCreate sets db.Set("audit_entity", parent)
	// 2. Parent's AfterCreate calls LogCreate(parent) ✓
	// 3. GORM creates children (FullSaveAssociations)
	// 4. Child's BeforeCreate does nothing (not Auditable) - parent still in db.Set!
	// 5. Child's AfterCreate reads db.Get("audit_entity") → gets PARENT
	// 6. BUG: Calls LogCreate(parent) again for each child!
	parent := &testParentEntity{
		Name: "Parent",
		Childs: []testChildEntity{
			{Name: "Child 1"},
			{Name: "Child 2"},
		},
	}

	// When: FullSaveAssociations로 부모와 자식을 함께 저장
	err := db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(parent).Error
	require.NoError(t, err)

	// Verify children were created
	var childCount int64
	db.Table("test_child_entities").Count(&childCount)
	require.Equal(t, int64(2), childCount, "자식 엔티티가 2개 생성되어야 함")

	// Then: LogCreate should be called exactly once for parent
	// BUG: With buggy hooks, LogCreate is called 3 times for parent:
	// - Once from parent's AfterCreate
	// - Once from each child's AfterCreate (reading stale audit_entity)
	parentCreateCalls := mockService.getCreateCallsByEntityType("test_parent")

	t.Logf("LogCreate was called %d times for test_parent (expected 1)", len(parentCreateCalls))
	for i, call := range parentCreateCalls {
		t.Logf("  Call %d: EntityType=%s, EntityID=%d", i, call.GetAuditEntityType(), call.GetAuditEntityID())
	}

	// Total create calls for debugging
	t.Logf("Total LogCreate calls: %d", len(mockService.createCalls))
	for i, call := range mockService.createCalls {
		t.Logf("  Total Call %d: EntityType=%s, EntityID=%d", i, call.GetAuditEntityType(), call.GetAuditEntityID())
	}

	// This assertion should FAIL with buggy code
	// Expected: 1 (parent created once)
	// Actual with bug: 3 (parent + 2 children's AfterCreate each call LogCreate for parent)
	assert.Equal(t, 1, len(parentCreateCalls), "LogCreate(parent)가 1번만 호출되어야 함. 실제: %d번 (중복 호출 버그)", len(parentCreateCalls))
}

func TestAfterCreate_신규엔티티생성시_CREATE로그정상생성(t *testing.T) {
	// Setup with mock service
	mockService := &mockAuditService{}
	db := setupTestDBWithMockHooks(t, mockService)
	ctx := SetUserContext(context.Background(), &[]uint{123}[0], "testuser")

	// Given: 새로운 부모 엔티티 생성 (자식 없음)
	parent := &testParentEntity{
		Name: "New Parent",
	}

	// When: 정상적인 CREATE 작업 수행
	err := db.WithContext(ctx).Create(parent).Error
	require.NoError(t, err)

	// Then: LogCreate should be called exactly once for parent
	parentCreateCalls := mockService.getCreateCallsByEntityType("test_parent")

	t.Logf("LogCreate was called %d times for test_parent", len(parentCreateCalls))

	assert.Equal(t, 1, len(parentCreateCalls), "LogCreate가 1번 호출되어야 함")
	if len(parentCreateCalls) > 0 {
		assert.Equal(t, "test_parent", parentCreateCalls[0].GetAuditEntityType())
		assert.Equal(t, parent.ID, parentCreateCalls[0].GetAuditEntityID())
	}
}
