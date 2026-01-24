package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/pkg/auth"
)

func TestAuthMiddleware(t *testing.T) {
	// Given: JWT 서비스 설정
	gin.SetMode(gin.TestMode)
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)

	tests := []struct {
		name           string
		setupAuth      func(*http.Request)
		expectedStatus int
		expectedBody   string
		checkContext   func(*gin.Context) bool
	}{
		{
			name: "유효한 토큰으로 요청 시 인증에 성공한다",
			setupAuth: func(req *http.Request) {
				token, _ := jwtService.GenerateAccessToken(1, "testuser", "ADMIN")
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"authenticated"}`,
			checkContext: func(c *gin.Context) bool {
				userID, exists := middleware.GetUserID(c)
				username, _ := middleware.GetUsername(c)
				role, _ := middleware.GetUserRole(c)
				return exists && userID == 1 && username == "testuser" && role == "ADMIN"
			},
		},
		{
			name: "Authorization 헤더가 없으면 401을 반환한다",
			setupAuth: func(req *http.Request) {
				// 헤더 설정 안함
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":null,"fieldErrors":null,"message":"인증 헤더가 없습니다"}`,
		},
		{
			name: "Bearer 접두사가 없으면 401을 반환한다",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "InvalidToken")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":null,"fieldErrors":null,"message":"잘못된 인증 헤더 형식"}`,
		},
		{
			name: "토큰이 비어있으면 401을 반환한다",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer ")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":null,"fieldErrors":null,"message":"토큰이 없습니다"}`,
		},
		{
			name: "만료된 토큰으로 요청 시 401을 반환한다",
			setupAuth: func(req *http.Request) {
				// 만료된 토큰 생성
				claims := &auth.Claims{
					Username:    "testuser",
					Authorities: "ADMIN",
				}
				claims.Subject = "1"
				claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("test-secret"))
				req.Header.Set("Authorization", "Bearer "+tokenString)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":null,"fieldErrors":null,"message":"유효하지 않거나 만료된 토큰"}`,
		},
		{
			name: "잘못된 서명의 토큰으로 요청 시 401을 반환한다",
			setupAuth: func(req *http.Request) {
				// 다른 시크릿으로 생성한 토큰
				wrongService := auth.NewJWTService("wrong-secret", 15*time.Minute, 7*24*time.Hour)
				token, _ := wrongService.GenerateAccessToken(1, "testuser", "ADMIN")
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":null,"fieldErrors":null,"message":"유효하지 않거나 만료된 토큰"}`,
		},
		{
			name: "Spring Boot 형식의 토큰이 정상적으로 파싱된다",
			setupAuth: func(req *http.Request) {
				// Spring Boot 스타일 토큰 (authorities 필드 사용)
				token, _ := jwtService.GenerateAccessToken(1, "testadmin", "SUPER_ADMIN")
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"authenticated"}`,
			checkContext: func(c *gin.Context) bool {
				role, _ := middleware.GetUserRole(c)
				return role == "SUPER_ADMIN"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: 라우터 설정
			router := gin.New()
			router.Use(middleware.AuthMiddleware(jwtService))
			router.GET("/test", func(c *gin.Context) {
				if tt.checkContext != nil {
					assert.True(t, tt.checkContext(c))
				}
				c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
			})

			// When: 요청 실행
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupAuth(req)
			router.ServeHTTP(w, req)

			// Then: 응답 검증
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestRoleMiddleware(t *testing.T) {
	// Given: 테스트 환경 설정
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       string
		allowedRoles   []string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "허용된 역할로 요청 시 접근이 가능하다",
			userRole:       "ADMIN",
			allowedRoles:   []string{"ADMIN", "SUPER_ADMIN"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"authorized"}`,
		},
		{
			name:           "SUPER_ADMIN은 ADMIN 권한이 필요한 곳에 접근 가능하다",
			userRole:       "SUPER_ADMIN",
			allowedRoles:   []string{"ADMIN", "SUPER_ADMIN"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"authorized"}`,
		},
		{
			name:           "권한이 없으면 403을 반환한다",
			userRole:       "USER",
			allowedRoles:   []string{"ADMIN", "SUPER_ADMIN"},
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"errors":null,"fieldErrors":null,"message":"권한이 부족합니다"}`,
		},
		{
			name:           "일반 사용자는 관리자 기능에 접근할 수 없다",
			userRole:       "USER",
			allowedRoles:   []string{"ADMIN"},
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"errors":null,"fieldErrors":null,"message":"권한이 부족합니다"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: 라우터 설정
			router := gin.New()
			router.Use(func(c *gin.Context) {
				// 인증 정보 설정
				c.Set(middleware.UserIDKey, uint(1))
				c.Set(middleware.UsernameKey, "testuser")
				c.Set(middleware.UserRoleKey, tt.userRole)
				c.Next()
			})
			router.Use(middleware.RoleMiddleware(tt.allowedRoles...))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "authorized"})
			})

			// When: 요청 실행
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			// Then: 응답 검증
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestAuthMiddleware_SecurityTests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)

	t.Run("쿠키에 토큰이 있어도 헤더에 없으면 인증 실패", func(t *testing.T) {
		// Given: 유효한 토큰 생성
		token, err := jwtService.GenerateAccessToken(1, "testuser", "USER")
		assert.NoError(t, err)

		router := gin.New()
		router.Use(middleware.AuthMiddleware(jwtService))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
		})

		// When: 쿠키에만 토큰을 넣고 요청
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:  "accessToken",
			Value: token,
		})
		router.ServeHTTP(w, req)

		// Then: 인증 실패
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"errors":null,"fieldErrors":null,"message":"인증 헤더가 없습니다"}`, w.Body.String())
	})

	t.Run("URL 파라미터에 토큰이 있어도 헤더에 없으면 인증 실패", func(t *testing.T) {
		// Given: 유효한 토큰 생성
		token, err := jwtService.GenerateAccessToken(1, "testuser", "USER")
		assert.NoError(t, err)

		router := gin.New()
		router.Use(middleware.AuthMiddleware(jwtService))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
		})

		// When: URL 파라미터에만 토큰을 넣고 요청
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test?token="+token, nil)
		router.ServeHTTP(w, req)

		// Then: 인증 실패
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"errors":null,"fieldErrors":null,"message":"인증 헤더가 없습니다"}`, w.Body.String())
	})

	t.Run("요청 본문에 토큰이 있어도 헤더에 없으면 인증 실패", func(t *testing.T) {
		// Given: 유효한 토큰 생성
		token, err := jwtService.GenerateAccessToken(1, "testuser", "USER")
		assert.NoError(t, err)

		router := gin.New()
		router.Use(middleware.AuthMiddleware(jwtService))
		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
		})

		// When: 요청 본문에만 토큰을 넣고 요청
		body := `{"token":"` + token + `"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Then: 인증 실패
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"errors":null,"fieldErrors":null,"message":"인증 헤더가 없습니다"}`, w.Body.String())
	})

	t.Run("대소문자 변형된 Bearer 헤더는 인증 실패", func(t *testing.T) {
		// Given: 유효한 토큰 생성
		token, err := jwtService.GenerateAccessToken(1, "testuser", "USER")
		assert.NoError(t, err)

		router := gin.New()
		router.Use(middleware.AuthMiddleware(jwtService))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
		})

		testCases := []string{
			"bearer " + token, // 소문자
			"BEARER " + token, // 대문자
			"Bearer" + token,  // 공백 없음
			"Token " + token,  // 다른 접두사
		}

		for _, authHeader := range testCases {
			// When: 잘못된 형식의 헤더로 요청
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", authHeader)
			router.ServeHTTP(w, req)

			// Then: 인증 실패
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("토큰에 공백이나 특수문자가 포함된 경우 처리", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.AuthMiddleware(jwtService))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
		})

		testCases := []string{
			"Bearer ",              // 토큰 없음
			"Bearer  ",             // 공백만 있음
			"Bearer \t",            // 탭 문자
			"Bearer \n",            // 개행 문자
			"Bearer invalid token", // 공백 포함 토큰
		}

		for _, authHeader := range testCases {
			// When: 잘못된 토큰으로 요청
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", authHeader)
			router.ServeHTTP(w, req)

			// Then: 인증 실패
			assert.NotEqual(t, http.StatusOK, w.Code)
		}
	})

	t.Run("변조된 JWT 헤더는 거부된다", func(t *testing.T) {
		// Given: 변조된 헤더를 가진 토큰
		// 정상 헤더: {"alg":"HS256","typ":"JWT"}
		// 변조 헤더: {"alg":"none","typ":"JWT"} - 서명 없음 시도
		tamperedToken := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiYXV0aG9yaXRpZXMiOiJVU0VSIiwic3ViIjoiMSIsImV4cCI6MTc1MDAwMDAwMCwiaWF0IjoxNzAwMDAwMDAwfQ."

		router := gin.New()
		router.Use(middleware.AuthMiddleware(jwtService))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
		})

		// When: 변조된 토큰으로 요청
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tamperedToken)
		router.ServeHTTP(w, req)

		// Then: 인증 실패
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"errors":null,"fieldErrors":null,"message":"유효하지 않거나 만료된 토큰"}`, w.Body.String())
	})

	t.Run("Multiple Authorization 헤더가 있는 경우 처리", func(t *testing.T) {
		// Given: 유효한 토큰과 잘못된 토큰
		validToken, err := jwtService.GenerateAccessToken(1, "testuser", "USER")
		assert.NoError(t, err)

		router := gin.New()
		router.Use(middleware.AuthMiddleware(jwtService))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
		})

		// When: 여러 Authorization 헤더가 있는 요청
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Add("Authorization", "Bearer invalid")
		req.Header.Add("Authorization", "Bearer "+validToken)
		router.ServeHTTP(w, req)

		// Then: 첫 번째 헤더를 사용하므로 인증 실패
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
