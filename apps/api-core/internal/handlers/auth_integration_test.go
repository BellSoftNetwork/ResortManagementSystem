//go:build integration

package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/database"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/handlers"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AuthIntegrationTestSuite struct {
	suite.Suite
	db          *gorm.DB
	redis       *redis.Client
	miniRedis   *miniredis.Miniredis
	router      *gin.Engine
	jwtService  *auth.JWTService
	authService services.AuthService
	userService services.UserService
}

func (suite *AuthIntegrationTestSuite) SetupSuite() {
	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.NoError(err)
	suite.db = db

	// Run migrations
	err = database.Migrate(db)
	suite.NoError(err)

	// Setup mini Redis
	mr, err := miniredis.Run()
	suite.NoError(err)
	suite.miniRedis = mr

	suite.redis = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Setup services
	suite.jwtService = auth.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)

	userRepo := repositories.NewUserRepository(db)
	loginAttemptRepo := repositories.NewLoginAttemptRepository(db)

	suite.authService = services.NewAuthService(userRepo, loginAttemptRepo, suite.jwtService, suite.redis)
	suite.userService = services.NewUserService(userRepo)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(suite.authService)
	userHandler := handlers.NewUserHandler(suite.userService)

	// Setup router
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.router.Use(middleware.ErrorHandler())

	// Setup routes
	api := suite.router.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
		}

		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware(suite.jwtService))
		{
			myRoutes := authenticated.Group("/my")
			{
				myRoutes.GET("", userHandler.GetCurrentUser)
				myRoutes.POST("", userHandler.GetCurrentUser)
				myRoutes.PATCH("", userHandler.UpdateCurrentUser)
			}
		}
	}
}

func (suite *AuthIntegrationTestSuite) TearDownSuite() {
	suite.miniRedis.Close()
}

func (suite *AuthIntegrationTestSuite) TearDownTest() {
	// Clean up database
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM login_attempts")
	// Clean up Redis
	suite.redis.FlushAll(context.Background())
}

func (suite *AuthIntegrationTestSuite) TestFullAuthenticationFlow() {
	// 1. 회원가입
	registerReq := dto.RegisterRequest{
		UserID:   "testuser123",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "testpassword123",
	}

	body, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)

	// 2. 로그인
	loginReq := dto.LoginRequest{
		Username: "testuser123",
		Password: "testpassword123",
	}

	body, _ = json.Marshal(loginReq)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var loginResp dto.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResp)
	suite.NoError(err)
	suite.NotEmpty(loginResp.AccessToken)
	suite.NotEmpty(loginResp.RefreshToken)

	// 3. 액세스 토큰으로 /my API 호출
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/my", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var userResp dto.UserResponse
	err = json.Unmarshal(w.Body.Bytes(), &userResp)
	suite.NoError(err)
	suite.Equal("testuser123", userResp.UserID)
	suite.Equal("test@example.com", userResp.Email)
	suite.Equal("Test User", userResp.Name)

	// 4. 리프레시 토큰으로 새 토큰 발급
	refreshReq := dto.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	}

	body, _ = json.Marshal(refreshReq)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var refreshResp dto.TokenResponse
	err = json.Unmarshal(w.Body.Bytes(), &refreshResp)
	suite.NoError(err)
	suite.NotEmpty(refreshResp.AccessToken)
	suite.NotEmpty(refreshResp.RefreshToken)

	// 5. 새 액세스 토큰으로 /my API 호출
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/my", nil) // POST 메서드도 지원하는지 확인
	req.Header.Set("Authorization", "Bearer "+refreshResp.AccessToken)

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *AuthIntegrationTestSuite) TestInvalidToken() {
	// 잘못된 토큰으로 /my API 호출
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/my", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)

	var errResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.NoError(err)
	suite.Equal("Invalid or expired token", errResp["message"])
}

func (suite *AuthIntegrationTestSuite) TestExpiredRefreshToken() {
	// 먼저 정상적으로 로그인
	// 회원가입 먼저
	registerReq := dto.RegisterRequest{
		UserID:   "testuser123",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "testpassword123",
	}

	body, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)

	// 로그인
	loginReq := dto.LoginRequest{
		Username: "testuser123",
		Password: "testpassword123",
	}

	body, _ = json.Marshal(loginReq)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)

	var loginResp dto.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)

	// Redis에서 refresh token 삭제
	suite.redis.Del(context.Background(), fmt.Sprintf("refresh_token:1"))

	// 리프레시 시도
	refreshReq := dto.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	}

	body, _ = json.Marshal(refreshReq)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)

	var errResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	suite.NoError(err)
	suite.Equal("Invalid or expired refresh token", errResp["message"])
}

func (suite *AuthIntegrationTestSuite) TestTokenStructureCompatibility() {
	// api-legacy와 동일한 토큰 구조인지 확인

	// 회원가입
	registerReq := dto.RegisterRequest{
		UserID:   "testadmin",
		Email:    "admin@example.com",
		Name:     "Admin User",
		Password: "adminpassword123",
	}

	body, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)

	// 로그인
	loginReq := dto.LoginRequest{
		Username: "testadmin",
		Password: "adminpassword123",
	}

	body, _ = json.Marshal(loginReq)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(w, req)

	var loginResp dto.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResp)

	// 토큰 파싱하여 claims 확인
	claims, err := suite.jwtService.ValidateAccessToken(loginResp.AccessToken)
	suite.NoError(err)
	suite.Equal("1", claims.Subject)          // user ID
	suite.Equal("testadmin", claims.Username) // UserID (not email!)
	suite.Equal("NORMAL", claims.Authorities) // role
}

func TestAuthIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}
