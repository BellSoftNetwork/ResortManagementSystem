package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_Liveness(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Liveness는 항상 UP 반환
	handler := NewHealthHandler(nil, nil)

	router := gin.New()
	router.GET("/actuator/health/liveness", handler.Liveness)

	req := httptest.NewRequest("GET", "/actuator/health/liveness", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 검증
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "UP", response["status"])
}

// 기본적인 응답 구조만 테스트
func TestHealthHandler_ResponseStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		endpoint string
		handler  func(*gin.Context)
	}{
		{
			name:     "Health endpoint returns proper structure",
			endpoint: "/actuator/health",
			handler: func(c *gin.Context) {
				handler := NewHealthHandler(nil, nil)
				handler.Health(c)
			},
		},
		{
			name:     "Readiness endpoint returns proper structure",
			endpoint: "/actuator/health/readiness",
			handler: func(c *gin.Context) {
				handler := NewHealthHandler(nil, nil)
				handler.Readiness(c)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET(tt.endpoint, tt.handler)

			req := httptest.NewRequest("GET", tt.endpoint, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 응답이 JSON인지 확인
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// status 필드가 있는지 확인
			_, hasStatus := response["status"]
			assert.True(t, hasStatus, "Response should have 'status' field")

			// components 필드가 있는지 확인
			_, hasComponents := response["components"]
			assert.True(t, hasComponents, "Response should have 'components' field")
		})
	}
}
