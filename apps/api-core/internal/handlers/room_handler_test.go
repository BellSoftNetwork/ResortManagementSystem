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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

// MockHistoryService는 HistoryService의 모킹 구현
type MockHistoryService struct {
	mock.Mock
}

func (m *MockHistoryService) GetRoomHistory(ctx context.Context, roomID uint, page, size int) ([]dto.RoomRevisionResponse, int64, error) {
	args := m.Called(ctx, roomID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]dto.RoomRevisionResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockHistoryService) GetReservationHistory(ctx context.Context, reservationID uint, page, size int) ([]dto.ReservationRevisionResponse, int64, error) {
	args := m.Called(ctx, reservationID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]dto.ReservationRevisionResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockHistoryService) GetDateBlockHistory(ctx context.Context, dateBlockID uint, page, size int) ([]dto.DateBlockRevisionResponse, int64, error) {
	args := m.Called(ctx, dateBlockID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]dto.DateBlockRevisionResponse), args.Get(1).(int64), args.Error(2)
}

// MockRoomService는 RoomService의 모킹 구현
type MockRoomService struct {
	mock.Mock
}

func (m *MockRoomService) GetByID(ctx context.Context, id uint) (*models.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomService) GetByIDWithGroup(ctx context.Context, id uint) (*models.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomService) GetAll(ctx context.Context, filter dto.RoomRepositoryFilter, page, size int, sort string) ([]models.Room, int64, error) {
	args := m.Called(ctx, filter, page, size, sort)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Room), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomService) Create(ctx context.Context, room *models.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomService) Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.Room, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoomService) GetAvailableRooms(ctx context.Context, startDate, endDate time.Time, excludeReservationID *uint) ([]models.Room, error) {
	args := m.Called(ctx, startDate, endDate, excludeReservationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Room), args.Error(1)
}

func TestRoomHandler_ListRooms(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		expectedPage   int
		expectedSize   int
		expectedSort   string
		mockRooms      []models.Room
		mockTotal      int64
		mockError      error
		expectedStatus int
		expectError    bool
	}{
		{
			name:         "정상적인 객실 목록 조회",
			queryParams:  "page=0&size=20",
			expectedPage: 0,
			expectedSize: 20,
			expectedSort: "",
			mockRooms: []models.Room{
				{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{
							BaseEntity: models.BaseEntity{ID: 1},
						},
					},
					Number:      "101",
					RoomGroupID: 1,
					Note:        "스탠다드 룸",
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
					Note:        "디럭스 룸",
					Status:      models.RoomStatusNormal,
				},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:         "타입 필터링된 객실 목록",
			queryParams:  "page=0&size=10&type=DELUXE",
			expectedPage: 0,
			expectedSize: 10,
			expectedSort: "",
			mockRooms: []models.Room{
				{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{
							BaseEntity: models.BaseEntity{ID: 2},
						},
					},
					Number:      "102",
					RoomGroupID: 2,
					Note:        "디럭스 룸",
					Status:      models.RoomStatusNormal,
				},
			},
			mockTotal:      1,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:         "정렬이 적용된 객실 목록",
			queryParams:  "page=0&size=20&sort=number,asc",
			expectedPage: 0,
			expectedSize: 20,
			expectedSort: "number,asc",
			mockRooms: []models.Room{
				{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{
							BaseEntity: models.BaseEntity{ID: 1},
						},
					},
					Number:      "101",
					RoomGroupID: 1,
					Note:        "스탠다드 룸",
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
					Note:        "디럭스 룸",
					Status:      models.RoomStatusNormal,
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
			mockRooms:      nil,
			mockTotal:      0,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockRoomService)
			mockUserService := new(MockUserService)
			mockHistoryService := new(MockHistoryService)
			handler := NewRoomHandler(mockService, mockUserService, mockHistoryService)

			mockService.On("GetAll",
				mock.Anything,
				mock.AnythingOfType("dto.RoomRepositoryFilter"),
				tt.expectedPage,
				tt.expectedSize,
				tt.expectedSort,
			).Return(tt.mockRooms, tt.mockTotal, tt.mockError)

			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/rooms", handler.ListRooms)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/rooms?%s", tt.queryParams), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError {
				var response struct {
					Values []dto.RoomResponse     `json:"values"`
					Page   map[string]interface{} `json:"page"`
					Filter interface{}            `json:"filter"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, int(tt.mockTotal), len(response.Values))
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestRoomHandler_GetRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		roomID         string
		mockRoom       *models.Room
		mockError      error
		expectedStatus int
	}{
		{
			name:   "정상적인 객실 조회",
			roomID: "1",
			mockRoom: &models.Room{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					BaseTimeEntity: models.BaseTimeEntity{
						BaseEntity: models.BaseEntity{ID: 1},
					},
				},
				Number:      "101",
				RoomGroupID: 1,
				Note:        "스탠다드 룸",
				Status:      models.RoomStatusNormal,
				RoomGroup: &models.RoomGroup{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{
							BaseEntity: models.BaseEntity{ID: 1},
						},
					},
					Name: "스탠다드 그룹",
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "존재하지 않는 객실",
			roomID:         "999",
			mockRoom:       nil,
			mockError:      services.ErrRoomNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "잘못된 ID 형식",
			roomID:         "invalid",
			mockRoom:       nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockRoomService)
			mockUserService := new(MockUserService)
			mockHistoryService := new(MockHistoryService)
			handler := NewRoomHandler(mockService, mockUserService, mockHistoryService)

			// Set up mock expectations if valid ID
			if tt.roomID != "invalid" {
				roomID := uint(1)
				if tt.roomID == "999" {
					roomID = uint(999)
				}
				mockService.On("GetByIDWithGroup", mock.Anything, roomID).Return(tt.mockRoom, tt.mockError)
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/rooms/:id", handler.GetRoom)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/rooms/%s", tt.roomID), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Debug: 실제 응답 내용 확인
				t.Logf("Response body: %s", w.Body.String())

				// Parse response
				var actualResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
				assert.NoError(t, err)

				// 실제 응답 구조 확인 (value 필드 사용)
				if data, ok := actualResponse["value"].(map[string]interface{}); ok {
					if number, exists := data["number"]; exists {
						assert.Equal(t, tt.mockRoom.Number, number)
					} else {
						t.Errorf("number field not found in response value")
					}
				} else {
					t.Errorf("value field not found or not a map in response")
				}
			}

			// Verify mock expectations
			if tt.roomID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestRoomHandler_CreateRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockRoomService)
		expectedStatus int
	}{
		{
			name: "정상적인 객실 생성",
			requestBody: dto.CreateRoomRequest{
				Number:      "103",
				Note:        "스위트 룸",
				Status:      "NORMAL",
				RoomGroupID: 1,
			},
			setupMock: func(m *MockRoomService) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(room *models.Room) bool {
					return room.Number == "103" && room.Note == "스위트 룸"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "중복된 객실 번호",
			requestBody: dto.CreateRoomRequest{
				Number:      "101",
				Note:        "스탠다드 룸",
				Status:      "NORMAL",
				RoomGroupID: 1,
			},
			setupMock: func(m *MockRoomService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.Room")).
					Return(services.ErrRoomNumberExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "잘못된 요청 데이터",
			requestBody:    "invalid json",
			setupMock:      func(m *MockRoomService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "필수 필드 누락",
			requestBody: dto.CreateRoomRequest{
				Note: "스탠다드 룸",
				// Number 필드 누락
			},
			setupMock:      func(m *MockRoomService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockRoomService)
			mockUserService := new(MockUserService)
			mockHistoryService := new(MockHistoryService)
			handler := NewRoomHandler(mockService, mockUserService, mockHistoryService)

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
			router.POST("/api/v1/rooms", handler.CreateRoom)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/rooms", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				// Parse response
				var response dto.RoomResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestRoomHandler_UpdateRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		roomID         string
		requestBody    interface{}
		mockRoom       *models.Room
		mockError      error
		expectedStatus int
	}{
		{
			name:   "정상적인 객실 업데이트",
			roomID: "1",
			requestBody: dto.UpdateRoomRequest{
				Note:   stringPtr("업데이트된 스탠다드 룸"),
				Status: stringPtr("CONSTRUCTION"),
			},
			mockRoom: &models.Room{
				BaseMustAuditEntity: models.BaseMustAuditEntity{
					BaseTimeEntity: models.BaseTimeEntity{
						BaseEntity: models.BaseEntity{ID: 1},
					},
				},
				Number:      "101",
				RoomGroupID: 1,
				Note:        "업데이트된 스탠다드 룸",
				Status:      models.RoomStatusConstruction,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "존재하지 않는 객실",
			roomID: "999",
			requestBody: dto.UpdateRoomRequest{
				Note: stringPtr("업데이트된 룸"),
			},
			mockRoom:       nil,
			mockError:      services.ErrRoomNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "잘못된 ID 형식",
			roomID:         "invalid",
			requestBody:    dto.UpdateRoomRequest{},
			mockRoom:       nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockRoomService)
			mockUserService := new(MockUserService)
			mockHistoryService := new(MockHistoryService)
			handler := NewRoomHandler(mockService, mockUserService, mockHistoryService)

			// Set up mock expectations if valid ID
			if tt.roomID != "invalid" {
				roomID := uint(1)
				if tt.roomID == "999" {
					roomID = uint(999)
				}
				mockService.On("Update", mock.Anything, roomID, mock.AnythingOfType("map[string]interface {}")).
					Return(tt.mockRoom, tt.mockError)
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
			router.PATCH("/api/v1/rooms/:id", handler.UpdateRoom)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%s", tt.roomID), bytes.NewBuffer(body))
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
					if number, exists := data["number"]; exists {
						assert.Equal(t, tt.mockRoom.Number, number)
					}
				}
			}

			// Verify mock expectations
			if tt.roomID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestRoomHandler_DeleteRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		roomID         string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "정상적인 객실 삭제",
			roomID:         "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "존재하지 않는 객실",
			roomID:         "999",
			mockError:      services.ErrRoomNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "예약이 있는 객실",
			roomID:         "2",
			mockError:      services.ErrRoomHasReservations,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "잘못된 ID 형식",
			roomID:         "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockRoomService)
			mockUserService := new(MockUserService)
			mockHistoryService := new(MockHistoryService)
			handler := NewRoomHandler(mockService, mockUserService, mockHistoryService)

			// Set up mock expectations if valid ID
			if tt.roomID != "invalid" {
				roomID := uint(1)
				if tt.roomID == "999" {
					roomID = uint(999)
				} else if tt.roomID == "2" {
					roomID = uint(2)
				}
				mockService.On("Delete", mock.Anything, roomID).Return(tt.mockError)
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
			router.DELETE("/api/v1/rooms/:id", handler.DeleteRoom)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/rooms/%s", tt.roomID), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify mock expectations
			if tt.roomID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestRoomHandler_GetRoomHistories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		roomID         string
		queryParams    string
		setupMocks     func(*MockRoomService, *MockUserService, *MockHistoryService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "객실 히스토리를 성공적으로 조회할 수 있다",
			roomID:      "1",
			queryParams: "page=0&size=10",
			setupMocks: func(mockRoomService *MockRoomService, mockUserService *MockUserService, mockHistoryService *MockHistoryService) {
				room := &models.Room{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
					},
					Number: "101",
				}
				mockRoomService.On("GetByIDWithGroup", mock.Anything, uint(1)).Return(room, nil)

				revisions := []dto.RoomRevisionResponse{
					{
						Entity: dto.RoomResponse{
							ID:     1,
							Number: "101",
							Status: "NORMAL",
						},
						HistoryType:      "CREATED",
						HistoryCreatedAt: dto.CustomTime{Time: time.Now()},
						UpdatedFields:    []string{},
					},
				}
				mockHistoryService.On("GetRoomHistory", mock.Anything, uint(1), 0, 10).Return(revisions, int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"historyType":"CREATED"`,
		},
		{
			name:        "존재하지 않는 객실 ID로 조회하면 404를 반환한다",
			roomID:      "999",
			queryParams: "page=0&size=10",
			setupMocks: func(mockRoomService *MockRoomService, mockUserService *MockUserService, mockHistoryService *MockHistoryService) {
				mockRoomService.On("GetByIDWithGroup", mock.Anything, uint(999)).Return(nil, services.ErrRoomNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"message":"존재하지 않는 객실"`,
		},
		{
			name:        "잘못된 객실 ID 형식이면 400을 반환한다",
			roomID:      "invalid",
			queryParams: "page=0&size=10",
			setupMocks: func(mockRoomService *MockRoomService, mockUserService *MockUserService, mockHistoryService *MockHistoryService) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"잘못된 객실 ID"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoomService := new(MockRoomService)
			mockUserService := new(MockUserService)
			mockHistoryService := new(MockHistoryService)

			handler := NewRoomHandler(mockRoomService, mockUserService, mockHistoryService)

			tt.setupMocks(mockRoomService, mockUserService, mockHistoryService)

			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/rooms/:id/histories", handler.GetRoomHistories)

			url := fmt.Sprintf("/rooms/%s/histories", tt.roomID)
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockRoomService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
			mockHistoryService.AssertExpectations(t)
		})
	}
}
