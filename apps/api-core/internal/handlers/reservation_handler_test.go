package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
)

// MockReservationService는 ReservationService의 모킹 구현
type MockReservationService struct {
	mock.Mock
}

func (m *MockReservationService) GetByID(ctx context.Context, id uint) (*models.Reservation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reservation), args.Error(1)
}

func (m *MockReservationService) GetByIDWithDetails(ctx context.Context, id uint) (*models.Reservation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reservation), args.Error(1)
}

func (m *MockReservationService) GetAll(ctx context.Context, filter dto.ReservationRepositoryFilter, page, size int, sort string) ([]models.Reservation, int64, error) {
	args := m.Called(ctx, filter, page, size, sort)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Reservation), args.Get(1).(int64), args.Error(2)
}

func (m *MockReservationService) GetStatistics(ctx context.Context, startDate, endDate time.Time, periodType string) ([]repositories.ReservationStatistics, error) {
	args := m.Called(ctx, startDate, endDate, periodType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repositories.ReservationStatistics), args.Error(1)
}

func (m *MockReservationService) Create(ctx context.Context, reservation *models.Reservation, roomIDs []uint) error {
	args := m.Called(ctx, reservation, roomIDs)
	return args.Error(0)
}

func (m *MockReservationService) Update(ctx context.Context, id uint, updates map[string]interface{}, roomIDs []uint, hasRoomsUpdate bool) (*models.Reservation, error) {
	args := m.Called(ctx, id, updates, roomIDs, hasRoomsUpdate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reservation), args.Error(1)
}

func (m *MockReservationService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReservationService) GetAvailableRooms(ctx context.Context, startDate, endDate time.Time, excludeReservationID *uint) ([]models.Room, error) {
	args := m.Called(ctx, startDate, endDate, excludeReservationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Room), args.Error(1)
}

func (m *MockReservationService) GetLastReservationForRoom(ctx context.Context, roomID uint) (*models.Reservation, error) {
	args := m.Called(ctx, roomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reservation), args.Error(1)
}

// MockUserService는 UserService의 모킹 구현
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetByUserID(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.User, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) GetAll(ctx context.Context, page, size int) ([]models.User, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) UpdatePassword(ctx context.Context, id uint, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func (m *MockUserService) IsUpdatableAccount(ctx context.Context, requestUser *models.User, targetUserID uint) (bool, error) {
	args := m.Called(ctx, requestUser, targetUserID)
	return args.Bool(0), args.Error(1)
}

func TestReservationHandler_ListReservations_Sorting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		queryParams      string
		expectedSort     string
		mockReservations []models.Reservation
		mockTotal        int64
		expectError      bool
	}{
		{
			name:         "가격 내림차순 정렬",
			queryParams:  "page=0&size=15&sort=price,desc",
			expectedSort: "price,desc",
			mockReservations: []models.Reservation{
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}}}, Price: 100000},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}}}, Price: 80000},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 3}}}, Price: 60000},
			},
			mockTotal:   3,
			expectError: false,
		},
		{
			name:         "가격 오름차순 정렬",
			queryParams:  "page=0&size=15&sort=price,asc",
			expectedSort: "price,asc",
			mockReservations: []models.Reservation{
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 3}}}, Price: 60000},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}}}, Price: 80000},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}}}, Price: 100000},
			},
			mockTotal:   3,
			expectError: false,
		},
		{
			name:         "체크아웃 날짜 오름차순 정렬",
			queryParams:  "page=0&size=15&sort=stayEndAt,asc",
			expectedSort: "stayEndAt,asc",
			mockReservations: []models.Reservation{
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}}}, StayEndAt: time.Now().AddDate(0, 0, 1)},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}}}, StayEndAt: time.Now().AddDate(0, 0, 2)},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 3}}}, StayEndAt: time.Now().AddDate(0, 0, 3)},
			},
			mockTotal:   3,
			expectError: false,
		},
		{
			name:         "정렬 파라미터 없을 때 기본값",
			queryParams:  "page=0&size=15",
			expectedSort: "",
			mockReservations: []models.Reservation{
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 3}}}},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}}}},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}}}},
			},
			mockTotal:   3,
			expectError: false,
		},
		{
			name:         "여러 필드 정렬",
			queryParams:  "page=0&size=15&sort=price,desc,stayEndAt,asc",
			expectedSort: "price,desc,stayEndAt,asc",
			mockReservations: []models.Reservation{
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}}}, Price: 100000},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}}}, Price: 100000},
				{BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 3}}}, Price: 80000},
			},
			mockTotal:   3,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockReservationService := new(MockReservationService)
			mockUserService := new(MockUserService)
			mockAuditService := new(MockHistoryService)
			handler := NewReservationHandler(mockReservationService, mockUserService, mockAuditService)

			// Set up mock expectations
			mockReservationService.On("GetAll",
				mock.Anything,
				mock.AnythingOfType("dto.ReservationRepositoryFilter"),
				0,
				15,
				tt.expectedSort,
			).Return(tt.mockReservations, tt.mockTotal, nil)

			// Mock user service for getUserSummary calls
			mockUserService.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(&models.User{
				BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
				UserID:         "testuser",
				Email:          stringPtr("test@example.com"),
				Name:           "Test User",
			}, nil).Maybe()

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/reservations", handler.ListReservations)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/reservations?%s", tt.queryParams), nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			if tt.expectError {
				assert.NotEqual(t, http.StatusOK, w.Code)
			} else {
				if w.Code != http.StatusOK {
					t.Logf("Response body: %s", w.Body.String())
				}
				assert.Equal(t, http.StatusOK, w.Code)

				// Parse response
				var response struct {
					Values []dto.ReservationResponse `json:"values"`
					Page   map[string]interface{}    `json:"page"`
					Filter interface{}               `json:"filter"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, int(tt.mockTotal), len(response.Values))
			}

			// Verify mock expectations
			mockReservationService.AssertExpectations(t)
		})
	}
}

func TestReservationHandler_ListReservations_WithFiltersAndSort(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockReservationService := new(MockReservationService)
	mockUserService := new(MockUserService)
	mockAuditService := new(MockHistoryService)
	handler := NewReservationHandler(mockReservationService, mockUserService, mockAuditService)

	// Test data
	stayStartAt := time.Now()
	stayEndAt := time.Now().AddDate(0, 0, 3)
	status := models.ReservationStatusNormal
	reservationType := models.ReservationTypeStay

	// expectedFilter는 mock 매칭에 사용됨
	_ = dto.ReservationRepositoryFilter{
		Status:    &status,
		Type:      &reservationType,
		StartDate: &stayStartAt,
		EndDate:   &stayEndAt,
	}

	mockReservations := []models.Reservation{
		{
			BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}}},
			Price:               100000,
			Status:              status,
			Type:                reservationType,
			StayStartAt:         stayStartAt,
			StayEndAt:           stayEndAt,
		},
		{
			BaseMustAuditEntity: models.BaseMustAuditEntity{BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 2}}},
			Price:               80000,
			Status:              status,
			Type:                reservationType,
			StayStartAt:         stayStartAt,
			StayEndAt:           stayEndAt,
		},
	}

	// Set up mock expectations
	mockReservationService.On("GetAll",
		mock.Anything,
		mock.MatchedBy(func(filter dto.ReservationRepositoryFilter) bool {
			return filter.Status != nil && *filter.Status == status &&
				filter.Type != nil && *filter.Type == reservationType
		}),
		0,
		15,
		"price,desc",
	).Return(mockReservations, int64(2), nil)

	// Mock user service
	mockUserService.On("GetByID", mock.Anything, mock.AnythingOfType("uint")).Return(&models.User{
		BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
		UserID:         "testuser",
		Email:          stringPtr("test@example.com"),
		Name:           "Test User",
	}, nil).Maybe()

	// Create test request
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.GET("/api/v1/reservations", handler.ListReservations)

	queryParams := fmt.Sprintf("page=0&size=15&sort=price,desc&stayStartAt=%s&stayEndAt=%s&status=NORMAL&type=STAY",
		stayStartAt.Format("2006-01-02"),
		stayEndAt.Format("2006-01-02"))

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/reservations?%s", queryParams), nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response struct {
		Values []dto.ReservationResponse     `json:"values"`
		Page   map[string]interface{}        `json:"page"`
		Filter dto.ReservationFilterResponse `json:"filter"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.Values))

	// Verify filter response
	assert.NotNil(t, response.Filter.Status)
	assert.Equal(t, "NORMAL", *response.Filter.Status)
	assert.NotNil(t, response.Filter.Type)
	assert.Equal(t, "STAY", *response.Filter.Type)

	// Verify mock expectations
	mockReservationService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestReservationHandler_GetStatistics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		mockStats      []repositories.ReservationStatistics
		mockError      error
		expectedStatus int
	}{
		{
			name:        "일별 통계 조회",
			queryParams: "startDate=2024-01-01&endDate=2024-01-31&periodType=DAILY",
			mockStats: []repositories.ReservationStatistics{
				{
					Period:           "2024-01-01",
					ReservationCount: 5,
					TotalRevenue:     500000,
					TotalGuests:      15,
					AverageStayDays:  2.5,
				},
				{
					Period:           "2024-01-02",
					ReservationCount: 8,
					TotalRevenue:     960000,
					TotalGuests:      24,
					AverageStayDays:  3.0,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:        "월별 통계 조회 - 다른 월",
			queryParams: "startDate=2024-02-01&endDate=2024-02-29&periodType=MONTHLY",
			mockStats: []repositories.ReservationStatistics{
				{
					Period:           "2024-02",
					ReservationCount: 150,
					TotalRevenue:     18000000,
					TotalGuests:      450,
					AverageStayDays:  3.0,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:        "월별 통계 조회",
			queryParams: "startDate=2024-01-01&endDate=2024-12-31&periodType=MONTHLY",
			mockStats: []repositories.ReservationStatistics{
				{
					Period:           "2024-01",
					ReservationCount: 150,
					TotalRevenue:     18000000,
					TotalGuests:      450,
					AverageStayDays:  3.2,
				},
				{
					Period:           "2024-02",
					ReservationCount: 140,
					TotalRevenue:     16800000,
					TotalGuests:      420,
					AverageStayDays:  3.0,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:        "연도별 통계 조회",
			queryParams: "startDate=2023-01-01&endDate=2024-12-31&periodType=YEARLY",
			mockStats: []repositories.ReservationStatistics{
				{
					Period:           "2023",
					ReservationCount: 1800,
					TotalRevenue:     216000000,
					TotalGuests:      5400,
					AverageStayDays:  3.0,
				},
				{
					Period:           "2024",
					ReservationCount: 2000,
					TotalRevenue:     240000000,
					TotalGuests:      6000,
					AverageStayDays:  3.0,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "필수 파라미터 누락 - startDate",
			queryParams:    "endDate=2024-01-31&periodType=DAILY",
			mockStats:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "필수 파라미터 누락 - endDate",
			queryParams:    "startDate=2024-01-01&periodType=DAILY",
			mockStats:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "periodType 없이 기본값 사용",
			queryParams: "startDate=2024-01-01&endDate=2024-01-31",
			mockStats: []repositories.ReservationStatistics{
				{
					Period:           "2024-01",
					ReservationCount: 150,
					TotalRevenue:     18000000,
					TotalGuests:      450,
					AverageStayDays:  3.0,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "잘못된 날짜 형식",
			queryParams:    "startDate=invalid&endDate=2024-01-31&periodType=DAILY",
			mockStats:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "잘못된 periodType",
			queryParams:    "startDate=2024-01-01&endDate=2024-01-31&periodType=INVALID",
			mockStats:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "시작일이 종료일보다 늦은 경우",
			queryParams:    "startDate=2024-01-31&endDate=2024-01-01&periodType=DAILY",
			mockStats:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "서비스 에러 발생",
			queryParams:    "startDate=2024-01-01&endDate=2024-01-31&periodType=DAILY",
			mockStats:      nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "빈 결과",
			queryParams:    "startDate=2024-01-01&endDate=2024-01-31&periodType=DAILY",
			mockStats:      []repositories.ReservationStatistics{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockReservationService := new(MockReservationService)
			mockUserService := new(MockUserService)
			mockAuditService := new(MockHistoryService)
			handler := NewReservationHandler(mockReservationService, mockUserService, mockAuditService)

			// Parse query params to determine if we should set up mock
			shouldMock := true
			if tt.queryParams != "" {
				params := make(map[string]string)
				for _, param := range strings.Split(tt.queryParams, "&") {
					parts := strings.Split(param, "=")
					if len(parts) == 2 {
						params[parts[0]] = parts[1]
					}
				}

				// Check if all required params are present and valid
				startDate, hasStart := params["startDate"]
				endDate, hasEnd := params["endDate"]
				periodType, hasPeriod := params["periodType"]

				if !hasStart || !hasEnd {
					shouldMock = false
				} else {
					// Check date format
					_, err1 := time.Parse("2006-01-02", startDate)
					_, err2 := time.Parse("2006-01-02", endDate)

					if err1 != nil || err2 != nil {
						shouldMock = false
					} else {
						// If periodType is provided, validate it
						if hasPeriod {
							validPeriods := map[string]bool{"DAILY": true, "MONTHLY": true, "YEARLY": true}
							if !validPeriods[periodType] {
								shouldMock = false
							}
						}

						// Check date order only if shouldMock is still true
						if shouldMock {
							start, _ := time.Parse("2006-01-02", startDate)
							end, _ := time.Parse("2006-01-02", endDate)
							if start.After(end) {
								shouldMock = false
							}
						}
					}
				}
			}

			// Set up mock expectations if valid request
			if shouldMock {
				mockReservationService.On("GetStatistics",
					mock.Anything,
					mock.AnythingOfType("time.Time"),
					mock.AnythingOfType("time.Time"),
					mock.AnythingOfType("string"),
				).Return(tt.mockStats, tt.mockError)
			}

			// Create test request
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/api/v1/reservation-statistics", handler.GetReservationStatistics)

			url := "/api/v1/reservation-statistics"
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				// Parse response with value wrapper
				var responseWrapper struct {
					Value dto.ReservationStatisticsResponse `json:"value"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &responseWrapper)
				assert.NoError(t, err)

				response := responseWrapper.Value
				assert.Equal(t, len(tt.mockStats), len(response.Stats))

				// Verify response data - only if stats exist
				if len(response.Stats) > 0 && len(tt.mockStats) > 0 {
					for i, stat := range tt.mockStats {
						if i < len(response.Stats) {
							assert.Equal(t, stat.Period, response.Stats[i].Period)
							assert.Equal(t, int(stat.ReservationCount), response.Stats[i].TotalReservations)
							assert.Equal(t, int(stat.TotalRevenue), response.Stats[i].TotalSales)
						}
					}
				}
			}

			// Verify mock expectations only for successful cases
			if shouldMock && tt.expectedStatus == http.StatusOK {
				mockReservationService.AssertExpectations(t)
			}
		})
	}
}

func TestReservationHandler_GetReservationHistories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		reservationID  string
		queryParams    string
		setupMocks     func(*MockReservationService, *MockUserService, *MockHistoryService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "예약 히스토리를 성공적으로 조회할 수 있다",
			reservationID: "1",
			queryParams:   "page=0&size=10",
			setupMocks: func(mockReservationService *MockReservationService, mockUserService *MockUserService, mockHistoryService *MockHistoryService) {
				reservation := &models.Reservation{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
					},
					Name: "테스트 예약",
				}
				mockReservationService.On("GetByIDWithDetails", mock.Anything, uint(1)).Return(reservation, nil)

				revisions := []dto.ReservationRevisionResponse{
					{
						Entity: dto.ReservationResponse{
							ID:     1,
							Name:   "테스트 예약",
							Status: "PENDING",
						},
						HistoryType:      "CREATED",
						HistoryCreatedAt: dto.CustomTime{Time: time.Now()},
						UpdatedFields:    []string{},
					},
					{
						Entity: dto.ReservationResponse{
							ID:     1,
							Name:   "테스트 예약",
							Status: "NORMAL",
						},
						HistoryType:      "UPDATED",
						HistoryCreatedAt: dto.CustomTime{Time: time.Now()},
						UpdatedFields:    []string{"status"},
					},
				}
				mockHistoryService.On("GetReservationHistory", mock.Anything, uint(1), 0, 10).Return(revisions, int64(2), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"historyType":"CREATED"`,
		},
		{
			name:          "존재하지 않는 예약 ID로 조회하면 404를 반환한다",
			reservationID: "999",
			queryParams:   "page=0&size=10",
			setupMocks: func(mockReservationService *MockReservationService, mockUserService *MockUserService, mockHistoryService *MockHistoryService) {
				mockReservationService.On("GetByIDWithDetails", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"message":"존재하지 않는 예약"`,
		},
		{
			name:          "잘못된 예약 ID 형식이면 400을 반환한다",
			reservationID: "invalid",
			queryParams:   "page=0&size=10",
			setupMocks: func(mockReservationService *MockReservationService, mockUserService *MockUserService, mockHistoryService *MockHistoryService) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"잘못된 예약 ID"`,
		},
		{
			name:          "히스토리 서비스 에러 시 500을 반환한다",
			reservationID: "1",
			queryParams:   "page=0&size=10",
			setupMocks: func(mockReservationService *MockReservationService, mockUserService *MockUserService, mockHistoryService *MockHistoryService) {
				reservation := &models.Reservation{
					BaseMustAuditEntity: models.BaseMustAuditEntity{
						BaseTimeEntity: models.BaseTimeEntity{BaseEntity: models.BaseEntity{ID: 1}},
					},
				}
				mockReservationService.On("GetByIDWithDetails", mock.Anything, uint(1)).Return(reservation, nil)
				mockHistoryService.On("GetReservationHistory", mock.Anything, uint(1), 0, 10).Return(nil, int64(0), errors.New("history service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"message":"예약 이력 조회 실패"`,
		},
		{
			name:          "잘못된 쿼리 파라미터로 400을 반환한다",
			reservationID: "1",
			queryParams:   "page=invalid&size=10",
			setupMocks: func(mockReservationService *MockReservationService, mockUserService *MockUserService, mockHistoryService *MockHistoryService) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"잘못된 쿼리 파라미터"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReservationService := new(MockReservationService)
			mockUserService := new(MockUserService)
			mockHistoryService := new(MockHistoryService)

			handler := NewReservationHandler(mockReservationService, mockUserService, mockHistoryService)

			tt.setupMocks(mockReservationService, mockUserService, mockHistoryService)

			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/reservations/:id/histories", handler.GetReservationHistories)

			url := fmt.Sprintf("/reservations/%s/histories", tt.reservationID)
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockReservationService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
			mockHistoryService.AssertExpectations(t)
		})
	}
}
