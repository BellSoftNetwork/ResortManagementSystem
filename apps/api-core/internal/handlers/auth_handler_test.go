package handlers

import (
	"bytes"
	"context"
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

// MockAuthService는 AuthService의 모킹 구현
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, userID, email, name, password string) (*models.User, error) {
	args := m.Called(ctx, userID, email, name, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, username, password, ipAddress string, deviceInfo *services.DeviceInfo) (*services.LoginResponse, error) {
	args := m.Called(ctx, username, password, ipAddress, deviceInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.LoginResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*services.TokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.TokenResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) IsDeviceChanged(ctx context.Context, username string, deviceInfo *services.DeviceInfo) (bool, error) {
	args := m.Called(ctx, username, deviceInfo)
	return args.Bool(0), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockUser       *models.User
		mockError      error
		expectedStatus int
	}{
		{
			name: "정상적인 회원가입",
			requestBody: dto.RegisterRequest{
				UserID:   "newuser",
				Email:    "newuser@example.com",
				Name:     "New User",
				Password: "password123",
			},
			mockUser: &models.User{
				BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
				UserID:         "newuser",
				Email:          "newuser@example.com",
				Name:           "New User",
				Role:           models.UserRoleNormal,
				Status:         models.UserStatusActive,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "중복된 사용자 ID",
			requestBody: dto.RegisterRequest{
				UserID:   "existinguser",
				Email:    "new@example.com",
				Name:     "New User",
				Password: "password123",
			},
			mockUser:       nil,
			mockError:      services.ErrUserAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name: "중복된 이메일",
			requestBody: dto.RegisterRequest{
				UserID:   "newuser2",
				Email:    "existing@example.com",
				Name:     "New User 2",
				Password: "password123",
			},
			mockUser:       nil,
			mockError:      services.ErrUserAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "잘못된 요청 데이터",
			requestBody:    "invalid json",
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "필수 필드 누락",
			requestBody: dto.RegisterRequest{
				// UserID 누락
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			},
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockAuthService)
			handler := NewAuthHandler(mockService)

			// Set up mock expectations for valid requests
			if tt.requestBody != "invalid json" {
				if req, ok := tt.requestBody.(dto.RegisterRequest); ok && req.UserID != "" {
					mockService.On("Register", mock.Anything, req.UserID, req.Email, req.Name, req.Password).
						Return(tt.mockUser, tt.mockError)
				}
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.POST("/api/v1/auth/register", handler.Register)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			if w.Code != tt.expectedStatus {
				t.Logf("Response body: %s", w.Body.String())
			}
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				// Parse response with value wrapper
				t.Logf("Response body for successful case: %s", w.Body.String())
				var wrapper struct {
					Value dto.UserResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &wrapper)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUser.UserID, wrapper.Value.UserID)
				assert.Equal(t, tt.mockUser.Email, wrapper.Value.Email)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockResponse   *services.LoginResponse
		mockError      error
		expectedStatus int
	}{
		{
			name: "정상적인 로그인",
			requestBody: dto.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockResponse: &services.LoginResponse{
				User: &models.User{
					BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
					UserID:         "testuser",
					Email:          "test@example.com",
					Name:           "Test User",
					Role:           models.UserRoleNormal,
				},
				AccessToken:          "access-token",
				RefreshToken:         "refresh-token",
				AccessTokenExpiresIn: 900000, // 15 minutes in milliseconds
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "잘못된 비밀번호",
			requestBody: dto.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockResponse:   nil,
			mockError:      services.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "존재하지 않는 사용자",
			requestBody: dto.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			mockResponse:   nil,
			mockError:      services.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "계정 잠김",
			requestBody: dto.LoginRequest{
				Username: "lockeduser",
				Password: "password123",
			},
			mockResponse:   nil,
			mockError:      services.ErrTooManyAttempts,
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "잘못된 요청 데이터",
			requestBody:    "invalid json",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockAuthService)
			handler := NewAuthHandler(mockService)

			// Set up mock expectations for valid requests
			if tt.requestBody != "invalid json" {
				if req, ok := tt.requestBody.(dto.LoginRequest); ok {
					// IP address from httptest.NewRequest defaults to 192.0.2.1
					mockService.On("Login", mock.Anything, req.Username, req.Password, "192.0.2.1", mock.AnythingOfType("*services.DeviceInfo")).
						Return(tt.mockResponse, tt.mockError)
				}
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.POST("/api/v1/auth/login", handler.Login)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			if w.Code != tt.expectedStatus {
				t.Logf("Response body: %s", w.Body.String())
			}
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response with value wrapper
				t.Logf("Response body for successful case: %s", w.Body.String())
				var wrapper struct {
					Value dto.LoginResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &wrapper)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse.AccessToken, wrapper.Value.AccessToken)
				assert.Equal(t, tt.mockResponse.RefreshToken, wrapper.Value.RefreshToken)
				assert.Equal(t, tt.mockResponse.AccessTokenExpiresIn, wrapper.Value.AccessTokenExpiresIn)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockResponse   *services.TokenResponse
		mockError      error
		expectedStatus int
	}{
		{
			name: "정상적인 토큰 갱신",
			requestBody: dto.RefreshTokenRequest{
				RefreshToken: "valid-refresh-token",
			},
			mockResponse: &services.TokenResponse{
				AccessToken:          "new-access-token",
				RefreshToken:         "new-refresh-token",
				AccessTokenExpiresIn: 900000, // 15 minutes in milliseconds
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "만료된 리프레시 토큰",
			requestBody: dto.RefreshTokenRequest{
				RefreshToken: "expired-refresh-token",
			},
			mockResponse:   nil,
			mockError:      errors.New("refresh token expired"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "유효하지 않은 리프레시 토큰",
			requestBody: dto.RefreshTokenRequest{
				RefreshToken: "invalid-refresh-token",
			},
			mockResponse:   nil,
			mockError:      errors.New("invalid refresh token"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "잘못된 요청 데이터",
			requestBody:    "invalid json",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "빈 리프레시 토큰",
			requestBody: dto.RefreshTokenRequest{
				RefreshToken: "",
			},
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockAuthService)
			handler := NewAuthHandler(mockService)

			// Set up mock expectations for valid requests
			if tt.requestBody != "invalid json" {
				if req, ok := tt.requestBody.(dto.RefreshTokenRequest); ok && req.RefreshToken != "" {
					mockService.On("RefreshToken", mock.Anything, req.RefreshToken).
						Return(tt.mockResponse, tt.mockError)
				}
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.POST("/api/v1/auth/refresh", handler.RefreshToken)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response with value wrapper
				var wrapper struct {
					Value dto.TokenResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &wrapper)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse.AccessToken, wrapper.Value.AccessToken)
				assert.Equal(t, tt.mockResponse.RefreshToken, wrapper.Value.RefreshToken)
				assert.Equal(t, tt.mockResponse.AccessTokenExpiresIn, wrapper.Value.AccessTokenExpiresIn)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         uint
		mockError      error
		expectedStatus int
	}{
		{
			name:           "정상적인 로그아웃",
			userID:         1,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "서비스 에러 발생",
			userID:         2,
			mockError:      errors.New("redis error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "인증되지 않은 요청",
			userID:         0, // 컨텍스트에 사용자 ID가 없음
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockAuthService)
			handler := NewAuthHandler(mockService)

			// Set up mock expectations for authenticated requests
			if tt.userID > 0 {
				mockService.On("Logout", mock.Anything, tt.userID).Return(tt.mockError)
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())

			// 인증된 사용자를 시뮬레이션하기 위한 미들웨어
			if tt.userID > 0 {
				router.Use(func(c *gin.Context) {
					c.Set(middleware.UserIDKey, tt.userID)
					c.Next()
				})
			}

			router.POST("/api/v1/auth/logout", handler.Logout)

			req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response with value wrapper
				var wrapper struct {
					Value map[string]string `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &wrapper)
				assert.NoError(t, err)
				assert.Equal(t, "Successfully logged out", wrapper.Value["message"])
			}

			// Verify mock expectations
			if tt.userID > 0 {
				mockService.AssertExpectations(t)
			}
		})
	}
}
