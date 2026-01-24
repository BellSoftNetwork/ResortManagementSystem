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
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

// stringPtr은 문자열에 대한 포인터를 반환하는 헬퍼 함수
func stringPtr(s string) *string {
	return &s
}

// MockRoomGroupService는 RoomGroupService의 모킹 구현
type MockRoomGroupService struct {
	mock.Mock
}

func (m *MockRoomGroupService) GetByID(ctx context.Context, id uint) (*models.RoomGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupService) GetByIDWithRooms(ctx context.Context, id uint, roomStatus *models.RoomStatus) (*models.RoomGroup, error) {
	args := m.Called(ctx, id, roomStatus)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupService) GetByIDWithFilteredRooms(ctx context.Context, id uint, filter repositories.RoomGroupRoomFilter) (*models.RoomGroup, error) {
	args := m.Called(ctx, id, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupService) GetByIDWithUsers(ctx context.Context, id uint) (*models.RoomGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupService) GetAllWithUsers(ctx context.Context, page, size int, sort string) ([]models.RoomGroup, int64, error) {
	args := m.Called(ctx, page, size, sort)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.RoomGroup), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomGroupService) GetAll(ctx context.Context, page, size int) ([]models.RoomGroup, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.RoomGroup), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomGroupService) Create(ctx context.Context, roomGroup *models.RoomGroup) error {
	args := m.Called(ctx, roomGroup)
	return args.Error(0)
}

func (m *MockRoomGroupService) Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.RoomGroup, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestRoomGroupHandler_ListRoomGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		expectedPage   int
		expectedSize   int
		mockGroups     []models.RoomGroup
		mockTotal      int64
		mockError      error
		expectedStatus int
		expectError    bool
	}{
		{
			name:         "정상적인 객실 그룹 목록 조회",
			queryParams:  "page=0&size=20",
			expectedPage: 0,
			expectedSize: 20,
			mockGroups: []models.RoomGroup{
				{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{
							BaseEntity: models.BaseEntity{ID: 1},
						},
					},
					Name:        "스탠다드 그룹",
					Description: "기본 객실 그룹",
				},
				{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{
							BaseEntity: models.BaseEntity{ID: 2},
						},
					},
					Name:        "디럭스 그룹",
					Description: "고급 객실 그룹",
				},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:         "페이지네이션 적용",
			queryParams:  "page=1&size=10",
			expectedPage: 1,
			expectedSize: 10,
			mockGroups: []models.RoomGroup{
				{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{
							BaseEntity: models.BaseEntity{ID: 11},
						},
					},
					Name:        "스위트 그룹",
					Description: "최고급 객실 그룹",
				},
			},
			mockTotal:      11,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "서비스 에러 발생",
			queryParams:    "page=0&size=20",
			expectedPage:   0,
			expectedSize:   20,
			mockGroups:     nil,
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
			mockGroups:     []models.RoomGroup{},
			mockTotal:      0,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockRoomGroupService)
			mockReservationService := new(MockReservationService)
			mockUserService := new(MockUserService)
			handler := NewRoomGroupHandler(mockService, mockReservationService, mockUserService)

			// Set up mock expectations
			mockService.On("GetAllWithUsers",
				mock.Anything,
				tt.expectedPage,
				tt.expectedSize,
				mock.AnythingOfType("string"),
			).Return(tt.mockGroups, tt.mockTotal, tt.mockError)

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/room-groups", handler.ListRoomGroups)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/room-groups?%s", tt.queryParams), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError {
				// Parse response
				var response struct {
					Values []dto.RoomGroupResponse `json:"values"`
					Page   map[string]interface{}  `json:"page"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockGroups), len(response.Values))
				assert.Equal(t, tt.mockTotal, int64(response.Page["totalElements"].(float64)))
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestRoomGroupHandler_GetRoomGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		groupID        string
		mockGroup      *models.RoomGroup
		mockError      error
		expectedStatus int
	}{
		{
			name:    "정상적인 객실 그룹 조회 (객실 포함)",
			groupID: "1",
			mockGroup: &models.RoomGroup{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					BaseTimeEntity: models.BaseTimeEntity{
						BaseEntity: models.BaseEntity{ID: 1},
					},
				},
				Name:        "스탠다드 그룹",
				Description: "기본 객실 그룹",
				Rooms: []models.Room{
					{
						BaseMustAuditEntity: models.BaseMustAuditEntity{
							BaseTimeEntity: models.BaseTimeEntity{
								BaseEntity: models.BaseEntity{ID: 1},
							},
						},
						Number:      "101",
						RoomGroupID: 1,
						Note:        "스탠다드 룸 101",
						Status:      models.RoomStatusNormal,
					},
					{
						BaseMustAuditEntity: models.BaseMustAuditEntity{
							BaseTimeEntity: models.BaseTimeEntity{
								BaseEntity: models.BaseEntity{ID: 2},
							},
						},
						Number:      "102",
						RoomGroupID: 1,
						Note:        "스탠다드 룸 102",
						Status:      models.RoomStatusNormal,
					},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "존재하지 않는 객실 그룹",
			groupID:        "999",
			mockGroup:      nil,
			mockError:      services.ErrRoomGroupNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "잘못된 ID 형식",
			groupID:        "invalid",
			mockGroup:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockRoomGroupService)
			mockReservationService := new(MockReservationService)
			mockUserService := new(MockUserService)
			handler := NewRoomGroupHandler(mockService, mockReservationService, mockUserService)

			// Set up mock expectations if valid ID
			if tt.groupID != "invalid" {
				groupID := uint(1)
				if tt.groupID == "999" {
					groupID = uint(999)
				}
				mockService.On("GetByIDWithFilteredRooms", mock.Anything, groupID, mock.AnythingOfType("repositories.RoomGroupRoomFilter")).Return(tt.mockGroup, tt.mockError)

				// Mock user service calls for room creation/update users
				if tt.mockGroup != nil && tt.mockError == nil {
					// Mock GetByID for any user ID (return error to trigger ID-only response)
					mockUserService.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(nil, errors.New("user not found"))

					// Mock GetLastReservationForRoom for each room
					for _, room := range tt.mockGroup.Rooms {
						mockReservationService.On("GetLastReservationForRoom", mock.Anything, room.ID).Return(nil, errors.New("no reservation"))
					}
				}
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/room-groups/:id", handler.GetRoomGroup)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/room-groups/%s", tt.groupID), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			if w.Code != tt.expectedStatus {
				t.Logf("Response body: %s", w.Body.String())
			}
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Always log response for debugging
			t.Logf("Response status: %d, body: %s", w.Code, w.Body.String())

			if tt.expectedStatus == http.StatusOK {
				// Parse response with value wrapper
				var wrapper struct {
					Value dto.RoomGroupResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &wrapper)
				assert.NoError(t, err)
				if err != nil {
					t.Logf("Failed to parse response: %v", err)
					t.Logf("Response body: %s", w.Body.String())
				}
				assert.Equal(t, tt.mockGroup.Name, wrapper.Value.Name)
				assert.Equal(t, len(tt.mockGroup.Rooms), len(wrapper.Value.Rooms))
			}

			// Verify mock expectations
			if tt.groupID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestRoomGroupHandler_CreateRoomGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockRoomGroupService)
		expectedStatus int
	}{
		{
			name: "정상적인 객실 그룹 생성",
			requestBody: dto.CreateRoomGroupRequest{
				Name:        "신규 스위트 그룹",
				Description: "새로운 스위트 객실 그룹",
			},
			setupMock: func(m *MockRoomGroupService) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(group *models.RoomGroup) bool {
					return group.Name == "신규 스위트 그룹" && group.Description == "새로운 스위트 객실 그룹"
				})).Return(nil).Run(func(args mock.Arguments) {
					// Set ID on creation
					group := args.Get(1).(*models.RoomGroup)
					group.ID = 3
				})

				// Mock the GetByIDWithUsers call after creation
				createdGroup := &models.RoomGroup{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{
							BaseEntity: models.BaseEntity{ID: 3},
						},
					},
					Name:        "신규 스위트 그룹",
					Description: "새로운 스위트 객실 그룹",
				}
				m.On("GetByIDWithUsers", mock.Anything, uint(3)).Return(createdGroup, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "중복된 그룹명",
			requestBody: dto.CreateRoomGroupRequest{
				Name:        "스탠다드 그룹",
				Description: "이미 존재하는 그룹명",
			},
			setupMock: func(m *MockRoomGroupService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.RoomGroup")).
					Return(services.ErrRoomGroupNameExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "잘못된 요청 데이터",
			requestBody:    "invalid json",
			setupMock:      func(m *MockRoomGroupService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "필수 필드 누락",
			requestBody: dto.CreateRoomGroupRequest{
				// Name 필드 누락
				Description: "설명만 있는 그룹",
			},
			setupMock:      func(m *MockRoomGroupService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockRoomGroupService)
			mockReservationService := new(MockReservationService)
			mockUserService := new(MockUserService)
			handler := NewRoomGroupHandler(mockService, mockReservationService, mockUserService)

			// Set up mock expectations
			tt.setupMock(mockService)

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			// Add middleware to set user context for authenticated endpoints
			router.Use(func(c *gin.Context) {
				c.Set(middleware.UserIDKey, uint(1))
				c.Set(middleware.UsernameKey, "testuser")
				c.Set(middleware.UserRoleKey, "ADMIN")
				c.Next()
			})
			router.POST("/api/v1/room-groups", handler.CreateRoomGroup)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/room-groups", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				// Parse response with value wrapper
				var wrapper struct {
					Value dto.RoomGroupResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &wrapper)
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestRoomGroupHandler_UpdateRoomGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		groupID        string
		requestBody    interface{}
		mockGroup      *models.RoomGroup
		mockError      error
		expectedStatus int
	}{
		{
			name:    "정상적인 객실 그룹 업데이트",
			groupID: "1",
			requestBody: dto.UpdateRoomGroupRequest{
				Name:        stringPtr("업데이트된 스탠다드 그룹"),
				Description: stringPtr("업데이트된 설명"),
			},
			mockGroup: &models.RoomGroup{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					BaseTimeEntity: models.BaseTimeEntity{
						BaseEntity: models.BaseEntity{ID: 1},
					},
				},
				Name:        "업데이트된 스탠다드 그룹",
				Description: "업데이트된 설명",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:    "부분 업데이트 (설명만)",
			groupID: "2",
			requestBody: dto.UpdateRoomGroupRequest{
				Description: stringPtr("새로운 설명만 업데이트"),
			},
			mockGroup: &models.RoomGroup{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					BaseTimeEntity: models.BaseTimeEntity{
						BaseEntity: models.BaseEntity{ID: 2},
					},
				},
				Name:        "디럭스 그룹",
				Description: "새로운 설명만 업데이트",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:    "존재하지 않는 객실 그룹",
			groupID: "999",
			requestBody: dto.UpdateRoomGroupRequest{
				Name: stringPtr("업데이트 시도"),
			},
			mockGroup:      nil,
			mockError:      services.ErrRoomGroupNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "잘못된 ID 형식",
			groupID:        "invalid",
			requestBody:    dto.UpdateRoomGroupRequest{},
			mockGroup:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockRoomGroupService)
			mockReservationService := new(MockReservationService)
			mockUserService := new(MockUserService)
			handler := NewRoomGroupHandler(mockService, mockReservationService, mockUserService)

			// Set up mock expectations if valid ID
			if tt.groupID != "invalid" {
				groupID := uint(1)
				if tt.groupID == "999" {
					groupID = uint(999)
				} else if tt.groupID == "2" {
					groupID = uint(2)
				}
				mockService.On("Update", mock.Anything, groupID, mock.AnythingOfType("map[string]interface {}")).
					Return(tt.mockGroup, tt.mockError)

				// Mock GetByIDWithUsers for successful updates
				if tt.mockGroup != nil && tt.mockError == nil {
					mockService.On("GetByIDWithUsers", mock.Anything, tt.mockGroup.ID).Return(tt.mockGroup, nil)
				}
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			// Add middleware to set user context for authenticated endpoints
			router.Use(func(c *gin.Context) {
				c.Set(middleware.UserIDKey, uint(1))
				c.Set(middleware.UsernameKey, "testuser")
				c.Set(middleware.UserRoleKey, "ADMIN")
				c.Next()
			})
			router.PATCH("/api/v1/room-groups/:id", handler.UpdateRoomGroup)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/room-groups/%s", tt.groupID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response with value wrapper
				var wrapper struct {
					Value dto.RoomGroupResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &wrapper)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockGroup.Name, wrapper.Value.Name)
				assert.Equal(t, tt.mockGroup.Description, wrapper.Value.Description)
			}

			// Verify mock expectations
			if tt.groupID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestRoomGroupHandler_DeleteRoomGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		groupID        string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "정상적인 객실 그룹 삭제",
			groupID:        "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "존재하지 않는 객실 그룹",
			groupID:        "999",
			mockError:      services.ErrRoomGroupNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "객실이 있는 그룹",
			groupID:        "2",
			mockError:      services.ErrRoomGroupHasRooms,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "잘못된 ID 형식",
			groupID:        "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockRoomGroupService)
			mockReservationService := new(MockReservationService)
			mockUserService := new(MockUserService)
			handler := NewRoomGroupHandler(mockService, mockReservationService, mockUserService)

			// Set up mock expectations if valid ID
			if tt.groupID != "invalid" {
				groupID := uint(1)
				if tt.groupID == "999" {
					groupID = uint(999)
				} else if tt.groupID == "2" {
					groupID = uint(2)
				}
				mockService.On("Delete", mock.Anything, groupID).Return(tt.mockError)
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			// Add middleware to set user context for authenticated endpoints
			router.Use(func(c *gin.Context) {
				c.Set(middleware.UserIDKey, uint(1))
				c.Set(middleware.UsernameKey, "testuser")
				c.Set(middleware.UserRoleKey, "ADMIN")
				c.Next()
			})
			router.DELETE("/api/v1/room-groups/:id", handler.DeleteRoomGroup)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/room-groups/%s", tt.groupID), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mock expectations
			if tt.groupID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}
