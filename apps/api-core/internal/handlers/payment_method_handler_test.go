package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

// ptr returns a pointer to the given value
func ptr[T any](v T) *T {
	return &v
}

// MockPaymentMethodService는 PaymentMethodService의 모킹 구현
type MockPaymentMethodService struct {
	mock.Mock
}

func (m *MockPaymentMethodService) GetByID(ctx context.Context, id uint) (*models.PaymentMethod, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodService) GetAll(ctx context.Context, page, size int, sort string) ([]models.PaymentMethod, int64, error) {
	args := m.Called(ctx, page, size, sort)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.PaymentMethod), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentMethodService) GetActive(ctx context.Context) ([]models.PaymentMethod, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodService) Create(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	args := m.Called(ctx, paymentMethod)
	return args.Error(0)
}

func (m *MockPaymentMethodService) Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.PaymentMethod, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPaymentMethodHandler_ListPaymentMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		expectedPage   int
		expectedSize   int
		expectedSort   string
		mockMethods    []models.PaymentMethod
		mockTotal      int64
		mockError      error
		expectedStatus int
		expectError    bool
	}{
		{
			name:         "정상적인 결제 방법 목록 조회",
			queryParams:  "page=0&size=20",
			expectedPage: 0,
			expectedSize: 20,
			expectedSort: "",
			mockMethods: []models.PaymentMethod{
				{
					BaseTimeEntity:           models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
					Name:                     "신용카드",
					CommissionRate:           0.025,
					RequireUnpaidAmountCheck: models.BitBool(false),
					IsDefaultSelect:          models.BitBool(true),
					Status:                   models.PaymentMethodStatusActive,
				},
				{
					BaseTimeEntity:           models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}},
					Name:                     "현금",
					CommissionRate:           0,
					RequireUnpaidAmountCheck: models.BitBool(true),
					IsDefaultSelect:          models.BitBool(false),
					Status:                   models.PaymentMethodStatusActive,
				},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:         "정렬이 적용된 목록 조회",
			queryParams:  "page=0&size=10&sort=name,asc",
			expectedPage: 0,
			expectedSize: 10,
			expectedSort: "name,asc",
			mockMethods: []models.PaymentMethod{
				{
					BaseTimeEntity:           models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 3}},
					Name:                     "계좌이체",
					CommissionRate:           0.015,
					RequireUnpaidAmountCheck: models.BitBool(false),
					IsDefaultSelect:          models.BitBool(false),
					Status:                   models.PaymentMethodStatusActive,
				},
				{
					BaseTimeEntity:           models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}},
					Name:                     "현금",
					CommissionRate:           0,
					RequireUnpaidAmountCheck: models.BitBool(true),
					IsDefaultSelect:          models.BitBool(false),
					Status:                   models.PaymentMethodStatusActive,
				},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "서비스 에러 발생",
			queryParams:    "page=0&size=20",
			expectedPage:   0,
			expectedSize:   20,
			expectedSort:   "",
			mockMethods:    nil,
			mockTotal:      0,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:           "빈 결과",
			queryParams:    "page=10&size=20",
			expectedPage:   10,
			expectedSize:   20,
			expectedSort:   "",
			mockMethods:    []models.PaymentMethod{},
			mockTotal:      0,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPaymentMethodService)
			handler := NewPaymentMethodHandler(mockService)

			// Set up mock expectations
			mockService.On("GetAll",
				mock.Anything,
				tt.expectedPage,
				tt.expectedSize,
				tt.expectedSort,
			).Return(tt.mockMethods, tt.mockTotal, tt.mockError)

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/payment-methods", handler.ListPaymentMethods)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/payment-methods?%s", tt.queryParams), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError {
				// Parse response
				var response struct {
					Values []dto.PaymentMethodResponse `json:"values"`
					Page   map[string]interface{}      `json:"page"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockMethods), len(response.Values))
				assert.Equal(t, tt.mockTotal, int64(response.Page["totalElements"].(float64)))
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestPaymentMethodHandler_GetPaymentMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		methodID       string
		mockMethod     *models.PaymentMethod
		mockError      error
		expectedStatus int
	}{
		{
			name:     "정상적인 결제 방법 조회",
			methodID: "1",
			mockMethod: &models.PaymentMethod{
				BaseTimeEntity:           models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
				Name:                     "신용카드",
				CommissionRate:           0.025,
				RequireUnpaidAmountCheck: models.BitBool(false),
				IsDefaultSelect:          models.BitBool(true),
				Status:                   models.PaymentMethodStatusActive,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "존재하지 않는 결제 방법",
			methodID:       "999",
			mockMethod:     nil,
			mockError:      services.ErrPaymentMethodNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "잘못된 ID 형식",
			methodID:       "invalid",
			mockMethod:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPaymentMethodService)
			handler := NewPaymentMethodHandler(mockService)

			// Set up mock expectations if valid ID
			if tt.methodID != "invalid" {
				methodID := uint(1)
				if tt.methodID == "999" {
					methodID = uint(999)
				}
				mockService.On("GetByID", mock.Anything, methodID).Return(tt.mockMethod, tt.mockError)
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/payment-methods/:id", handler.GetPaymentMethod)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/payment-methods/%s", tt.methodID), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var responseWrapper struct {
					Value dto.PaymentMethodResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &responseWrapper)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockMethod.Name, responseWrapper.Value.Name)
				assert.Equal(t, tt.mockMethod.CommissionRate, responseWrapper.Value.CommissionRate)
			}

			// Verify mock expectations
			if tt.methodID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestPaymentMethodHandler_CreatePaymentMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockPaymentMethodService)
		expectedStatus int
	}{
		{
			name: "정상적인 결제 방법 생성",
			requestBody: dto.CreatePaymentMethodRequest{
				Name:                     "간편결제",
				CommissionRate:           0.03,
				RequireUnpaidAmountCheck: false,
			},
			setupMock: func(m *MockPaymentMethodService) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(method *models.PaymentMethod) bool {
					return method.Name == "간편결제" && method.CommissionRate == 0.03
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "중복된 결제 방법명",
			requestBody: dto.CreatePaymentMethodRequest{
				Name:                     "신용카드",
				CommissionRate:           0.025,
				RequireUnpaidAmountCheck: false,
			},
			setupMock: func(m *MockPaymentMethodService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.PaymentMethod")).
					Return(services.ErrPaymentMethodNameExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "잘못된 요청 데이터",
			requestBody:    "invalid json",
			setupMock:      func(m *MockPaymentMethodService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "필수 필드 누락",
			requestBody: dto.CreatePaymentMethodRequest{
				// Name 필드 누락
				CommissionRate:           0,
				RequireUnpaidAmountCheck: true,
			},
			setupMock:      func(m *MockPaymentMethodService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "잘못된 수수료율",
			requestBody: dto.CreatePaymentMethodRequest{
				Name:                     "잘못된 수수료",
				CommissionRate:           -1, // 음수 수수료율
				RequireUnpaidAmountCheck: false,
			},
			setupMock:      func(m *MockPaymentMethodService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPaymentMethodService)
			handler := NewPaymentMethodHandler(mockService)

			// Set up mock expectations
			tt.setupMock(mockService)

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.POST("/api/v1/payment-methods", handler.CreatePaymentMethod)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/payment-methods", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				// Parse response
				var responseWrapper struct {
					Value dto.PaymentMethodResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &responseWrapper)
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestPaymentMethodHandler_UpdatePaymentMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		methodID       string
		requestBody    interface{}
		mockMethod     *models.PaymentMethod
		mockError      error
		expectedStatus int
	}{
		{
			name:     "정상적인 결제 방법 업데이트",
			methodID: "1",
			requestBody: dto.UpdatePaymentMethodRequest{
				Name:           ptr("업데이트된 신용카드"),
				CommissionRate: ptr(0.035),
			},
			mockMethod: &models.PaymentMethod{
				BaseTimeEntity:           models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
				Name:                     "업데이트된 신용카드",
				CommissionRate:           0.035,
				RequireUnpaidAmountCheck: models.BitBool(false),
				IsDefaultSelect:          models.BitBool(true),
				Status:                   models.PaymentMethodStatusActive,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:     "활성화 상태만 변경",
			methodID: "2",
			requestBody: dto.UpdatePaymentMethodRequest{
				CommissionRate: ptr(0.01),
			},
			mockMethod: &models.PaymentMethod{
				BaseTimeEntity:           models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}},
				Name:                     "현금",
				CommissionRate:           0.01,
				RequireUnpaidAmountCheck: models.BitBool(true),
				IsDefaultSelect:          models.BitBool(false),
				Status:                   models.PaymentMethodStatusActive,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:     "존재하지 않는 결제 방법",
			methodID: "999",
			requestBody: dto.UpdatePaymentMethodRequest{
				Name: ptr("업데이트 시도"),
			},
			mockMethod:     nil,
			mockError:      services.ErrPaymentMethodNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "잘못된 ID 형식",
			methodID:       "invalid",
			requestBody:    dto.UpdatePaymentMethodRequest{},
			mockMethod:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPaymentMethodService)
			handler := NewPaymentMethodHandler(mockService)

			// Set up mock expectations if valid ID
			if tt.methodID != "invalid" {
				methodID := uint(1)
				if tt.methodID == "999" {
					methodID = uint(999)
				} else if tt.methodID == "2" {
					methodID = uint(2)
				}
				mockService.On("Update", mock.Anything, methodID, mock.AnythingOfType("map[string]interface {}")).
					Return(tt.mockMethod, tt.mockError)
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.PATCH("/api/v1/payment-methods/:id", handler.UpdatePaymentMethod)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/payment-methods/%s", tt.methodID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var responseWrapper struct {
					Value dto.PaymentMethodResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &responseWrapper)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockMethod.Name, responseWrapper.Value.Name)
				assert.Equal(t, tt.mockMethod.CommissionRate, responseWrapper.Value.CommissionRate)
			}

			// Verify mock expectations
			if tt.methodID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestPaymentMethodHandler_DeletePaymentMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		methodID       string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "정상적인 결제 방법 삭제",
			methodID:       "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "존재하지 않는 결제 방법",
			methodID:       "999",
			mockError:      services.ErrPaymentMethodNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "사용 중인 결제 방법",
			methodID:       "2",
			mockError:      services.ErrPaymentMethodInUse,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "잘못된 ID 형식",
			methodID:       "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPaymentMethodService)
			handler := NewPaymentMethodHandler(mockService)

			// Set up mock expectations if valid ID
			if tt.methodID != "invalid" {
				methodID := uint(1)
				if tt.methodID == "999" {
					methodID = uint(999)
				} else if tt.methodID == "2" {
					methodID = uint(2)
				}
				mockService.On("Delete", mock.Anything, methodID).Return(tt.mockError)
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.DELETE("/api/v1/payment-methods/:id", handler.DeletePaymentMethod)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/payment-methods/%s", tt.methodID), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mock expectations
			if tt.methodID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}
