package dto

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuditLogQuery_EntityTypeOptional(t *testing.T) {
	// Given: 테스트 환경 설정
	gin.SetMode(gin.TestMode)

	t.Run("entityType이 없어도 validation 통과해야 한다", func(t *testing.T) {
		// Given: entityType이 없는 쿼리 스트링
		req := httptest.NewRequest("GET", "/api/v1/audit?startDate=2026-01-01&endDate=2026-01-31", nil)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = req

		// When: 쿼리 바인딩 실행
		var query AuditLogQuery
		err := c.ShouldBindQuery(&query)

		// Then: 에러가 없어야 하며 EntityType은 빈 문자열이어야 함
		assert.NoError(t, err)
		assert.Empty(t, query.EntityType)
		assert.Equal(t, "2026-01-01", query.StartDate)
		assert.Equal(t, "2026-01-31", query.EndDate)
	})

	t.Run("entityType이 있으면 정상 동작해야 한다", func(t *testing.T) {
		// Given: entityType이 있는 쿼리 스트링
		req := httptest.NewRequest("GET", "/api/v1/audit?entityType=room&startDate=2026-01-01", nil)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = req

		// When: 쿼리 바인딩 실행
		var query AuditLogQuery
		err := c.ShouldBindQuery(&query)

		// Then: 에러가 없어야 하며 EntityType은 "room"이어야 함
		assert.NoError(t, err)
		assert.Equal(t, "room", query.EntityType)
		assert.Equal(t, "2026-01-01", query.StartDate)
	})

	t.Run("모든 필드가 설정된 경우 정상 동작해야 한다", func(t *testing.T) {
		// Given: 모든 필드가 있는 쿼리 스트링
		req := httptest.NewRequest("GET", "/api/v1/audit?entityType=reservation&startDate=2026-01-01&endDate=2026-01-31&action=UPDATE&userId=1&entityId=100", nil)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = req

		// When: 쿼리 바인딩 실행
		var query AuditLogQuery
		err := c.ShouldBindQuery(&query)

		// Then: 모든 필드가 정상적으로 바인딩되어야 함
		assert.NoError(t, err)
		assert.Equal(t, "reservation", query.EntityType)
		assert.Equal(t, "2026-01-01", query.StartDate)
		assert.Equal(t, "2026-01-31", query.EndDate)
		assert.Equal(t, "UPDATE", query.Action)
		assert.NotNil(t, query.UserID)
		assert.Equal(t, uint(1), *query.UserID)
		assert.NotNil(t, query.EntityID)
		assert.Equal(t, uint(100), *query.EntityID)
	})
}
