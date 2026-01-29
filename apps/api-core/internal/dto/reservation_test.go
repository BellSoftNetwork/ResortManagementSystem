package dto_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
)

func TestReservationFilter(t *testing.T) {
	// Given: 테스트 환경 설정
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		queryString string
		expectError bool
		validate    func(*dto.ReservationFilter)
	}{
		{
			name:        "날짜 필터가 정상적으로 파싱된다",
			queryString: "stayStartAt=2025-05-26&stayEndAt=2025-07-06&status=NORMAL&type=STAY",
			expectError: false,
			validate: func(f *dto.ReservationFilter) {
				assert.NotNil(t, f.StayStartAt)
				assert.NotNil(t, f.StayEndAt)
				assert.Equal(t, "2025-05-26", f.StayStartAt.Format("2006-01-02"))
				assert.Equal(t, "2025-07-06", f.StayEndAt.Format("2006-01-02"))
				assert.Equal(t, "NORMAL", *f.Status)
				assert.Equal(t, "STAY", *f.Type)
			},
		},
		{
			name:        "size=200 같은 큰 값도 페이지네이션과 함께 사용 가능하다",
			queryString: "stayStartAt=2025-05-26&stayEndAt=2025-07-06&status=NORMAL&type=STAY&size=200",
			expectError: false,
			validate: func(f *dto.ReservationFilter) {
				// ReservationFilter 자체는 size를 포함하지 않음
				assert.NotNil(t, f.StayStartAt)
				assert.NotNil(t, f.StayEndAt)
			},
		},
		{
			name:        "월별 예약 조회가 가능하다",
			queryString: "stayStartAt=2025-06-01&stayEndAt=2025-06-30",
			expectError: false,
			validate: func(f *dto.ReservationFilter) {
				assert.NotNil(t, f.StayStartAt)
				assert.NotNil(t, f.StayEndAt)
				assert.Equal(t, time.Month(6), f.StayStartAt.Month())
				assert.Equal(t, time.Month(6), f.StayEndAt.Month())
			},
		},
		{
			name:        "잘못된 상태값은 검증 에러가 발생한다",
			queryString: "status=INVALID",
			expectError: true,
		},
		{
			name:        "잘못된 타입값은 검증 에러가 발생한다",
			queryString: "type=INVALID",
			expectError: true,
		},
		{
			name:        "검색어로 예약을 필터링할 수 있다",
			queryString: "search=홍길동",
			expectError: false,
			validate: func(f *dto.ReservationFilter) {
				assert.Equal(t, "홍길동", f.Search)
			},
		},
		{
			name:        "객실 ID로 필터링할 수 있다",
			queryString: "roomId=101",
			expectError: false,
			validate: func(f *dto.ReservationFilter) {
				assert.NotNil(t, f.RoomID)
				assert.Equal(t, uint(101), *f.RoomID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: HTTP 요청 생성
			req := httptest.NewRequest("GET", "/test?"+tt.queryString, nil)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = req

			// When: 필터 바인딩 실행
			var filter dto.ReservationFilter
			err := c.ShouldBindQuery(&filter)

			// Then: 결과 검증
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(&filter)
				}
			}
		})
	}
}

func TestReservationStatisticsQuery(t *testing.T) {
	// Given: 테스트 환경 설정
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		queryString  string
		expectError  bool
		errorMessage string
		validate     func(*dto.ReservationStatisticsQuery)
	}{
		{
			name:        "월별 통계 조회가 정상적으로 파싱된다",
			queryString: "startDate=2023-07-01&endDate=2025-06-30&periodType=MONTHLY",
			expectError: false,
			validate: func(q *dto.ReservationStatisticsQuery) {
				assert.Equal(t, "2023-07-01", q.StartDate.Format("2006-01-02"))
				assert.Equal(t, "2025-06-30", q.EndDate.Format("2006-01-02"))
				assert.Equal(t, "MONTHLY", q.PeriodType)
			},
		},
		{
			name:        "일별 통계 조회가 가능하다",
			queryString: "startDate=2025-06-01&endDate=2025-06-30&periodType=DAILY",
			expectError: false,
			validate: func(q *dto.ReservationStatisticsQuery) {
				assert.Equal(t, "DAILY", q.PeriodType)
			},
		},
		{
			name:        "연도별 통계 조회가 가능하다",
			queryString: "startDate=2020-01-01&endDate=2025-12-31&periodType=YEARLY",
			expectError: false,
			validate: func(q *dto.ReservationStatisticsQuery) {
				assert.Equal(t, "YEARLY", q.PeriodType)
			},
		},
		{
			name:         "시작일이 없으면 에러가 발생한다",
			queryString:  "endDate=2025-06-30&periodType=MONTHLY",
			expectError:  true,
			errorMessage: "required",
		},
		{
			name:         "종료일이 없으면 에러가 발생한다",
			queryString:  "startDate=2023-07-01&periodType=MONTHLY",
			expectError:  true,
			errorMessage: "required",
		},
		{
			name:         "잘못된 기간 타입은 에러가 발생한다",
			queryString:  "startDate=2023-07-01&endDate=2025-06-30&periodType=INVALID",
			expectError:  true,
			errorMessage: "oneof",
		},
		{
			name:        "기간 타입 없이도 조회가 가능하다 (기본값 사용)",
			queryString: "startDate=2023-07-01&endDate=2025-06-30",
			expectError: false,
			validate: func(q *dto.ReservationStatisticsQuery) {
				assert.Equal(t, "", q.PeriodType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: HTTP 요청 생성
			req := httptest.NewRequest("GET", "/test?"+tt.queryString, nil)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = req

			// When: 쿼리 바인딩 실행
			var query dto.ReservationStatisticsQuery
			err := c.ShouldBindQuery(&query)

			// Then: 결과 검증
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(&query)
				}
			}
		})
	}
}

func TestCreateReservationRequest(t *testing.T) {
	// Given: 정상적인 예약 생성 요청 데이터
	validRequest := dto.CreateReservationRequest{
		PaymentMethodID: 1,
		RoomIDs:         []uint{101, 102},
		Name:            "홍길동",
		Phone:           "010-1234-5678",
		PeopleCount:     4,
		StayStartAt:     dto.JSONTime{Time: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)},
		StayEndAt:       dto.JSONTime{Time: time.Date(2025, 7, 3, 0, 0, 0, 0, time.UTC)},
		Price:           200000,
		Deposit:         100000,
		PaymentAmount:   100000,
		Note:            "조용한 방 요청",
		Type:            "STAY",
	}

	tests := []struct {
		name        string
		modify      func(*dto.CreateReservationRequest)
		expectError bool
		errorField  string
	}{
		{
			name:        "정상적인 예약 생성 요청은 검증을 통과한다",
			modify:      func(r *dto.CreateReservationRequest) {},
			expectError: false,
		},
		{
			name: "이름이 비어있으면 에러가 발생한다",
			modify: func(r *dto.CreateReservationRequest) {
				r.Name = ""
			},
			expectError: true,
			errorField:  "Name",
		},
		{
			name: "이름이 30자를 초과하면 에러가 발생한다",
			modify: func(r *dto.CreateReservationRequest) {
				r.Name = "아주긴이름아주긴이름아주긴이름아주긴이름아주긴이름아주긴이름아" // 31자
			},
			expectError: true,
			errorField:  "Name",
		},
		{
			name: "전화번호가 비어있어도 정상 처리된다 (optional)",
			modify: func(r *dto.CreateReservationRequest) {
				r.Phone = ""
			},
			expectError: false,
		},
		{
			name: "월세 타입도 정상적으로 처리된다",
			modify: func(r *dto.CreateReservationRequest) {
				r.Type = "MONTHLY_RENT"
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: 요청 데이터 준비 및 JSON 변환
			req := validRequest
			tt.modify(&req)

			// JSON으로 변환
			jsonData, err := json.Marshal(req)
			assert.NoError(t, err)

			// When: gin 컨텍스트로 바인딩 검증
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/test", bytes.NewReader(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			var boundReq dto.CreateReservationRequest
			err = c.ShouldBindJSON(&boundReq)

			// Then: 결과 검증
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorField != "" && err != nil {
					assert.Contains(t, err.Error(), tt.errorField)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateReservationRequest_JSONStringTypeRejection(t *testing.T) {
	// Given: 테스트 환경 설정
	gin.SetMode(gin.TestMode)

	t.Run("price 필드가 문자열로 전달되면 바인딩 에러가 발생한다", func(t *testing.T) {
		// Given: price가 문자열인 JSON
		jsonBody := `{"name":"테스트","phone":"010-1234-5678","price":"100000","stayStartAt":"2025-07-01","stayEndAt":"2025-07-03"}`

		// When
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = req

		var boundReq dto.CreateReservationRequest
		err := c.ShouldBindJSON(&boundReq)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot unmarshal string")
	})
}
