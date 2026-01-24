package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.bellsoft.net/rms/api-core/internal/handlers"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

func TestAuthHandler_BruteForceProtection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("30분 이내 5번 초과 실패 시 계정 잠김", func(t *testing.T) {
		// Given: 인증 서비스 모킹
		mockAuthService := new(MockAuthService)
		authHandler := handlers.NewAuthHandler(mockAuthService)

		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/v1/auth/login", authHandler.Login)

		// 처음 5번은 인증 실패로 설정
		mockAuthService.On("Login", mock.Anything, "testuser", "wrongpassword", "127.0.0.1", mock.Anything).
			Return(nil, services.ErrInvalidCredentials).Times(5)

		// 6번째부터는 Too Many Attempts 에러로 설정
		mockAuthService.On("Login", mock.Anything, "testuser", mock.Anything, "127.0.0.1", mock.Anything).
			Return(nil, services.ErrTooManyAttempts).Times(1)

		loginData := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		var statusCodes []int

		// When: 6번 연속 로그인 시도
		for i := 0; i < 6; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Forwarded-For", "127.0.0.1")
			router.ServeHTTP(w, req)

			statusCodes = append(statusCodes, w.Code)
		}

		// Then: 첫 5번은 401, 마지막은 429
		for i := 0; i < 5; i++ {
			assert.Equal(t, http.StatusUnauthorized, statusCodes[i], "로그인 시도 %d는 401이어야 함", i+1)
		}
		assert.Equal(t, http.StatusTooManyRequests, statusCodes[5], "6번째 로그인 시도는 429여야 함")
	})

	t.Run("계정 잠김 후 정상 비밀번호로도 로그인 불가", func(t *testing.T) {
		// Given: 계정이 이미 잠긴 상태
		mockAuthService := new(MockAuthService)
		authHandler := handlers.NewAuthHandler(mockAuthService)

		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/v1/auth/login", authHandler.Login)

		// 계정 잠김 상태로 설정
		mockAuthService.On("Login", mock.Anything, "testuser", "correctpassword", "127.0.0.1", mock.Anything).
			Return(nil, services.ErrTooManyAttempts)

		// When: 정상 비밀번호로 로그인 시도
		loginData := map[string]string{
			"username": "testuser",
			"password": "correctpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "127.0.0.1")
		router.ServeHTTP(w, req)

		// Then: 429 Too Many Requests 응답
		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "로그인 시도 횟수 초과. 잠시 후 다시 시도하세요")
	})

	t.Run("서로 다른 IP에서의 로그인 시도는 독립적으로 처리", func(t *testing.T) {
		// Given: 두 개의 다른 IP 주소
		mockAuthService := new(MockAuthService)
		authHandler := handlers.NewAuthHandler(mockAuthService)

		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/v1/auth/login", authHandler.Login)

		// IP1에서는 계정 잠김
		mockAuthService.On("Login", mock.Anything, "testuser", "wrongpassword", "192.168.1.100", mock.Anything).
			Return(nil, services.ErrTooManyAttempts)

		// IP2에서는 정상 로그인 가능
		expectedUser := &models.User{
			UserID: "testuser",
			Name:   "Test User",
			Role:   models.UserRoleNormal,
		}
		expectedUser.ID = 1

		mockAuthService.On("Login", mock.Anything, "testuser", "correctpassword", "192.168.1.200", mock.Anything).
			Return(&services.LoginResponse{
				User:                 expectedUser,
				AccessToken:          "valid.access.token",
				RefreshToken:         "valid.refresh.token",
				AccessTokenExpiresIn: time.Now().Add(15 * time.Minute).Unix(),
			}, nil)

		// When: IP1에서 로그인 시도 (계정 잠김)
		loginData1 := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		jsonData1, _ := json.Marshal(loginData1)

		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("X-Forwarded-For", "192.168.1.100")
		router.ServeHTTP(w1, req1)

		// And: IP2에서 로그인 시도 (정상 로그인)
		loginData2 := map[string]string{
			"username": "testuser",
			"password": "correctpassword",
		}
		jsonData2, _ := json.Marshal(loginData2)

		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData2))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("X-Forwarded-For", "192.168.1.200")
		router.ServeHTTP(w2, req2)

		// Then: IP1은 잠김, IP2는 성공
		assert.Equal(t, http.StatusTooManyRequests, w1.Code)
		assert.Equal(t, http.StatusOK, w2.Code)
	})

	t.Run("서로 다른 사용자의 로그인 시도는 독립적으로 처리", func(t *testing.T) {
		// Given: 두 명의 다른 사용자
		mockAuthService := new(MockAuthService)
		authHandler := handlers.NewAuthHandler(mockAuthService)

		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/v1/auth/login", authHandler.Login)

		// user1은 계정 잠김
		mockAuthService.On("Login", mock.Anything, "user1", "wrongpassword", "127.0.0.1", mock.Anything).
			Return(nil, services.ErrTooManyAttempts)

		// user2는 정상 로그인
		expectedUser2 := &models.User{
			UserID: "user2",
			Name:   "User Two",
			Role:   models.UserRoleNormal,
		}
		expectedUser2.ID = 2

		mockAuthService.On("Login", mock.Anything, "user2", "correctpassword", "127.0.0.1", mock.Anything).
			Return(&services.LoginResponse{
				User:                 expectedUser2,
				AccessToken:          "user2.access.token",
				RefreshToken:         "user2.refresh.token",
				AccessTokenExpiresIn: time.Now().Add(15 * time.Minute).Unix(),
			}, nil)

		// When: user1 로그인 시도 (잠김)
		loginData1 := map[string]string{
			"username": "user1",
			"password": "wrongpassword",
		}
		jsonData1, _ := json.Marshal(loginData1)

		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("X-Forwarded-For", "127.0.0.1")
		router.ServeHTTP(w1, req1)

		// And: user2 로그인 시도 (성공)
		loginData2 := map[string]string{
			"username": "user2",
			"password": "correctpassword",
		}
		jsonData2, _ := json.Marshal(loginData2)

		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData2))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("X-Forwarded-For", "127.0.0.1")
		router.ServeHTTP(w2, req2)

		// Then: user1은 잠김, user2는 성공
		assert.Equal(t, http.StatusTooManyRequests, w1.Code)
		assert.Equal(t, http.StatusOK, w2.Code)
	})

	t.Run("로그인 시도 제한 응답 메시지 확인", func(t *testing.T) {
		// Given: 계정 잠김 상태
		mockAuthService := new(MockAuthService)
		authHandler := handlers.NewAuthHandler(mockAuthService)

		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/v1/auth/login", authHandler.Login)

		mockAuthService.On("Login", mock.Anything, "testuser", "password", "127.0.0.1", mock.Anything).
			Return(nil, services.ErrTooManyAttempts)

		// When: 로그인 시도
		loginData := map[string]string{
			"username": "testuser",
			"password": "password",
		}
		jsonData, _ := json.Marshal(loginData)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "127.0.0.1")
		router.ServeHTTP(w, req)

		// Then: 429 응답과 적절한 메시지
		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// 응답 메시지가 보안상 중요한 정보를 노출하지 않는지 확인
		message, ok := response["message"].(string)
		assert.True(t, ok)
		assert.NotEmpty(t, message)

		// 내부 구현 세부사항이 노출되지 않았는지 확인
		assert.NotContains(t, message, "database")
		assert.NotContains(t, message, "redis")
		assert.NotContains(t, message, "internal")
	})
}
