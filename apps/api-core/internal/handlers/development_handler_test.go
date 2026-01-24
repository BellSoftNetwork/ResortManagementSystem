package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

// MockDevelopmentService는 DevelopmentService의 모킹 구현
type MockDevelopmentService struct {
	mock.Mock
}

func (m *MockDevelopmentService) GenerateTestData(dataType string, reservationOptions *services.ReservationGenerationOptions) (map[string]interface{}, error) {
	args := m.Called(dataType, reservationOptions)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestDevelopmentHandler_GenerateTestData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       models.UserRole
		apiProfile     string
		requestBody    interface{}
		mockResponse   map[string]interface{}
		mockError      error
		expectedStatus int
	}{
		{
			name:       "필수 데이터 생성",
			userRole:   models.UserRoleSuperAdmin,
			apiProfile: "local",
			requestBody: dto.GenerateTestDataRequest{
				Type: dto.TestDataTypeEssential,
			},
			mockResponse: map[string]interface{}{
				"paymentMethods": map[string]interface{}{
					"created": 4,
				},
				"roomGroups": map[string]interface{}{
					"created": 7,
				},
				"rooms": map[string]interface{}{
					"created": 50,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:       "예약 데이터 생성",
			userRole:   models.UserRoleSuperAdmin,
			apiProfile: "local",
			requestBody: dto.GenerateTestDataRequest{
				Type: dto.TestDataTypeReservation,
			},
			mockResponse: map[string]interface{}{
				"reservations": map[string]interface{}{
					"created":      18,
					"oneDayStay":   10,
					"twoDaysStay":  4,
					"monthlyRent":  4,
					"totalCreated": 18,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "잘못된 요청 데이터",
			userRole:       models.UserRoleSuperAdmin,
			apiProfile:     "local",
			requestBody:    "invalid json",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "서비스 에러 발생",
			userRole:   models.UserRoleSuperAdmin,
			apiProfile: "local",
			requestBody: dto.GenerateTestDataRequest{
				Type: dto.TestDataTypeAll,
			},
			mockResponse:   nil,
			mockError:      errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:       "일반 사용자가 접근 시도",
			userRole:   models.UserRoleNormal,
			apiProfile: "local",
			requestBody: dto.GenerateTestDataRequest{
				Type: dto.TestDataTypeEssential,
			},
			mockResponse: map[string]interface{}{
				"paymentMethods": 0,
				"roomGroups":     0,
				"rooms":          0,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK, // Note: Role middleware is not included in unit test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockDevelopmentService)
			handler := NewDevelopmentHandler(mockService)

			// Set up mock expectations for valid requests
			if tt.requestBody != "invalid json" {
				req, ok := tt.requestBody.(dto.GenerateTestDataRequest)
				if ok {
					mockService.On("GenerateTestData", string(req.Type), mock.Anything).
						Return(tt.mockResponse, tt.mockError)
				}
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())

			// 인증된 사용자를 시뮬레이션하기 위한 미들웨어
			router.Use(func(c *gin.Context) {
				c.Set(middleware.UserIDKey, uint(1))
				c.Set(middleware.UserRoleKey, tt.userRole.String())
				c.Next()
			})

			router.POST("/api/v1/dev/generate-test-data", handler.GenerateTestData)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/dev/generate-test-data", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var response dto.GenerateTestDataResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Test data generated successfully", response.Message)
				assert.NotNil(t, response.Data)
			}

			// Verify mock expectations
			if tt.requestBody != "invalid json" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

// TestDevelopmentHandler_MiddlewareProtection is skipped as middleware testing is done separately
func TestDevelopmentHandler_MiddlewareProtection(t *testing.T) {
	t.Skip("Middleware protection testing is done in middleware tests")
}

func TestDevelopmentHandler_GenerateEssentialData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       models.UserRole
		mockResponse   map[string]interface{}
		mockError      error
		expectedStatus int
	}{
		{
			name:     "필수 데이터 생성 성공",
			userRole: models.UserRoleSuperAdmin,
			mockResponse: map[string]interface{}{
				"paymentMethods": map[string]interface{}{
					"created": 4,
				},
				"roomGroups": map[string]interface{}{
					"created": 7,
				},
				"rooms": map[string]interface{}{
					"created": 50,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "서비스 에러 발생",
			userRole:       models.UserRoleSuperAdmin,
			mockResponse:   nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockDevelopmentService)
			handler := NewDevelopmentHandler(mockService)

			// Set up mock expectations
			mockService.On("GenerateTestData", "essential", (*services.ReservationGenerationOptions)(nil)).
				Return(tt.mockResponse, tt.mockError)

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.POST("/api/v1/dev/generate-essential-data", handler.GenerateEssentialData)

			req := httptest.NewRequest("POST", "/api/v1/dev/generate-essential-data", nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Essential data generated successfully", response["message"])
				assert.NotNil(t, response["data"])
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
