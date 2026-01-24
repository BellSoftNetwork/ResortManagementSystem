package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gorm.io/gorm"
)

func TestErrorHandler(t *testing.T) {
	// Given: 테스트 환경 설정
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupError     func(*gin.Context)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "레코드를 찾을 수 없을 때 404 응답을 반환한다",
			setupError: func(c *gin.Context) {
				c.Error(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"message": "존재하지 않는 데이터",
			},
		},
		{
			name: "중복 키 에러 발생 시 409 응답을 반환한다",
			setupError: func(c *gin.Context) {
				c.Error(gorm.ErrDuplicatedKey)
			},
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"message": "이미 존재하는 데이터",
			},
		},
		{
			name: "검증 에러 발생 시 400 응답과 필드 에러를 반환한다",
			setupError: func(c *gin.Context) {
				type TestStruct struct {
					Size int `validate:"min=1,max=100"`
				}
				validate := validator.New()
				err := validate.Struct(&TestStruct{Size: 200})
				c.Error(err)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "잘못된 요청",
				"fieldErrors": []string{
					"'Size'은(는) 최대 100 이하여야 합니다 (요청 값: 200)",
				},
			},
		},
		{
			name: "일반 에러 발생 시 500 응답을 반환한다",
			setupError: func(c *gin.Context) {
				c.Error(errors.New("unexpected error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"message": "서버 오류",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: 라우터 설정
			router := gin.New()
			router.Use(middleware.ErrorHandler())
			router.GET("/test", func(c *gin.Context) {
				tt.setupError(c)
			})

			// When: 요청 실행
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			// Then: 응답 검증
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 간단한 검증을 위해 문자열로 확인
			body := w.Body.String()
			for key, value := range tt.expectedBody {
				if key == "message" {
					assert.Contains(t, body, value.(string))
				}
			}
		})
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	// Given: 테스트 환경 설정
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		panicValue     interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "문자열 패닉 발생 시 500 응답을 반환한다",
			panicValue:     "panic occurred",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"panic occurred","errors":null,"fieldErrors":null}`,
		},
		{
			name:           "기타 패닉 발생 시 일반 에러 메시지를 반환한다",
			panicValue:     123,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"서버 오류","errors":null,"fieldErrors":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: 라우터 설정
			router := gin.New()
			router.Use(middleware.RecoveryMiddleware())
			router.GET("/test", func(c *gin.Context) {
				panic(tt.panicValue)
			})

			// When: 요청 실행
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			// Then: 응답 검증
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
