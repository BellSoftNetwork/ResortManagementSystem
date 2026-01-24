package dto_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
)

func TestPaginationQuery(t *testing.T) {
	// Given: 테스트 환경 설정
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryString    string
		expectedPage   int
		expectedSize   int
		expectedSort   string
		expectedOffset int
		shouldError    bool
		errorMessage   string
	}{
		{
			name:           "기본값으로 page=0, size=20이 설정된다",
			queryString:    "",
			expectedPage:   0,
			expectedSize:   20,
			expectedSort:   "",
			expectedOffset: 0,
			shouldError:    false,
		},
		{
			name:           "Spring Boot 스타일의 0-based 페이징이 작동한다",
			queryString:    "page=0&size=15",
			expectedPage:   0,
			expectedSize:   15,
			expectedSort:   "",
			expectedOffset: 0,
			shouldError:    false,
		},
		{
			name:           "페이지 번호가 정상적으로 offset으로 변환된다",
			queryString:    "page=2&size=10",
			expectedPage:   2,
			expectedSize:   10,
			expectedSort:   "",
			expectedOffset: 20, // 2 * 10
			shouldError:    false,
		},
		{
			name:           "큰 size 값(1000)도 허용된다",
			queryString:    "page=0&size=1000",
			expectedPage:   0,
			expectedSize:   1000,
			expectedSort:   "",
			expectedOffset: 0,
			shouldError:    false,
		},
		{
			name:           "최대 size(2000)까지 허용된다",
			queryString:    "page=0&size=2000",
			expectedPage:   0,
			expectedSize:   2000,
			expectedSort:   "",
			expectedOffset: 0,
			shouldError:    false,
		},
		{
			name:         "size가 최대값(2000)을 초과하면 에러가 발생한다",
			queryString:  "page=0&size=2001",
			shouldError:  true,
			errorMessage: "Size",
		},
		{
			name:         "size가 0이면 에러가 발생한다",
			queryString:  "page=0&size=0",
			shouldError:  true,
			errorMessage: "Size",
		},
		{
			name:         "page가 음수면 에러가 발생한다",
			queryString:  "page=-1&size=20",
			shouldError:  true,
			errorMessage: "Page",
		},
		{
			name:           "정렬 파라미터가 정상적으로 파싱된다",
			queryString:    "page=0&size=20&sort=name,asc",
			expectedPage:   0,
			expectedSize:   20,
			expectedSort:   "name,asc",
			expectedOffset: 0,
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: HTTP 요청 생성
			req := httptest.NewRequest("GET", "/test?"+tt.queryString, nil)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = req

			// When: 쿼리 바인딩 실행
			var query dto.PaginationQuery
			err := c.ShouldBindQuery(&query)

			// Then: 결과 검증
			if tt.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPage, query.Page)
				assert.Equal(t, tt.expectedSize, query.Size)
				assert.Equal(t, tt.expectedSort, query.Sort)
				assert.Equal(t, tt.expectedOffset, query.GetOffset())
				assert.Equal(t, tt.expectedSize, query.GetLimit())
			}
		})
	}
}

func TestPaginationQuery_CompatibilityWithSpringBoot(t *testing.T) {
	// Given: Spring Boot에서 전송하는 실제 쿼리 파라미터들
	springBootQueries := []struct {
		name           string
		queryString    string
		expectedOffset int
		description    string
	}{
		{
			name:           "첫 페이지 요청",
			queryString:    "page=0&size=15&sort=number,asc",
			expectedOffset: 0,
			description:    "Spring Boot의 첫 페이지는 0부터 시작한다",
		},
		{
			name:           "두 번째 페이지 요청",
			queryString:    "page=1&size=15&sort=name,desc",
			expectedOffset: 15,
			description:    "두 번째 페이지의 offset은 15여야 한다",
		},
		{
			name:           "대용량 데이터 조회",
			queryString:    "page=0&size=1000",
			expectedOffset: 0,
			description:    "통계나 대시보드용 대용량 조회가 가능해야 한다",
		},
	}

	for _, tt := range springBootQueries {
		t.Run(tt.name, func(t *testing.T) {
			// Given: Spring Boot 스타일 요청
			req := httptest.NewRequest("GET", "/test?"+tt.queryString, nil)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = req

			// When: 파싱 실행
			var query dto.PaginationQuery
			err := c.ShouldBindQuery(&query)

			// Then: Spring Boot와 호환되는지 검증
			assert.NoError(t, err, tt.description)
			assert.Equal(t, tt.expectedOffset, query.GetOffset(), tt.description)
		})
	}
}
