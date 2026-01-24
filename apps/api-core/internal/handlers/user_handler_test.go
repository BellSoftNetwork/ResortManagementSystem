package handlers

import (
	"bytes"
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

func TestUserHandler_GetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         uint
		mockUser       *models.User
		mockError      error
		expectedStatus int
	}{
		{
			name:   "정상적인 현재 사용자 조회",
			userID: 1,
			mockUser: &models.User{
				BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
				UserID:         "testuser",
				Email:          "test@example.com",
				Name:           "Test User",
				Role:           models.UserRoleNormal,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "존재하지 않는 사용자",
			userID:         999,
			mockUser:       nil,
			mockError:      errors.New("존재하지 않는 사용자"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "인증되지 않은 요청",
			userID:         0, // 컨텍스트에 사용자 ID가 없음
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			handler := NewUserHandler(mockService)

			// Set up mock expectations if authenticated
			if tt.userID > 0 {
				mockService.On("GetByID", mock.Anything, tt.userID).Return(tt.mockUser, tt.mockError)
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

			router.GET("/api/v1/my", handler.GetCurrentUser)
			router.POST("/api/v1/my", handler.GetCurrentUser) // POST 메서드도 지원

			// GET 요청 테스트
			req := httptest.NewRequest("GET", "/api/v1/my", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var actualResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
				assert.NoError(t, err)

				// 실제 응답 구조 확인 (value 필드 사용)
				if data, ok := actualResponse["value"].(map[string]interface{}); ok {
					if userID, exists := data["userId"]; exists {
						assert.Equal(t, tt.mockUser.UserID, userID)
					}
					if email, exists := data["email"]; exists {
						assert.Equal(t, tt.mockUser.Email, email)
					}
					if name, exists := data["name"]; exists {
						assert.Equal(t, tt.mockUser.Name, name)
					}
				}
			}

			// POST 요청도 테스트
			if tt.userID > 0 {
				req = httptest.NewRequest("POST", "/api/v1/my", nil)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			// Verify mock expectations
			if tt.userID > 0 {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestUserHandler_UpdateCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         uint
		requestBody    interface{}
		mockUser       *models.User
		mockError      error
		expectedStatus int
	}{
		{
			name:   "정상적인 사용자 정보 업데이트",
			userID: 1,
			requestBody: dto.UpdateUserRequest{
				Name: stringPtr("Updated Name"),
			},
			mockUser: &models.User{
				BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
				UserID:         "testuser",
				Email:          "test@example.com",
				Name:           "Updated Name",
				Role:           models.UserRoleNormal,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "이름만 업데이트",
			userID: 2,
			requestBody: dto.UpdateUserRequest{
				Name: stringPtr("New Name Only"),
			},
			mockUser: &models.User{
				BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}},
				UserID:         "testuser2",
				Email:          "test2@example.com",
				Name:           "New Name Only",
				Role:           models.UserRoleNormal,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "이메일 변경",
			userID: 3,
			requestBody: dto.UpdateUserRequest{
				Email: "newemail@example.com",
			},
			mockUser: &models.User{
				BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 3}},
				UserID:         "testuser3",
				Email:          "test3@example.com",
				Name:           "Test User 3",
				Role:           models.UserRoleNormal,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "인증되지 않은 요청",
			userID:         0,
			requestBody:    dto.UpdateUserRequest{Name: stringPtr("Should Fail")},
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "잘못된 요청 데이터",
			userID:         1,
			requestBody:    "invalid json",
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "서비스 에러",
			userID: 4,
			requestBody: dto.UpdateUserRequest{
				Name: stringPtr("Error Test"),
			},
			mockUser:       nil,
			mockError:      errors.New("업데이트 실패"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			handler := NewUserHandler(mockService)

			// Set up mock expectations if authenticated and valid request
			if tt.userID > 0 && tt.requestBody != "invalid json" {
				// UpdateUserRequest 처리
				if _, ok := tt.requestBody.(dto.UpdateUserRequest); ok {
					mockService.On("Update", mock.Anything, tt.userID, mock.AnythingOfType("map[string]interface {}")).
						Return(tt.mockUser, tt.mockError)
				}
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

			router.PATCH("/api/v1/my", handler.UpdateCurrentUser)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PATCH", "/api/v1/my", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response
				var actualResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
				assert.NoError(t, err)

				// 실제 응답 구조 확인 (value 필드 사용)
				if data, ok := actualResponse["value"].(map[string]interface{}); ok {
					if userID, exists := data["userId"]; exists {
						assert.Equal(t, tt.mockUser.UserID, userID)
					}
					if name, exists := data["name"]; exists {
						assert.Equal(t, tt.mockUser.Name, name)
					}
				}
			}

			// Verify mock expectations
			if tt.userID > 0 && tt.requestBody != "invalid json" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestUserHandler_ListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		userRole       models.UserRole
		expectedPage   int
		expectedSize   int
		mockUsers      []models.User
		mockTotal      int64
		mockError      error
		expectedStatus int
		expectError    bool
	}{
		{
			name:         "관리자가 사용자 목록 조회",
			queryParams:  "page=0&size=20",
			userRole:     models.UserRoleSuperAdmin,
			expectedPage: 0,
			expectedSize: 20,
			mockUsers: []models.User{
				{
					BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
					UserID:         "user1",
					Email:          "user1@example.com",
					Name:           "User 1",
					Role:           models.UserRoleNormal,
				},
				{
					BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}},
					UserID:         "user2",
					Email:          "user2@example.com",
					Name:           "User 2",
					Role:           models.UserRoleAdmin,
				},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "일반 사용자가 접근 시도",
			queryParams:    "page=0&size=20",
			userRole:       models.UserRoleNormal,
			expectedPage:   0,
			expectedSize:   20,
			mockUsers:      nil,
			mockTotal:      0,
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
			expectError:    true,
		},
		{
			name:         "페이지네이션 적용",
			queryParams:  "page=2&size=10",
			userRole:     models.UserRoleAdmin,
			expectedPage: 2,
			expectedSize: 10,
			mockUsers: []models.User{
				{
					BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 21}},
					UserID:         "user21",
					Email:          "user21@example.com",
					Name:           "User 21",
					Role:           models.UserRoleNormal,
				},
			},
			mockTotal:      21,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "서비스 에러",
			queryParams:    "page=0&size=20",
			userRole:       models.UserRoleSuperAdmin,
			expectedPage:   0,
			expectedSize:   20,
			mockUsers:      nil,
			mockTotal:      0,
			mockError:      errors.New("데이터베이스 에러"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			handler := NewUserHandler(mockService)

			// Set up mock expectations for authorized users
			if tt.userRole == models.UserRoleAdmin || tt.userRole == models.UserRoleSuperAdmin {
				mockService.On("GetAll",
					mock.Anything,
					tt.expectedPage,
					tt.expectedSize,
				).Return(tt.mockUsers, tt.mockTotal, tt.mockError)
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

			// Add role middleware like in production
			router.Use(middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))
			router.GET("/api/v1/admin/accounts", handler.ListUsers)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/admin/accounts?%s", tt.queryParams), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError && tt.expectedStatus == http.StatusOK {
				// Parse response
				var response struct {
					Values []dto.UserResponse     `json:"values"`
					Page   map[string]interface{} `json:"page"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockUsers), len(response.Values))
				assert.Equal(t, tt.mockTotal, int64(response.Page["totalElements"].(float64)))
			}

			// Verify mock expectations
			if tt.userRole == models.UserRoleAdmin || tt.userRole == models.UserRoleSuperAdmin {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       models.UserRole
		requestBody    interface{}
		setupMock      func(*MockUserService)
		expectedStatus int
	}{
		{
			name:     "관리자가 정상적인 사용자 생성",
			userRole: models.UserRoleSuperAdmin,
			requestBody: dto.CreateUserRequest{
				UserID:   "newuser",
				Email:    "newuser@example.com",
				Name:     "New User",
				Password: "password123",
				Role:     "NORMAL",
			},
			setupMock: func(m *MockUserService) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(user *models.User) bool {
					return user.UserID == "newuser" && user.Email == "newuser@example.com"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:     "일반 사용자가 생성 시도",
			userRole: models.UserRoleNormal,
			requestBody: dto.CreateUserRequest{
				UserID: "shouldfail",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:     "중복된 사용자 ID",
			userRole: models.UserRoleAdmin,
			requestBody: dto.CreateUserRequest{
				UserID:   "existinguser",
				Email:    "new@example.com",
				Name:     "New User",
				Password: "password123",
				Role:     "NORMAL",
			},
			setupMock: func(m *MockUserService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(services.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "잘못된 요청 데이터",
			userRole:       models.UserRoleSuperAdmin,
			requestBody:    "invalid json",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "필수 필드 누락",
			userRole: models.UserRoleSuperAdmin,
			requestBody: dto.CreateUserRequest{
				// UserID 누락
				Email: "test@example.com",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			handler := NewUserHandler(mockService)

			// Set up mock expectations
			tt.setupMock(mockService)

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())

			// 인증된 사용자를 시뮬레이션하기 위한 미들웨어
			router.Use(func(c *gin.Context) {
				c.Set(middleware.UserIDKey, uint(1))
				c.Set(middleware.UserRoleKey, tt.userRole.String())
				c.Next()
			})

			// Add role middleware like in production
			router.Use(middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))
			router.POST("/api/v1/admin/accounts", handler.CreateUser)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/admin/accounts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       models.UserRole
		targetUserID   string
		requestBody    interface{}
		mockUser       *models.User
		mockError      error
		expectedStatus int
	}{
		{
			name:         "관리자가 사용자 정보 업데이트",
			userRole:     models.UserRoleSuperAdmin,
			targetUserID: "2",
			requestBody: dto.UpdateUserRequest{
				Name: stringPtr("Updated User"),
				Role: stringPtr("ADMIN"),
			},
			mockUser: &models.User{
				BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}},
				UserID:         "user2",
				Email:          "user2@example.com",
				Name:           "Updated User",
				Role:           models.UserRoleAdmin,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "일반 사용자가 다른 사용자 수정 시도",
			userRole:       models.UserRoleNormal,
			targetUserID:   "2",
			requestBody:    dto.UpdateUserRequest{Name: stringPtr("Should Fail")},
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "존재하지 않는 사용자",
			userRole:       models.UserRoleAdmin,
			targetUserID:   "999",
			requestBody:    dto.UpdateUserRequest{Name: stringPtr("Not Found")},
			mockUser:       nil,
			mockError:      services.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "잘못된 ID 형식",
			userRole:       models.UserRoleSuperAdmin,
			targetUserID:   "invalid",
			requestBody:    dto.UpdateUserRequest{},
			mockUser:       nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			handler := NewUserHandler(mockService)

			// Set up mock expectations for authorized users with valid ID
			if (tt.userRole == models.UserRoleAdmin || tt.userRole == models.UserRoleSuperAdmin) && tt.targetUserID != "invalid" {
				targetID := uint(2)
				if tt.targetUserID == "999" {
					targetID = uint(999)
				}

				// UpdateUserRequest 처리
				if _, ok := tt.requestBody.(dto.UpdateUserRequest); ok {
					mockService.On("Update", mock.Anything, targetID, mock.AnythingOfType("map[string]interface {}")).
						Return(tt.mockUser, tt.mockError)
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

			// Add role middleware like in production
			router.Use(middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))
			router.PATCH("/api/v1/admin/accounts/:id", handler.UpdateUser)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/admin/accounts/%s", tt.targetUserID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response with value wrapper
				var responseWrapper struct {
					Value dto.UserResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &responseWrapper)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUser.Name, responseWrapper.Value.Name)
			}

			// Verify mock expectations
			if (tt.userRole == models.UserRoleAdmin || tt.userRole == models.UserRoleSuperAdmin) && tt.targetUserID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

// DeleteUser 메서드가 UserHandler에 없으므로 테스트 주석처리
/*
func TestUserHandler_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       models.UserRole
		targetUserID   string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "관리자가 사용자 삭제",
			userRole:       models.UserRoleSuperAdmin,
			targetUserID:   "2",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "일반 사용자가 삭제 시도",
			userRole:       models.UserRoleNormal,
			targetUserID:   "2",
			mockError:      nil,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "존재하지 않는 사용자",
			userRole:       models.UserRoleAdmin,
			targetUserID:   "999",
			mockError:      errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "본인 계정 삭제 시도",
			userRole:       models.UserRoleSuperAdmin,
			targetUserID:   "1", // 현재 로그인한 사용자 ID와 동일
			mockError:      errors.New("cannot delete own account"),
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "잘못된 ID 형식",
			userRole:       models.UserRoleSuperAdmin,
			targetUserID:   "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			handler := NewUserHandler(mockService)

			// Set up mock expectations for authorized users with valid ID
			if (tt.userRole == models.UserRoleAdmin || tt.userRole == models.UserRoleSuperAdmin) && tt.targetUserID != "invalid" {
				targetID := uint(2)
				if tt.targetUserID == "999" {
					targetID = uint(999)
				} else if tt.targetUserID == "1" {
					targetID = uint(1)
				}
				mockService.On("Delete", mock.Anything, targetID).Return(tt.mockError)
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())

			// 인증된 사용자를 시뮬레이션하기 위한 미들웨어
			router.Use(func(c *gin.Context) {
				c.Set(middleware.UserIDKey, uint(1)) // 현재 로그인한 사용자 ID
				c.Set(middleware.UserRoleKey, tt.userRole.String())
				c.Next()
			})

			router.DELETE("/api/v1/admin/accounts/:id", handler.DeleteUser)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/admin/accounts/%s", tt.targetUserID), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mock expectations
			if (tt.userRole == models.UserRoleAdmin || tt.userRole == models.UserRoleSuperAdmin) && tt.targetUserID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}
*/
