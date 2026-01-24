package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

// MockConfigService는 ConfigService의 모킹 구현
type MockConfigService struct {
	mock.Mock
}

func (m *MockConfigService) GetConfig() *dto.ConfigResponse {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dto.ConfigResponse)
}

func (m *MockConfigService) GetEnvironment() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetEnvInfo() services.EnvInfo {
	args := m.Called()
	return args.Get(0).(services.EnvInfo)
}

func (m *MockConfigService) GetAppConfig(hasUsers bool) services.AppConfig {
	args := m.Called(hasUsers)
	return args.Get(0).(services.AppConfig)
}

func (m *MockConfigService) GetEnvironmentResponse() *dto.EnvironmentResponse {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dto.EnvironmentResponse)
}

func TestConfigHandler_GetConfig(t *testing.T) {
	t.Skip("ConfigHandler is not used, MainHandler handles /config and /env endpoints")
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		mockReturn       *dto.ConfigResponse
		expectedStatus   int
		expectedResponse *dto.ConfigResponse
	}{
		{
			name: "정상적인 설정 조회",
			mockReturn: &dto.ConfigResponse{
				API: dto.APIConfig{
					Host:    "localhost",
					Port:    8080,
					Profile: "local",
				},
				Database: dto.DatabaseConfig{
					Host:     "localhost",
					Port:     3306,
					Database: "resort_management",
				},
				Redis: dto.RedisConfig{
					Host: "localhost",
					Port: 6379,
				},
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &dto.ConfigResponse{
				API: dto.APIConfig{
					Host:    "localhost",
					Port:    8080,
					Profile: "local",
				},
				Database: dto.DatabaseConfig{
					Host:     "localhost",
					Port:     3306,
					Database: "resort_management",
				},
				Redis: dto.RedisConfig{
					Host: "localhost",
					Port: 6379,
				},
			},
		},
		{
			name:             "빈 설정 반환",
			mockReturn:       &dto.ConfigResponse{},
			expectedStatus:   http.StatusOK,
			expectedResponse: &dto.ConfigResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockConfigService)
			handler := NewConfigHandler(mockService)

			// Set up mock expectations
			mockService.On("GetConfig").Return(tt.mockReturn)

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/config", handler.GetConfig)

			req := httptest.NewRequest("GET", "/api/v1/config", nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response dto.ConfigResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResponse, &response)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestConfigHandler_GetEnvironment(t *testing.T) {
	t.Skip("ConfigHandler is not used, MainHandler handles /config and /env endpoints")
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		mockReturn       *dto.EnvironmentResponse
		expectedStatus   int
		expectedResponse *dto.EnvironmentResponse
	}{
		{
			name: "정상적인 환경 정보 조회",
			mockReturn: &dto.EnvironmentResponse{
				Profile:  "local",
				Hostname: "test-host",
				Version:  "1.0.0",
				Uptime:   "1h2m3s",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &dto.EnvironmentResponse{
				Profile:  "local",
				Hostname: "test-host",
				Version:  "1.0.0",
				Uptime:   "1h2m3s",
			},
		},
		{
			name: "프로덕션 환경 정보",
			mockReturn: &dto.EnvironmentResponse{
				Profile:  "production",
				Hostname: "prod-server-01",
				Version:  "2.5.0",
				Uptime:   "120h30m",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: &dto.EnvironmentResponse{
				Profile:  "production",
				Hostname: "prod-server-01",
				Version:  "2.5.0",
				Uptime:   "120h30m",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockConfigService)
			handler := NewConfigHandler(mockService)

			// Set up mock expectations
			mockService.On("GetEnvironmentResponse").Return(tt.mockReturn)

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/env", handler.GetEnvironment)

			req := httptest.NewRequest("GET", "/api/v1/env", nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response dto.EnvironmentResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResponse, &response)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
