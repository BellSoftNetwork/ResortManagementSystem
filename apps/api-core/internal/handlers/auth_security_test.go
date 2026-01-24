package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.bellsoft.net/rms/api-core/internal/handlers"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/auth"
)

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

func TestAuthHandler_SecurityTests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CSRF 공격 방어 - 쿠키에 토큰을 저장해도 인증되지 않는다", func(t *testing.T) {
		// Given: JWT 서비스 설정
		jwtService := auth.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)

		// 유효한 토큰 생성
		validToken, err := jwtService.GenerateAccessToken(1, "testuser", "USER")
		assert.NoError(t, err)

		// 라우터 설정 - 인증이 필요한 엔드포인트
		router := gin.New()
		router.Use(middleware.AuthMiddleware(jwtService))
		router.GET("/api/v1/my", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// When: 쿠키에 토큰을 저장하고 악의적인 사이트에서 요청 시뮬레이션
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/my", nil)
		req.AddCookie(&http.Cookie{
			Name:  "accessToken",
			Value: validToken,
		})
		req.Header.Set("Origin", "https://malicious-site.com")
		router.ServeHTTP(w, req)

		// Then: 인증 실패 (CSRF 방어)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("XSS 공격 방어 - 악성 스크립트가 포함된 요청도 정상 처리", func(t *testing.T) {
		// Given: 인증 서비스 모킹
		mockAuthService := new(MockAuthService)
		authHandler := handlers.NewAuthHandler(mockAuthService)

		// 로그인 성공 응답 설정
		expectedUser := &models.User{
			UserID: "testuser",
			Name:   "Test User",
			Email:  "test@example.com",
			Role:   models.UserRoleNormal,
		}
		expectedUser.ID = 1

		mockAuthService.On("Login", mock.Anything, "testuser", "password123", "127.0.0.1", mock.Anything).
			Return(&services.LoginResponse{
				User:                 expectedUser,
				AccessToken:          "mock.access.token",
				RefreshToken:         "mock.refresh.token",
				AccessTokenExpiresIn: time.Now().Add(15 * time.Minute).Unix(),
			}, nil)

		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/v1/auth/login", authHandler.Login)

		// When: 악성 스크립트가 포함된 User-Agent로 로그인 시도
		maliciousUserAgent := "<script>alert('XSS')</script>"
		loginData := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", maliciousUserAgent)
		req.Header.Set("X-Forwarded-For", "127.0.0.1")
		router.ServeHTTP(w, req)

		// Then: 로그인 성공 (악성 스크립트는 실행되지 않음)
		assert.Equal(t, http.StatusOK, w.Code)

		// 응답에 토큰이 포함되어 있는지 확인
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		if value, ok := response["value"].(map[string]interface{}); ok {
			assert.Contains(t, value, "accessToken")
		} else {
			assert.Contains(t, response, "accessToken")
		}
	})

	t.Run("Race Condition 테스트 - 동시 로그인 요청 처리", func(t *testing.T) {
		// Given: 인증 서비스 모킹
		mockAuthService := new(MockAuthService)
		authHandler := handlers.NewAuthHandler(mockAuthService)

		expectedUser := &models.User{
			UserID: "testuser",
			Name:   "Test User",
			Role:   models.UserRoleNormal,
		}
		expectedUser.ID = 1

		// 동시 요청에 대해 모든 호출이 성공하도록 설정
		mockAuthService.On("Login", mock.Anything, "testuser", "password123", mock.AnythingOfType("string"), mock.Anything).
			Return(&services.LoginResponse{
				User:                 expectedUser,
				AccessToken:          "mock.access.token",
				RefreshToken:         "mock.refresh.token",
				AccessTokenExpiresIn: time.Now().Add(15 * time.Minute).Unix(),
			}, nil)

		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/v1/auth/login", authHandler.Login)

		// When: 여러 goroutine에서 동시에 로그인 요청
		const numConcurrentRequests = 10
		var wg sync.WaitGroup
		var successCount int64

		for i := 0; i < numConcurrentRequests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				loginData := map[string]string{
					"username": "testuser",
					"password": "password123",
				}
				jsonData, _ := json.Marshal(loginData)

				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Forwarded-For", "127.0.0.1")
				router.ServeHTTP(w, req)

				if w.Code == http.StatusOK {
					atomic.AddInt64(&successCount, 1)
				}
			}()
		}

		wg.Wait()

		// Then: 모든 요청이 성공적으로 처리되었는지 확인
		assert.Equal(t, int64(numConcurrentRequests), atomic.LoadInt64(&successCount))
	})

	t.Run("동시 로그인 허용 - 서로 다른 디바이스에서 로그인 가능", func(t *testing.T) {
		// Given: JWT 서비스 설정
		jwtService := auth.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
		mockAuthService := new(MockAuthService)
		authHandler := handlers.NewAuthHandler(mockAuthService)

		expectedUser := &models.User{
			UserID: "testuser",
			Name:   "Test User",
			Role:   models.UserRoleNormal,
		}
		expectedUser.ID = 1

		// 첫 번째 디바이스 로그인
		mockAuthService.On("Login", mock.Anything, "testuser", "password123", "192.168.1.100", mock.Anything).
			Return(&services.LoginResponse{
				User:                 expectedUser,
				AccessToken:          "first.device.token",
				RefreshToken:         "first.device.refresh",
				AccessTokenExpiresIn: time.Now().Add(15 * time.Minute).Unix(),
			}, nil).Once()

		// 두 번째 디바이스 로그인
		mockAuthService.On("Login", mock.Anything, "testuser", "password123", "192.168.1.200", mock.Anything).
			Return(&services.LoginResponse{
				User:                 expectedUser,
				AccessToken:          "second.device.token",
				RefreshToken:         "second.device.refresh",
				AccessTokenExpiresIn: time.Now().Add(15 * time.Minute).Unix(),
			}, nil).Once()

		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/v1/auth/login", authHandler.Login)
		router.Use(middleware.AuthMiddleware(jwtService))
		router.GET("/api/v1/my", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		loginData := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		// When: 첫 번째 디바이스에서 로그인
		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		req1.Header.Set("X-Forwarded-For", "192.168.1.100")
		router.ServeHTTP(w1, req1)

		// And: 두 번째 디바이스에서 로그인
		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0)")
		req2.Header.Set("X-Forwarded-For", "192.168.1.200")
		router.ServeHTTP(w2, req2)

		// Then: 두 디바이스 모두 로그인 성공
		assert.Equal(t, http.StatusOK, w1.Code)
		assert.Equal(t, http.StatusOK, w2.Code)

		// 각각 다른 토큰을 받았는지 확인
		var response1, response2 map[string]interface{}
		json.Unmarshal(w1.Body.Bytes(), &response1)
		json.Unmarshal(w2.Body.Bytes(), &response2)

		// value 필드에서 accessToken 추출
		var token1, token2 interface{}
		if value1, ok := response1["value"].(map[string]interface{}); ok {
			token1 = value1["accessToken"]
		} else {
			token1 = response1["accessToken"]
		}
		if value2, ok := response2["value"].(map[string]interface{}); ok {
			token2 = value2["accessToken"]
		} else {
			token2 = response2["accessToken"]
		}

		assert.NotEqual(t, token1, token2)
	})

	t.Run("HTTP 헤더 조작 방어", func(t *testing.T) {
		// Given: 헤더 조작이 포함된 요청
		mockAuthService := new(MockAuthService)
		authHandler := handlers.NewAuthHandler(mockAuthService)

		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.POST("/api/v1/auth/login", authHandler.Login)

		// 잘못된 로그인으로 설정
		mockAuthService.On("Login", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, services.ErrInvalidCredentials)

		loginData := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		// When: 다양한 헤더 조작 시도
		manipulatedHeaders := map[string]string{
			"X-Forwarded-For":     "127.0.0.1, 192.168.1.1, 10.0.0.1", // IP 스푸핑 시도
			"X-Real-IP":           "admin.localhost",                  // 관리자 IP로 위장
			"X-Originating-IP":    "192.168.1.1",
			"X-Remote-IP":         "10.0.0.1",
			"X-Cluster-Client-IP": "172.16.0.1",
		}

		for headerName, headerValue := range manipulatedHeaders {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set(headerName, headerValue)
			router.ServeHTTP(w, req)

			// Then: 헤더 조작과 관계없이 로그인 실패 (잘못된 비밀번호)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		}
	})
}
